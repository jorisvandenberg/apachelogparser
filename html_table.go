package main

import (
	"html/template"
	"os"
)

type Table struct {
	Pagetitle       string
	Pagedescription string
	Pagecontent     []string
	Pagefooter      string
	Headers         map[string]string
	Data            []map[string]string
}

func createtable(args Args, outputfilename string, htmltitle string, myTable Table, section string, order int) {
	logger("creating a table with name " + outputfilename)
	t, err := template.New("mytemplate").Parse(templatedb["table_tmpl"])
	if err != nil {
		panic(err)
	}
	var outputHTMLFile *os.File
	if outputHTMLFile, err = os.Create(args.Outputs.Outputpath + outputfilename); err != nil {
		panic(err)
	}

	if err = t.Execute(outputHTMLFile, myTable); err != nil {
		panic(err)
	}
	defer outputHTMLFile.Close()

	MyPageForIndex := page_forindex{
		Title:   htmltitle,
		Url:     outputfilename,
		Section: section,
		Order:   order,
	}
	indexpages = append(indexpages, MyPageForIndex)
	logger("finished creating a table with name " + outputfilename)
}
