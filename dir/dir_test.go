package dir_test

import (
	"fmt"

	"github.com/rwxrob/fs/dir"
)

func ExampleEntries() {
	list := dir.Entries("testdata")
	fmt.Println(list)
	// Output:
	// [testdata/file testdata/other]
}

func ExampleExists() {
	fmt.Println(dir.Exists("testdata/exists"))
	fmt.Println(dir.Exists("testdata"))
	// Output:
	// false
	// true
}
