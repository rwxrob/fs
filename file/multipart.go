package file

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"
)

// Multipart is meant to contain the delimited sections of output and can
// be marshalled into a single delimited string safely and automatically
// simply by using it in a string context.
type Multipart map[string]string

// MarshalText fulfills the encoding.TextMarshaler interface by
// delimiting each section of output with a unique delimiter line that
// contains a space and the key for each section. Order of sections is
// indeterminate officially (but consistent for testing, per Go). The
// special "end" delimiter is always the last line.
func (o Multipart) MarshalText() ([]byte, error) {
	var out string
	delim := _base32()
	if delim == "" {
		return nil, fmt.Errorf(`unable to get random data`)
	}
	for k, v := range o {
		out += delim + " " + k + "\n" + v + "\n"
	}
	out += delim + " end"
	return []byte(out), nil
}

// UnmarshalText fulfills the encoding.TextUnmarshaler interface by
// sensing the delimiter as the first text field (up to the first space)
// and using that delimiter to parse the remaining data into the
// key/value pairs ending when either the end of text is encountered or
// the special "end" delimiter is read.
func (o Multipart) UnmarshalText(text []byte) error {

	s := bufio.NewScanner(bytes.NewReader(text))
	if !s.Scan() {
		return fmt.Errorf(`failed to scan first line`)
	}

	f := strings.Fields(s.Text())
	if len(f) < 2 {
		return fmt.Errorf(`first line is not delimiter`)
	}

	if f[1] == "end" {
		return nil
	}

	delim := f[0]
	cur := f[1]

	for s.Scan() {
		line := s.Text()

		// delimiter?
		if strings.HasPrefix(line, delim) {
			f := strings.Fields(line)
			if len(f) < 2 {
				return fmt.Errorf(`delimiter missing key`)
			}
			o[cur] = o[cur][:len(o[cur])-1]
			if f[1] == `end` {
				return nil
			}
			cur = f[1]
			continue
		}

		o[cur] += line + "\n"

	}

	return nil
}

// String fulfills the fmt.Stringer interface by calling MarshalText.
func (o Multipart) String() string {
	buf, err := o.MarshalText()
	if err != nil {
		return ""
	}
	return string(buf)
}

// from github.com/rwxrob/uniq
func _base32() string {
	byt := make([]byte, 20)
	_, err := rand.Read(byt)
	if err != nil {
		return ""
	}
	return base32.HexEncoding.EncodeToString(byt)
}
