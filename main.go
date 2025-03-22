package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var FS = os.DirFS("./static")

func main() {
	/* This could be used to give the user a nice index of the site
	
	direntry, _ := fs.ReadDir(FS, ".")
	
	for _, i := range direntry {
		fmt.Println(i)
	}*/

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		err := sendMarkdown(w, "/index")
		if err != nil {
			fmt.Println(err)
			notFound(w, r)
		}
	})

	http.HandleFunc("GET /info", func(w http.ResponseWriter, r *http.Request) {
		r.Header["X-Forwarded-For"] = append(r.Header["X-Forwarded-For"], r.RemoteAddr)
		fmt.Println(r.URL.Path)
		fmt.Fprintf(w, "You are %s visiting %s at %s using %s.", r.Header["X-Forwarded-For"][0], r.Host, r.URL.Path, r.Proto)
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		data, err := getStaticFile(r.URL.Path)
		if err == nil {
			w.Write(data)
			return
		}

		err = sendMarkdown(w, r.URL.Path)
		if err == nil {
			return
		}

		fmt.Println(err)
		notFound(w, r)
	})

	http.HandleFunc("/", notFound)

	log.Print("Listening on :8080")
	path, _ := os.Getwd()
	log.Print("PWD: ", path)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println("404: " + r.URL.Path)
	http.NotFound(w, r)
}

func getStaticFile(filename string) ([]byte, error) {
	return fs.ReadFile(FS, filename[1:])
}