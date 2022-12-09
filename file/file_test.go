package file_test

import (
	"fmt"
	"log"
	"net/http"
	ht "net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fs"
	"github.com/rwxrob/fs/file"
)

func ExampleTouch_create() {
	fmt.Println(fs.NotExists("testdata/foo"))
	file.Touch("testdata/foo")
	fmt.Println(fs.Exists("testdata/foo"))
	os.Remove("testdata/foo")
	// Output:
	// true
	// true
}

func ExampleTouch_update() {

	// first create it and capture the time as a string
	file.Touch("testdata/tmpfile")
	u1 := fs.ModTime("testdata/tmpfile")
	log.Print(u1)

	// touch it and capture the new time
	file.Touch("testdata/tmpfile")
	u2 := fs.ModTime("testdata/tmpfile")
	log.Print(u2)

	// check that they are not equiv
	fmt.Println(u1 == u2)

	// Output:
	// false
}

func ExampleFetch() {

	// serve get
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `random file content`)
		})
	svr := ht.NewServer(handler)
	defer svr.Close()
	defer os.Remove(`testdata/file`)

	// not found
	handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		})
	notfound := ht.NewServer(handler)
	defer notfound.Close()

	if err := file.Fetch(svr.URL, `testdata/file`); err != nil {
		fmt.Println(err)
	}

	it, _ := os.ReadFile(`testdata/file`)
	fmt.Println(string(it))

	if err := file.Fetch(notfound.URL, `testdata/file`); err != nil {
		fmt.Println(err)
	}

	// Output:
	// random file content
	// 404 Not Found
}

func ExampleReplace() {

	// serve get
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `something random`)
	})
	svr := ht.NewServer(handler)
	defer svr.Close()

	// create a file to replace
	os.Create(`testdata/replaceme`)
	defer os.Remove(`testdata/replaceme`)
	os.Chmod(`testdata/replaceme`, 0400)

	// show info about control file
	info, err := os.Stat(`testdata/replaceme`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info.Mode())
	fmt.Println(info.Size())

	// replace it with local url
	file.Replace(`testdata/replaceme`, svr.URL)

	// check that it is new
	info, err = os.Stat(`testdata/replaceme`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info.Mode())
	fmt.Println(info.Size())

	// Output:
	// -r--------
	// 0
	// -r--------
	// 16

}

func ExampleExists() {
	fmt.Println(file.Exists("testdata/exists"))
	fmt.Println(file.Exists("testdata"))
	// Output:
	// true
	// false
}

func ExampleHereOrAbove_here() {
	dir, _ := os.Getwd()
	defer func() { os.Chdir(dir) }()
	os.Chdir("testdata/adir")

	path, err := file.HereOrAbove("afile")
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

	path, err := file.HereOrAbove("anotherfile")
	if err != nil {
		fmt.Println(err)
	}
	d := strings.Split(path, string(filepath.Separator))
	fmt.Println(d[len(d)-2:])

	// Output:
	// [testdata anotherfile]

}

func ExampleHead() {

	lines, err := file.Head(`testdata/headtail`, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [one two]
}

func ExampleHead_over() {

	lines, err := file.Head(`testdata/headtail`, 20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [one two three four five]
}

func ExampleTail() {

	lines, err := file.Tail(`testdata/headtail`, 2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [four five]
}

func ExampleTail_over() {

	lines, err := file.Tail(`testdata/headtail`, 20)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [one two three four five]
}

func ExampleTail_negative() {

	lines, err := file.Tail(`testdata/headtail`, -2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lines)

	// Output:
	// [three four five]
}

/*
func ExampleRepaceAllString() {
	err := file.ReplaceAllString(`testdata/headtail`, `three`, `THREE`)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// ignored
}
*/

func ExampleFindString() {
	str, err := file.FindString(`testdata/headtail`, `thre+`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)
	// Output:
	// three
}

func ExampleOverwrite() {
	err := file.Overwrite(`testdata/overwritten`, `hello`)
	defer os.Remove(`testdata/overwritten`)
	if err != nil {
		fmt.Println(err)
	}
	file.Cat(`testdata/overwritten`)
	// Output:
	// hello
}

func ExampleIsEmpty() {
	fmt.Println(file.IsEmpty(`testdata/overwritten`))
	fmt.Println(file.IsEmpty(`testdata/ovewritten`))
	file.Touch(`testdata/emptyfile`)
	fmt.Println(file.IsEmpty(`testdata/emptyfile`))
	// Output:
	// false
	// false
	// true
}

func ExampleSize() {
	fmt.Println(file.Size(`testdata/headtail`))
	// Output:
	// 24
}
