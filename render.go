package main

import (
	"os"
	"strings"
)

var htmlTemplate string

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