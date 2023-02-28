package main

import (
	"html/template"
	"os"
)

type page_forindex struct {
	Title    string
	Url      string
	Textpre  string
	Textpost string
}

var indexpages []page_forindex

func createindex(args args) {
	t, err := template.New("mytemplate").Parse(templatedb["html_index"])
	if err != nil {
		panic(err)
	}
	var outputHTMLFile *os.File
	if outputHTMLFile, err = os.Create(args.outputpath + "index.html"); err != nil {
		panic(err)
	}

	if err = t.Execute(outputHTMLFile, indexpages); err != nil {
		panic(err)
	}
	defer outputHTMLFile.Close()
}