package main

import (
	"html/template"
	"os"
	"sort"
)

type page_forindex struct {
	Title    string
	Section  string
	Order    int
	Url      string
	Textpre  string
	Textpost string
}

var indexpages []page_forindex

func removeDuplicates(pages []page_forindex) []page_forindex {
    encountered := map[page_forindex]bool{} 
    result := []page_forindex{} 

    for _, page := range pages {
        if encountered[page] == false { 
            encountered[page] = true 
            result = append(result, page) 
        }
    }

    return result // return the new slice with unique structs
}

func createindex(args Args) {
	logger("creating an index file")
	indexpages = removeDuplicates(indexpages)
	sort.Slice(indexpages, func(i, j int) bool {
		return indexpages[i].Order < indexpages[j].Order
	})

	groups := make(map[string][]page_forindex)
	for _, p := range indexpages {
		groups[p.Section] = append(groups[p.Section], p)
	}

	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	groups_sorted := make(map[string][]page_forindex)
	for _, k := range keys {
		groups_sorted[k] = groups[k]
	}

	t, err := template.New("mytemplate").Parse(templatedb["html_index"])
	if err != nil {
		panic(err)
	}
	var outputHTMLFile *os.File
	if outputHTMLFile, err = os.Create(args.Outputs.Outputpath + "index.html"); err != nil {
		panic(err)
	}

	if err = t.Execute(outputHTMLFile, groups_sorted); err != nil {
		panic(err)
	}
	defer outputHTMLFile.Close()
	logger("finished creating an index file")
}
