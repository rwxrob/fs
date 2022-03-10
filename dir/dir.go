package dir

import (
	"os"
	"path/filepath"

	"github.com/rwxrob/json"
)

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
