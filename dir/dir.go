package dir

import (
	"io/fs"
	"os"
	"path/filepath"

	_fs "github.com/rwxrob/fs"
)

// DefaultPerms are defaults for new directory creation.
var DefaultPerms = 0700

// Create a new directory with the DefaultPerms creating any new
// directories as well (see os.MkdirAll)
func Create(path string) error {
	return os.MkdirAll(path, fs.FileMode(DefaultPerms))
}

// In returns a slice of strings with all the files in the directory
// at that path joined to their path (as is usually wanted). Returns an
// empty slice if empty or path doesn't point to a directory. See List.
func Entries(path string) []string {
	var list []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return list
	}
	for _, f := range entries {
		list = append(list, filepath.Join(path, f.Name()))
	}
	return list
}

// EntriesWithSlash returns Entries passed to AddSlash so that all
// directories will have a trailing slash.
func EntriesWithSlash(path string) []string {
	return AddSlash(Entries(path))
}

// AddSlash adds a filepath.Separator to the end of all entries passed
// that are directories.
func AddSlash(entries []string) []string {
	var list []string
	for _, entry := range entries {
		if _fs.IsDir(entry) {
			entry += string(filepath.Separator)
		}
		list = append(list, entry)
	}
	return list
}

// Exists calls fs.Exists and further confirms that the path is
// a directory and not a file.
func Exists(path string) bool { return _fs.Exists(path) && _fs.IsDir(path) }

// Name returns the current working directory name or an empty string.
func Name() string {
	wd, _ := os.Getwd()
	return filepath.Base(wd)
}
