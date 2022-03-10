package dir_test

import "github.com/rwxrob/fs/dir"

func ExampleEntries() {
	list := dir.Entries("testdata")
	list.Print()
	// Output:
	// ["testdata/file","testdata/other"]
}
