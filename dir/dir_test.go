package dir_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fs/dir"
)

func ExampleEntries() {
	list := dir.Entries("testdata")
	fmt.Println(list)
	// Output:
	// [testdata/adir testdata/file testdata/other]

}

func ExampleExists() {
	fmt.Println(dir.Exists("testdata/exists"))
	fmt.Println(dir.Exists("testdata"))
	// Output:
	// false
	// true
}

func ExampleName() {
	os.Chdir(`testdata`)
	fmt.Println(dir.Name())
	os.Chdir(`..`)
	// Output:
	// testdata
}

func ExampleHereOrAbove_here() {
	_dir, _ := os.Getwd()
	defer func() { os.Chdir(_dir) }()
	os.Chdir("testdata/adir")

	path, err := dir.HereOrAbove("anotherdir")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [adir anotherdir]

}

func ExampleHereOrAbove_above() {
	_dir, _ := os.Getwd()
	defer func() { os.Chdir(_dir) }()
	os.Chdir("testdata/adir/anotherdir")

	path, err := dir.HereOrAbove("adir")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [testdata adir]

}

/*
func ExampleLatestChange() {
	path, _ := dir.LatestChange("testdata")
	fmt.Println(path)
	// Output:
	// whichever
}
*/
