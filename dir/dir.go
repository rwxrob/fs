package dir

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rwxrob/json"
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
func Entries(path string) json.Array {
	list := json.Array{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return list
	}
	for _, f := range entries {
		list = append(list, filepath.Join(path, f.Name()))
	}
	return list
}
