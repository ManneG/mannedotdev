package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"github.com/gomarkdown/markdown"
)

func sendMarkdown(w http.ResponseWriter, path string) error {
	md, err := getMarkdown(path)
	if err != nil {
		return err
	}
	html := markdown.ToHTML(md, nil, nil)
	out, err := render(string(html), nil)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, out)
	return nil
}

func getMarkdown(filename string) ([]byte, error) {
	FS := os.DirFS("./static")
	return fs.ReadFile(FS, filename[1:] + ".md")
}