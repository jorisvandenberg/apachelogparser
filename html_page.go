package main

import (
	"html/template"
	"os"
)

type HtmlPage struct {
	Pagetitle       string
	Pagedescription string
	Paragraphs      []string
}

func createhtmltable(args Args, outputfilename string, myHtmlPage HtmlPage, section string, order int) {
	logger("creating a html file with name " + outputfilename)
	t, err := template.New("mytemplate").Parse(templatedb["html_page"])
	if err != nil {
		panic(err)
	}
	var outputHTMLFile *os.File
	if outputHTMLFile, err = os.Create(args.Outputs.Outputpath + outputfilename); err != nil {
		panic(err)
	}

	if err = t.Execute(outputHTMLFile, myHtmlPage); err != nil {
		panic(err)
	}
	defer outputHTMLFile.Close()

	MyPageForIndex := page_forindex{
		Title:   myHtmlPage.Pagetitle,
		Url:     outputfilename,
		Section: section,
		Order:   order,
	}
	indexpages = append(indexpages, MyPageForIndex)
	logger("finished creating a html file with name " + outputfilename)
}
