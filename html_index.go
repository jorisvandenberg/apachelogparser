package main

import (
	"html/template"
	"os"
	"sort"
)

type page_forindex struct {
	Title    string
	Section	 string
	Order 	 int
	Url      string
	Textpre  string
	Textpost string
}

var indexpages []page_forindex

func createindex(args args) {
	sort.Slice(indexpages, func(i, j int) bool {
		return indexpages[i].Order < indexpages[j].Order
	  })

	groups := make(map[string][]page_forindex)
	  for _, p := range indexpages {
		  groups[p.Section] = append(groups[p.Section], p)
	  }

	t, err := template.New("mytemplate").Parse(templatedb["html_index"])
	if err != nil {
		panic(err)
	}
	var outputHTMLFile *os.File
	if outputHTMLFile, err = os.Create(args.outputpath + "index.html"); err != nil {
		panic(err)
	}

	if err = t.Execute(outputHTMLFile, groups); err != nil {
		panic(err)
	}
	defer outputHTMLFile.Close()
}
