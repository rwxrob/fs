package fs_test

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fs"
)

//go:embed all:testdata/testfs
var testfs embed.FS

func ExampleExtractEmbed_confirm_Default_Read() {

	// go:embed all:testdata/testfs
	// var testfs embed.FS

	foo, err := testfs.Open("testdata/testfs/foo")
	if err != nil {
		fmt.Println(err)
	}
	info, err := foo.Stat()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info.Mode())

	// Output:
	// -r--r--r--
}

func ExampleExtractEmbed() {

	// go:embed all:testdata/testfs
	// var testfs embed.FS
	defer os.RemoveAll("testdata/testfsout")

	if err := fs.ExtractEmbed(testfs,
		"testdata/testfs", "testdata/testfsout"); err != nil {
		fmt.Println(err)
	}

	stuff := []string{
		`foo`,
		`_notignored`,
		`.secret`,
		`dir`,
		`dir/README.md`,
	}

	for _, i := range stuff {
		f, err := os.Stat(filepath.Join("testdata/testfsout", i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("%v %v\n", i, f.Mode())
	}

	// Output:
	// foo -rw-------
	// _notignored -rw-------
	// .secret -rw-------
	// dir drwx------
	// dir/README.md -rw-------

}

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
