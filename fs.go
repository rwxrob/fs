package fs

import (
	"log"
	"os"
)

// IsDir returns true if the indicated path is a directory. Returns
// false and logs if unable to determine.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		log.Print(err)
		return false
	}
	return info.IsDir()
}
