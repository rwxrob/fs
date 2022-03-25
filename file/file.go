package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rwxrob/fs"
)

// Touch creates a new file at path or updates the time stamp of existing.
func Touch(path string) error {
	if fs.NotExists(path) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}
	now := time.Now().Local()
	if err := os.Chtimes(path, now, now); err != nil {
		return err
	}
	return nil
}

// Replace replaces a file at a specified location with another
// successfully retrieved file from the specified URL or file path and
// duplicates the original files permissions. Only http and https URLs
// are currently supported. For security reasons, no backup copy of the
// replaced executable is kept. Also no checksum validation of the file
// is performed (which is fine in most cases where the connection has
// been secured with HTTPS).
func Replace(orig, url string) error {
	if err := Fetch(url, orig+`.new`); err != nil {
		return err
	}
	if err := fs.DupPerms(orig, orig+`.new`); err != nil {
		return err
	}
	if err := os.Rename(orig, orig+`.orig`); err != nil {
		return err
	}
	if err := os.Rename(orig+`.new`, orig); err != nil {
		return err
	}
	if err := os.Remove(orig + `.orig`); err != nil {
		return err
	}
	return nil
}

// Fetch fetches the specified file at the give "from" URL and saves it
// "to" the specified file path. The name is *not* inferred. If
// timeouts, status, and contexts are required use the net/http package
// instead. Will block until the entire file is downloaded. For more
// involved downloading needs consider the github.com/cavaliercoder/grab
// package.
func Fetch(from, to string) error {

	file, err := os.Create(to)
	defer file.Close()
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Get(from)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return fmt.Errorf(res.Status)
	}

	if _, err := io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}
