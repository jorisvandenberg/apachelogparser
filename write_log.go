package main


func writelog(args args) {
	//var mylog []string
	if (args.writelog) {
		var newpage HtmlPage
		newpage.Pagetitle = "logggggs"
		newpage.Pagedescription = "here are the logs :)"
		newpage.Paragraphs = mylog
		createhtmltable(args, "logs.html",  newpage)
	}
}