package main

import (
	"github.com/go-echarts/go-echarts/v2/components"
	"io"
	"os"
)

func writedemographs(args Args) {
	logger("i was asked to write some demo pages")
	demobarchart(args)
	demotable(args)
	demolinegraph(args)
	demoboxplot(args)
	demopiechart(args)
	demowritemulti(args)
	demowritehtmlpage(args)
	debugstatfunct(args)
	logger("finished writing some demo pages")
}

func debugstatfunct(args Args) {
	parameters := []interface{}{"0"}
	tableheaders := map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_5": "NB RAW HITS",
	}
	sqlreturnvalues := []interface{}{
		[]interface{}{"year", "int"},
		[]interface{}{"month", "int"},
		[]interface{}{"day", "int"},
		[]interface{}{"hour", "int"},
		[]interface{}{"nm raw hits", "int"},
	}
	genstats(args, "demotext for log", "stat_perday_hits_unique_2xx_3xx", "stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed", parameters, tableheaders, sqlreturnvalues)
}

func demowritehtmlpage(args Args) {
	var newpage HtmlPage
	newpage.Pagetitle = "my demo page"
	newpage.Pagedescription = "this is a simple demo page"
	newpage.Paragraphs = append(newpage.Paragraphs, "this is the first paragraph")
	newpage.Paragraphs = append(newpage.Paragraphs, "and this is the second paragraph")
	createhtmltable(args, "demosimplepage.html", newpage, "demo", 99)
}

func demowritemulti(args Args) {
	XValues := []string{"Januari", "Februari", "March", "April"}
	YValues := make(map[string]int)
	YValues["a"] = 5
	YValues["b"] = 9
	YValues["c"] = 15
	page := components.NewPage()
	page.AddCharts(
		createpiechart(XValues, YValues, "this is a demo", "yup, a demo", args, "magniet1.html", "demos", 999, false),
		createpiechart(XValues, YValues, "this is a demo", "yup, a demo", args, "magniet2.html", "demos", 999, false),
	)
	f, err := os.Create("./output/piechartz.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
	MyPageForIndex := page_forindex{
		Title:   "blahblah",
		Url:     "piechartz.html",
		Section: "demo",
		Order:   99,
	}
	indexpages = append(indexpages, MyPageForIndex)
}

func demopiechart(args Args) {
	XValues := []string{"Januari", "Februari", "March", "April"}
	YValues := make(map[string]int)
	YValues["a"] = 5
	YValues["b"] = 9
	YValues["c"] = 15
	createpiechart(XValues, YValues, "this is a demo", "yup, a demo", args, "demopiegraph.html", "demos", 99)
}

func demolinegraph(args Args) {
	XValues := []string{"Apple", "Banana", "Peach ", "Lemon"}
	YValues := make(map[string][]int)
	YValues["This year"] = append(YValues["This year"], 5)
	YValues["This year"] = append(YValues["This year"], 15)
	YValues["This year"] = append(YValues["This year"], 7)
	YValues["This year"] = append(YValues["This year"], 9)
	YValues["Last year"] = append(YValues["Last year"], 19)
	YValues["Last year"] = append(YValues["Last year"], 25)
	YValues["Last year"] = append(YValues["Last year"], 17)
	YValues["Last year"] = append(YValues["Last year"], 14)
	PreChartText = "first chart"
	PostChartText = "end first chart"
	createlinegraph(XValues, YValues, "this is a demo", "yup, a demo", args, "1demolinegraph.html", "demos", 99)
	PreChartText = "second chart"
	PostChartText = "end second chart"
	createlinegraph(XValues, YValues, "this is a demo", "yup, a demo", args, "2demolinegraph.html", "demos", 99)
	PreChartText = ""
	PostChartText = ""
}

func demobarchart(args Args) {
	XValues := []string{"Januari", "Februari", "March", "April"}
	YValues := make(map[string][]int)
	YValues["This year"] = append(YValues["This year"], 5)
	YValues["This year"] = append(YValues["This year"], 15)
	YValues["This year"] = append(YValues["This year"], 7)
	YValues["This year"] = append(YValues["This year"], 9)
	YValues["Last year"] = append(YValues["Last year"], 19)
	YValues["Last year"] = append(YValues["Last year"], 25)
	YValues["Last year"] = append(YValues["Last year"], 17)
	YValues["Last year"] = append(YValues["Last year"], 14)
	createbarchart(XValues, YValues, "this is a demo", "yup, a demo", args, "demobarchart.html", "demos", 99)
}

func demoboxplot(args Args) {
	XValues := []string{"Januari", "Februari", "March", "April"}
	YValues := make(map[string][][]int)
	YValues["This year"] = append(YValues["This year"], []int{7, 14, 8, 9, 99, 44})
	YValues["This year"] = append(YValues["This year"], []int{8, 12, 7, 8, 94, 45})
	YValues["This year"] = append(YValues["This year"], []int{8, 12, 7, 8, 94, 45})
	YValues["This year"] = append(YValues["This year"], []int{8, 12, 7, 8, 94, 45})

	YValues["Last year"] = append(YValues["Last year"], []int{7, 14, 8, 9, 99, 44})
	YValues["Last year"] = append(YValues["Last year"], []int{8, 12, 7, 8, 94, 45})
	YValues["Last year"] = append(YValues["Last year"], []int{8, 12, 7, 8, 94, 45})
	YValues["Last year"] = append(YValues["Last year"], []int{8, 12, 7, 8, 94, 45})

	createboxplot(XValues, YValues, args, "my demo boxplot", "demoboxplot.html", "demos", 999)
}
func demotable(args Args) {
	MyHeaders := map[string]string{
		"Title_1": "kolom 1",
		"Title_2": "kolom 2",
	}

	myTable := Table{
		Pagetitle:       "tadaa",
		Pagedescription: "blahblah",
		Pagefooter:      "the end",
		Pagecontent:     []string{"one", "this istwo", "three"},
		Headers:         MyHeaders,
		Data:            []map[string]string{},
	}

	MyData := map[string]string{
		"Value_1": "record 1.1",
		"Value_2": "record 1.2",
	}
	myTable.Data = append(myTable.Data, MyData)
	MyData = map[string]string{
		"Value_1": "record 2.1",
		"Value_2": "record 2.2",
	}
	myTable.Data = append(myTable.Data, MyData)

	createtable(args, "demotable.html", "this is a title", myTable, "demos", 999)
}
