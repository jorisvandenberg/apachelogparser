package main

import (
	"fmt"
)

func writelog(args Args) {
	if args.Generals.Writelog {
		var newpage HtmlPage
		newpage.Pagetitle = "detailed runlogs of the statistics application"
		newpage.Pagedescription = "here are the logs :)"
		newpage.Paragraphs = mylog
		createhtmltable(args, "logs.html", newpage, "9. logs", 999)
	}
	if args.Commandlines.Debug {
		for _, printdebug := range mylog {
			fmt.Printf("DEBUG => %s\n", printdebug)
		}
	}
}
