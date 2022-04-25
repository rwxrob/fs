package fs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fs"
)

/*
func ExampleIsDir() {
	fmt.Println(fs.IsDir("testdata"))
	fmt.Println(fs.IsDir("testdata/fail"))
	// Output:
	// true
	// false
}

func ExampleDupPerms() {
	os.Mkdir(`testdata/some`, 0000)
	defer os.Remove(`testdata/some`)
	stats, _ := os.Stat(`testdata/orig`)
	fmt.Println(stats.Mode())
	stats, _ = os.Stat(`testdata/some`)
	fmt.Println(stats.Mode())
	err := fs.DupPerms(`testdata/orig`, `testdata/some`)
	if err != nil {
		fmt.Println(err)
	}
	stats, _ = os.Stat(`testdata/some`)
	fmt.Println(stats.Mode())
	// Output:
	// drw-------
	// d---------
	// drw-------
}
*/

func ExampleExists() {
	fmt.Println(fs.Exists("testdata/file")) // use NotExists instead of !
	// Output:
	// false
}

func ExampleNotExists() {
	fmt.Println(fs.NotExists("testdata/nope")) // use Exists instead of !
	// Output:
	// true
}

func ExampleModTime() {
	fmt.Println(fs.ModTime("testdata/file").IsZero())
	fmt.Println(fs.ModTime("testdata/none"))
	// Output:
	// true
	// 0001-01-01 00:00:00 +0000 UTC
}

func ExampleHereOrAbove_here() {
	dir, _ := os.Getwd()
	defer func() { os.Chdir(dir) }()
	os.Chdir("testdata/adir")

	path, err := fs.HereOrAbove("afile")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [adir afile]

}

func ExampleHereOrAbove_above() {
	dir, _ := os.Getwd()
	defer func() { os.Chdir(dir) }()
	os.Chdir("testdata/adir")

	path, err := fs.HereOrAbove("anotherfile")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [testdata anotherfile]

}
