package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func createpiechart(XValues []string, YValues map[string]int, title string, subtitle string, args args, filename string, section string, order int, writehtml ...bool) *charts.Pie {
	logger("creating a piechart with title '" + title+ "' and filename " + filename)
	pie := charts.NewPie()
	pie.Renderer = newSnippetRenderer(pie, pie.Validate)
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: title}),
	)

	items := make([]opts.PieData, 0)
	for serienaam, serievalue := range YValues {
		items = append(items, opts.PieData{Name: serienaam, Value: serievalue})
	}
	pie.AddSeries("pie", items).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
		)
	writehtml_optional := true
	if len(writehtml) > 0 {
		writehtml_optional = writehtml[0]
	}
	if writehtml_optional {
		f, _ := os.Create(args.outputpath + filename)
		_ = pie.Render(f)

		MyPageForIndex := page_forindex{
			Title:   title,
			Url:     filename,
			Section: section,
			Order:   order,
		}
		indexpages = append(indexpages, MyPageForIndex)
	}
	logger("finished creating a piechart with title '" + title+ "' and filename " + filename)
	return pie

}
