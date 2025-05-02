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
		err := sendMarkdownFile(w, "/index")
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
			mType := mime.TypeByExtension(path.Ext(r.URL.Path))
			w.Header().Set("Content-Type", mType)
			if mType == "" {
				fmt.Println("No matching MIME type for " + path.Ext(r.URL.Path))
			}
			w.Write(data)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		err = sendMarkdownFile(w, r.URL.Path)
		if err == nil {
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
	dirPath := r.URL.Path
	var direntry []fs.DirEntry
	var err error

	for {
		dirPath = path.Dir(dirPath)	
		direntry, err = fs.ReadDir(FS, path.Clean("./" + dirPath))
		
		if(err == nil) {
			break;
		}
	}
	
	if string(dirPath[len(dirPath)-1]) != "/" {
		dirPath += "/"
	}
	isExistingDirectory := dirPath == r.URL.Path
	
	if isExistingDirectory {
		md := getMarkdownIndex(dirPath, direntry)
		sendMarkdown(w, []byte(md))
		return
	}

	fmt.Println("404: " + r.URL.Path)
	//http.NotFound(w, r)
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
	return fs.ReadFile(FS, filename[1:])
}

func getMarkdownIndex(path string, entries []fs.DirEntry) string {
	md := "# Index of " + path + "\n\n"
	md += fmt.Sprintf("* [%[2]s](%[1]s%[2]s)\n", path, "..")

	for _, e := range entries {
		relPath := e.Name()
		if e.IsDir() {
			relPath += "/"
		}
		if relPath[len(relPath)-3:] == ".md" {
			relPath = relPath[:len(relPath)-3]
		} else {
		}
		md += fmt.Sprintf("* [%[2]s](%[1]s%[2]s)\n", path, relPath)
	}

	return md + "\n"
}