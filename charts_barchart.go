package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
)

func createbarchart(XValues []string, YValues map[string][]int, title string, subtitle string, args args, filename string, writehtml ...bool) *charts.Bar {
	bar := charts.NewBar()
	bar.Renderer = newSnippetRenderer(bar, bar.Validate)
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
		//https://github.com/go-echarts/go-echarts/blob/ad9b214d3d71d5a1a0e02c2a706df9b23acdcbf6/charts/base.go#L21
		//https://github.com/go-echarts/go-echarts/blob/1dee3e5ca83599ebae7d62f4afb3e1113ebda667/opts/global.go#L20
		//https://github.com/go-echarts/go-echarts/blob/1dee3e5ca83599ebae7d62f4afb3e1113ebda667/opts/global.go#L157
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
		f, _ := os.Create(args.outputpath + filename)
		_ = bar.Render(f)
		MyPageForIndex := page_forindex{
			Title:   title,
			Url:     filename,
			Section: "graphs",
			Order:   3,
		}
		indexpages = append(indexpages, MyPageForIndex)
	}
	return bar
}
