package fs_test

import (
	"fmt"

	"github.com/rwxrob/fs"
)

func ExampleIsDir() {
	fmt.Println(fs.IsDir("testdata"))
	fmt.Println(fs.IsDir("testdata/fail"))
	// Output:
	// true
	// false
}
