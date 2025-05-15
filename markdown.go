package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"text/template"
)

var tmpl = template.Must(template.ParseFiles("template.html"))

// [_metadata_:group1]:# "group2"
var rMetadata = regexp.MustCompile(`(^|[^ ]) {0,3}\[_metadata_:(\w+)\]:\n?[\r\t\f\v ]*\S+\n?[\r\t\f\v ]+"(.+)"`)

type Page struct {
	Content  string
	Index    []Index
	DirPath  string
	BaseName string
	Metadata map[string]string
}

type Index struct {
	IsDir bool
	Href  string
	Path  string
}

func NewPage() *Page {
	p := Page{
		Content:  "",
		DirPath:  "",
		BaseName: "",
		Index:    nil,
		Metadata: make(map[string]string)}
	p.Metadata["Title"] = "Manne.dev"
	return &p
}

func (p *Page) setContentHTML(html string) *Page {
	p.Content = html
	return p
}

func (p *Page) setContentMarkdown(md []byte) *Page {
	p.Content = string(markdown.ToHTML(md, nil, nil))
	return p
}

func (p *Page) setIndex(url string) *Page {
	dirPath, entries := getClosestIndex(url)

	p.DirPath = dirPath
	p.Index = make([]Index, len(entries))

	if dirPath != url {
		p.BaseName = path.Base(url)
	}

	for n, e := range entries {
		i := Index{
			IsDir: e.IsDir(),
			Href:  e.Name(),
			Path:  e.Name()}
		if i.IsDir {
			i.Href += "/"
		}
		if path.Ext(e.Name()) == ".md" {
			i.Href = i.Href[:len(i.Href)-3]
			i.Path = i.Href
		}
		p.Index[n] = i
	}
	return p
}

func (p *Page) Send(w http.ResponseWriter) {
	metadatatags := rMetadata.FindAllStringSubmatch(p.Content, -1)
	for _, tag := range metadatatags {
		p.Metadata[tag[2]] = tag[3]
	}

	tmpl.Execute(w, p)
}

func sendMarkdownFile(w http.ResponseWriter, path string) error {
	md, err := getMarkdown(path)
	if err != nil {
		return err
	}
	sendMarkdown(w, md)
	return nil
}

func sendMarkdown(w http.ResponseWriter, md []byte) {
	NewPage().setContentMarkdown(md).Send(w)
}

func getMarkdown(filename string) ([]byte, error) {
	return getStaticFile(filename + ".md")
}

func getClosestIndex(dirPath string) (string, []fs.DirEntry) {
	var entries []fs.DirEntry
	var err error
	for {
		dirPath = path.Dir(dirPath)
		entries, err = fs.ReadDir(FS, path.Clean("./"+dirPath))

		if err == nil {
			break
		}
	}

	if string(dirPath[len(dirPath)-1]) != "/" {
		dirPath += "/"
	}

	return dirPath, entries
}

func getMarkdownIndex(path string, entries []fs.DirEntry) string {
	md := "# Index of " + path + "\n\n"
	if path != "/" {
		md += fmt.Sprintf("* [%[2]s](%[1]s%[2]s)\n", path, "..")
	}

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
