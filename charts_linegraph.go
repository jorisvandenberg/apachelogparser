package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func createlinegraph(XValues []string, YValues map[string][]int, title string, subtitle string, args args, filename string, section string, order int, writehtml ...bool) *charts.Line {
	logger("creating a linegraph with title '" + title+ "' and filename " + filename)
	line := charts.NewLine()
	line.Renderer = newSnippetRenderer(line, line.Validate)
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
			//Link:     "https://github.com/go-echarts/go-echarts",
		}),
	)
	for serienaam, serievalues := range YValues {
		items := make([]opts.LineData, 0)
		for _, serievalue := range serievalues {
			items = append(items, opts.LineData{Value: serievalue})
		}
		line.SetXAxis(XValues).AddSeries(serienaam, items).SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: true}))
	}

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle:  title,
			Width:      `95vw`,
			Height:     `95vh`,
			AssetsHost: args.assethost,
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:  true,
			Top:   "10%",
			Right: "5%",
		}),
	)

	writehtml_optional := true
	if len(writehtml) > 0 {
		writehtml_optional = writehtml[0]
	}
	if writehtml_optional {
		f, _ := os.Create(args.outputs.outputpath + filename)
		_ = line.Render(f)
		MyPageForIndex := page_forindex{
			Title:   title,
			Url:     filename,
			Section: section,
			Order:   order,
		}
		indexpages = append(indexpages, MyPageForIndex)
	}
	logger("finished creating a linegraph with title '" + title+ "' and filename " + filename)
	return line

}
