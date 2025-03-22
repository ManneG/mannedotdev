package main

import (
	"net/http"
	"regexp"
	"text/template"
	"github.com/gomarkdown/markdown"
)

var tmpl = template.Must(template.ParseFiles("template.html"))
// [_metadata_:group1]:# "group2"
var rMetadata = regexp.MustCompile(`(^|[^ ]) {0,3}\[_metadata_:(\w+)\]:\n?[\r\t\f\v ]*\S+\n?[\r\t\f\v ]+"(.+)"`)


type Page struct {
	Content string
	Metadata map[string]string
}

func NewPage(content string) Page {
	p := Page{
		Content: content,
		Metadata: make(map[string]string)}
	p.Metadata["Title"] = "Manne.dev"
	return p
}

func sendMarkdown(w http.ResponseWriter, path string) error {
	md, err := getMarkdown(path)
	if err != nil {
		return err
	}
	html := markdown.ToHTML(md, nil, nil)

	p := NewPage(string(html))

	metadatatags := rMetadata.FindAllStringSubmatch(string(md), -1)
	for _, tag := range metadatatags {
		p.Metadata[tag[2]] = tag[3]
	}

	tmpl.Execute(w, p)

	return nil
}

func getMarkdown(filename string) ([]byte, error) {
	return getStaticFile(filename + ".md")
}