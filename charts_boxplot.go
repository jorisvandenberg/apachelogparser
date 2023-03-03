package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func createboxplot(XValues []string, YValues map[string][][]int, args args, title string, filename string, writehtml ...bool) *charts.BoxPlot {
	bp := charts.NewBoxPlot()
	bp.Renderer = newSnippetRenderer(bp, bp.Validate)
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
	writehtml_optional := true
	if len(writehtml) > 0 {
		writehtml_optional = writehtml[0]
	}
	if writehtml_optional {
		f, _ := os.Create(args.outputpath + filename)
		_ = bp.Render(f)
		MyPageForIndex := page_forindex{
			Title: title,
			Url:   filename,
		}
		indexpages = append(indexpages, MyPageForIndex)
	}
	return bp
}
