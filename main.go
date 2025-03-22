package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var htmlTemplate string

func main() {
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

	http.HandleFunc("GET /blog", func(w http.ResponseWriter, r *http.Request) {
		err := sendMarkdown(w, r.URL.Path)
		if err != nil {
			fmt.Println(err)
			notFound(w, r)
		}
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

type renderOptions struct {
	title string
}

func render(content string, options *renderOptions) (string, error) {
	if htmlTemplate == "" {
		bytes, err := os.ReadFile("template.html")
		if err != nil {
			return "", err
		}
		htmlTemplate = string(bytes)
	}

	title := "Manne.dev"
	if options != nil {
		title = options.title
	}

	temp := strings.Replace(htmlTemplate, "<content>", content, 1)
	temp = strings.Replace(temp, "<titlecontent>", title, 1)

	return temp, nil
}