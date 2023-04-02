package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func createboxplot(XValues []string, YValues map[string][][]int, args Args, title string, filename string, section string, order int, writehtml ...bool) *charts.BoxPlot {
	logger("creating a boxplot with title '" + title + "' and filename " + filename)
	bp := charts.NewBoxPlot()
	bp.Renderer = newSnippetRenderer(bp, bp.Validate)
	bp.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
		charts.WithInitializationOpts(opts.Initialization{
			AssetsHost: args.Outputs.Assethost,
		}),
	)

	for serienaam, serievalues := range YValues {
		items := make([]opts.BoxPlotData, 0)
		for _, serievalue := range serievalues {
			items = append(items, opts.BoxPlotData{Value: serievalue})
		}
		bp.SetXAxis(XValues).AddSeries(serienaam, items)

		f, _ := os.Create(args.Outputs.Outputpath + filename)
		_ = bp.Render(f)

	}
	writehtml_optional := true
	if len(writehtml) > 0 {
		writehtml_optional = writehtml[0]
	}
	if writehtml_optional {
		f, _ := os.Create(args.Outputs.Outputpath + filename)
		_ = bp.Render(f)
		MyPageForIndex := page_forindex{
			Title:   title,
			Url:     filename,
			Section: section,
			Order:   order,
		}
		indexpages = append(indexpages, MyPageForIndex)
	}
	logger("finished creating a boxplot with title '" + title + "' and filename " + filename)
	return bp
}
