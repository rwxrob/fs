package fs

import (
	"embed"
	"io/fs"
	_fs "io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Tilde2Home expands a Tilde (~) prefix into a proper os.UserHomeDir path. If
// it cannot find os.UserHomeDir simple returns unchanged path. Will not
// expand for specific users (~username).
func Tilde2Home(dir string) string {
	if !strings.HasPrefix(dir, `~`) {
		return dir
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		return path.Join(home, dir[1:])
	}
	return dir
}

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

// DupPerms duplicates the perms of one file or directory onto the second.
func DupPerms(orig, clone string) error {
	ostats, err := os.Stat(orig)
	if err != nil {
		return err
	}
	_, err = os.Stat(clone)
	if err != nil {
		return nil
	}
	return os.Chmod(clone, ostats.Mode())
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

// Exists returns true if the given path was absolutely found to exist
// on the system. A false return value means either the file does not
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

// HereOrAbove returns the full path to the file or directory if it is
// found in the current working directory, or if not exists in any
// parent directory recursively.
func HereOrAbove(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for ; len(dir) > 0 && dir != "/"; dir = filepath.Dir(dir) {
		path := filepath.Join(dir, name)
		if Exists(path) {
			return path, nil
		}
	}
	return "", nil
}

// IsDirFS simply a shortcut for fs.Stat().IsDir(). Only returns true if
// the path is a directory. If not a directory (or an error prevented
// confirming it is a directory) then returns false.
func IsDirFS(fsys _fs.FS, path string) bool {
	info, err := _fs.Stat(fsys, path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

var (
	ExtractFilePerms = _fs.FileMode(0600)
	ExtractDirPerms  = _fs.FileMode(0700)
)

// ExtractEmbed walks the embedded file system and duplicates its
// structure into the target directory beginning with root. Since
// embed.FS drops all permissions information from the original files
// all files and directories are written as read/write (0600) for the
// current effective user (0600 for file, 0700 for directories). This
// default can be changed by setting the package variables
// ExtractFilePerms and ExtractDirPerms. Note that each embedded file is
// full buffered into memory before writing.
func ExtractEmbed(it embed.FS, root, target string) error {
	return _fs.WalkDir(it, root,

		func(path string, i _fs.DirEntry, err error) error {

			to := filepath.Join(target, strings.TrimPrefix(path, root))

			if i.IsDir() {
				return os.MkdirAll(to, ExtractDirPerms)
			}

			buf, err := _fs.ReadFile(it, path)
			if err != nil {
				return err
			}

			return os.WriteFile(to, buf, ExtractFilePerms)
		})

}

// Paths returns a list of full paths to each of the directories or
// files from the root but with the path to the root stripped resulting
// in relative paths.
func RelPaths(it fs.FS, root string) []string {
	var paths []string
	_fs.WalkDir(it, root,
		func(path string, i _fs.DirEntry, err error) error {
			to := strings.TrimPrefix(path, root)
			if to == "" {
				return nil
			}
			paths = append(paths, to[1:])
			return nil
		})
	return paths
}

// LatestChange walks the directory rooted at root looking at each file or
// directory within it recursively (including itself) and returns the
// time of the most recent change along with its full path. LastChange
// returns nil tile and empty string if root does not exist.
func LatestChange(root string) (string, fs.FileInfo) {
	latest := struct {
		Info fs.FileInfo
		Path string
	}{}
	err := filepath.WalkDir(root, func(p string, f fs.DirEntry, _ error) error {
		if latest.Path == "" {
			i, err := f.Info()
			if err != nil {
				return err
			}
			latest.Info = i
			latest.Path = p
			return nil
		}
		i, err := f.Info()
		if err != nil {
			return nil
		}
		if i.ModTime().After(latest.Info.ModTime()) {
			latest.Info = i
			latest.Path = p
		}
		return nil
	})
	if err != nil {
		return "", nil
	}
	return latest.Path, latest.Info
}
