package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func createbarchart(XValues []string, YValues map[string][]int, title string, subtitle string, args args, filename string) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    title,
		Subtitle: subtitle,
		Top:      "5%",
	}))
	for serienaam, serievalues := range YValues {
		items := make([]opts.BarData, 0)
		for _, serievalue := range serievalues {
			items = append(items, opts.BarData{Value: serievalue})
		}
		bar.SetXAxis(XValues).AddSeries(serienaam, items)
	}
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: title,
			Width: `95vw`,
			Height: `95vh`, 
			AssetsHost: args.assethost,
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:  true,
			Top:   "10%",
			Right: "5%",
		}),
	)
	f, _ := os.Create(args.outputpath + filename)
	_ = bar.Render(f)

	MyPageForIndex := page_forindex{
		Title: title,
		Url:   filename,
	}
	indexpages = append(indexpages, MyPageForIndex)
}