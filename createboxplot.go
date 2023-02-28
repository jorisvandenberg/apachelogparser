package main

import (
	"os"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)



func createboxplot(XValues []string, YValues map[string][][]int, args args, title string, filename string)  {
	bp := charts.NewBoxPlot()
	bp.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
	)

	for serienaam, serievalues := range YValues {
		items := make([]opts.BoxPlotData, 0)
		for _, serievalue := range serievalues {
			items = append(items, opts.BoxPlotData{Value: serievalue})
		}
		bp.SetXAxis(XValues).AddSeries(serienaam, items)
	
		f, _ := os.Create(args.outputpath + filename)
		_ = bp.Render(f)
		
}
MyPageForIndex := page_forindex{
	Title: title,
	Url:   filename,
}
indexpages = append(indexpages, MyPageForIndex)
}
