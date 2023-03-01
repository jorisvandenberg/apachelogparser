package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)


func createpiechart(XValues []string, YValues map[string]int, title string, subtitle string, args args, filename string) {
	pie := charts.NewPie()
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
		f, _ := os.Create(args.outputpath + filename)
		_ = pie.Render(f)

		MyPageForIndex := page_forindex{
			Title: title,
			Url:   filename,
		}
		indexpages = append(indexpages, MyPageForIndex)

}