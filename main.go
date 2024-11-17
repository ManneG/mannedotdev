package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	http.HandleFunc("GET /info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You are %s visiting %s using %s.", r.RemoteAddr, r.Host, r.Proto)
	})

	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
