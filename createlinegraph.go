package main

import (
	"os"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)




func createlinegraph(XValues []string, YValues map[string][]int, title string, subtitle string, args args, filename string) {
		line := charts.NewLine()
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
			line.SetXAxis(XValues).AddSeries(serienaam, items).SetSeriesOptions(charts.WithLabelOpts(opts.Label{Show: true,}),)
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
		
	f, _ := os.Create(args.outputpath + filename)
			_ = line.Render(f)
}