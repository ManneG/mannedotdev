package main

import (
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
)

var FS = os.DirFS("./static")

func main() {
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		md, err := getMarkdown("/index")
		if err != nil {
			fmt.Println("index errored")
			return
		}

		p := NewPage()
		p.setContentMarkdown(md)
		p.setIndex("/index")
		p.Send(w)
	})

	http.HandleFunc("GET /info", func(w http.ResponseWriter, r *http.Request) {
		r.Header["X-Forwarded-For"] = append(r.Header["X-Forwarded-For"], r.RemoteAddr)
		fmt.Println(r.URL.Path)
		fmt.Fprintf(w, "You are %s visiting %s at %s using %s.", r.Header["X-Forwarded-For"][0], r.Host, r.URL.Path, r.Proto)
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		data, err := getStaticFile(r.URL.Path)
		if err == nil {
			mType := mime.TypeByExtension(path.Ext(r.URL.Path))
			w.Header().Set("Content-Type", mType)
			if mType == "" {
				fmt.Println("No matching MIME type for " + path.Ext(r.URL.Path))
			}
			w.Write(data)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		md, err := getMarkdown(r.URL.Path)
		if err == nil {
			NewPage().setContentMarkdown(md).setIndex(r.URL.Path).Send(w)
			return
		}

		notFound(w, r)
	})

	http.HandleFunc("/", notFound)

	log.Print("Listening on :8080")
	path, _ := os.Getwd()
	log.Print("PWD: ", path)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	dirPath, direntry := getClosestIndex(r.URL.Path)

	fileInfo, err := fs.Stat(FS, path.Clean(r.URL.Path[1:]))
	isExistingDirectory := err == nil && fileInfo.IsDir()

	if isExistingDirectory {
		md := getMarkdownIndex(dirPath, direntry)
		NewPage().setContentMarkdown([]byte(md)).setIndex(r.URL.Path).Send(w)
		return
	}

	fmt.Println("404: " + r.URL.Path)

	headers := w.Header()
	headers.Del("Content-Length")
	headers.Set("Content-Type", "text/html; charset=utf-8")
	headers.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(404)

	md := "# Error 404: Not Found\n\n"
	md += "This file or directory does not exist\n\n#"
	md += getMarkdownIndex(dirPath, direntry)
	md += "[_metadata_:Title]:# \"404 Not Found\""
	sendMarkdown(w, []byte(md))
}

func getStaticFile(filename string) ([]byte, error) {
	return fs.ReadFile(FS, path.Clean("./"+filename))
}
