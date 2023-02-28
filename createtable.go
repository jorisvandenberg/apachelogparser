package main

import (
	"html/template"
	"os"
)

type Table struct {
	Pagetitle       string
	Pagedescription string
	Headers         map[string]string
	Data            []map[string]string
}

func createtable(args args, outputfilename string, htmltitle string, myTable Table) {
	t, err := template.New("mytemplate").Parse(templatedb["table_tmpl"])
	if err != nil {
		panic(err)
	}
	var outputHTMLFile *os.File
	if outputHTMLFile, err = os.Create(args.outputpath + outputfilename); err != nil {
		panic(err)
	}

	if err = t.Execute(outputHTMLFile, myTable); err != nil {
		panic(err)
	}
	defer outputHTMLFile.Close()

	MyPageForIndex := page_forindex{
		Title: htmltitle,
		Url:   outputfilename,
	}
	indexpages = append(indexpages, MyPageForIndex)
}