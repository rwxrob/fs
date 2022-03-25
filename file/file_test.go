package file_test

import (
	"fmt"
	"log"
	"net/http"
	ht "net/http/httptest"
	"os"

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
