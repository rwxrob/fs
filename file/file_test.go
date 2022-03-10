package file_test

import (
	"fmt"
	"log"
	"os"

	"github.com/rwxrob/fs/file"
)

func ExampleExists() {
	fmt.Println(file.Exists("testdata/file")) // use NotExists instead of !
	// Output:
	// true
}

func ExampleNotExists() {
	fmt.Println(file.NotExists("testdata/nope")) // use Exists instead of !
	// Output:
	// true
}

func ExampleTouch_create() {
	fmt.Println(file.NotExists("testdata/foo"))
	file.Touch("testdata/foo")
	fmt.Println(file.Exists("testdata/foo"))
	os.Remove("testdata/foo")
	// Output:
	// true
	// true
}

func ExampleModTime() {
	fmt.Println(file.ModTime("testdata/file").IsZero())
	fmt.Println(file.ModTime("testdata/none"))
	// Output:
	// false
	// 0001-01-01 00:00:00 +0000 UTC
}

func ExampleTouch_update() {

	// first create it and capture the time as a string
	file.Touch("testdata/tmpfile")
	u1 := file.ModTime("testdata/tmpfile")
	log.Print(u1)

	// touch it and capture the new time
	file.Touch("testdata/tmpfile")
	u2 := file.ModTime("testdata/tmpfile")
	log.Print(u2)

	// check that they are not equiv
	fmt.Println(u1 == u2)

	// Output:
	// false
}
