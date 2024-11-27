package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		http.ServeFile(w, r, "./index.html")
	})

	http.HandleFunc("GET /info", func(w http.ResponseWriter, r *http.Request) {
		r.Header["X-Forwarded-For"] = append(r.Header["X-Forwarded-For"], r.RemoteAddr)
		fmt.Println(r.URL.Path)
		fmt.Fprintf(w, "You are %s visiting %s at %s using %s.", r.Header["X-Forwarded-For"][0], r.Host, r.URL.Path, r.Proto)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("404: " + r.URL.Path)
		http.NotFound(w, r)
	})

	log.Print("Listening on :8080")
	path, _ := os.Getwd()
	log.Print("PWD: ", path)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

