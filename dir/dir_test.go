package dir_test

import (
	"fmt"
	"os"

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

func ExampleName() {
	os.Chdir(`testdata`)
	fmt.Println(dir.Name())
	os.Chdir(`..`)
	// Output:
	// testdata
}
