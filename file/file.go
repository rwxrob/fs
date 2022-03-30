package file

import (
	"fmt"
	"io"
	_fs "io/fs"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/rwxrob/fs"
)

// DefaultPerms for new file creation.
var DefaultPerms = 0600

// Touch creates a new file at path or updates the time stamp of
// existing. If a new file is needed creates it with 0600 permissions
// (instead of 0666 as default os.Create does). The directory must
// already exist.
func Touch(path string) error {
	if fs.NotExists(path) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		file.Close()
		return os.Chmod(path, _fs.FileMode(DefaultPerms))
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

// execute executes the given arguments using a syscall so as to hand over
// all the associated running process references and resources.
func execute(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing name of executable")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	// exits the program unless there is an error
	return syscall.Exec(path, args, os.Environ())
}

// Edit opens the file at the given path for editing searching for an
// editor on the system using the following (in order of priority):
//
// * VISUAL
// * EDITOR
// * code
// * kakoune
// * vim
// * vi
// * nano
//
// Currently, only architectures that support syscall.Exec are used but
// all supported Go architectures are planned.
func Edit(path string) error {
	ed := os.Getenv("VISUAL")
	if ed != "" {
		return execute(ed, path)
	}
	ed = os.Getenv("EDITOR")
	if ed != "" {
		return execute(ed, path)
	}
	ed, _ = exec.LookPath("code")
	if ed != "" {
		return execute(ed, path)
	}
	ed, _ = exec.LookPath("kak")
	if ed != "" {
		return execute(ed, path)
	}
	ed, _ = exec.LookPath("vim")
	if ed != "" {
		return execute(ed, path)
	}
	ed, _ = exec.LookPath("vi")
	if ed != "" {
		return execute(ed, path)
	}
	ed, _ = exec.LookPath("nano")
	if ed != "" {
		return execute(ed, path)
	}
	return fmt.Errorf("unable to find editor")
}

// Exists calls fs.Exists and further confirms that the file is a file
// and not a directory.
func Exists(path string) bool { return fs.Exists(path) && !fs.IsDir(path) }
