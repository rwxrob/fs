package file

import (
	"os"
	"time"
)

// Exists returns true if the given path was absolutely found to exist on
// the system. A false return value means either the file does not
// exists or it was not able to determine if it exists or not. Use
// NotExists instead.
//
// WARNING: do not use this function if a definitive check for the
// non-existence of a file is required since the possible indeterminate
// error state is a possibility. These checks are also not atomic on
// many file systems so avoid this usage for pseudo-semaphore designs
// and depend on file locks.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// NotExists definitively returns true if the given path does not exist.
// See Exists as well.
func NotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// Touch creates a new file at path or updates the time stamp of existing.
func Touch(path string) error {
	if NotExists(path) {
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

// ModTime returns the FileInfo.ModTime() from Stat() as a convenience
// or returns a zero time if not. See time.IsZero.
func ModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err == nil {
		return info.ModTime()
	}
	var t time.Time
	return t
}
