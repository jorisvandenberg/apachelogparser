package main

import (
	//"io"
	"math/rand"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	//"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	itemCntLine = 6
	fruits      = []string{"Apple", "Banana", "Peach ", "Lemon", "Pear", "Cherry"}
)

func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < itemCntLine; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func createlinegraph(args args) {
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title:    "title and label options",
				Subtitle: "go-echarts is an awesome chart library written in Golang",
				Link:     "https://github.com/go-echarts/go-echarts",
			}),
		)
		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				PageTitle:  "mijn pagina titel",
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
		line.SetXAxis(fruits).
			AddSeries("Category A", generateLineItems()).
			SetSeriesOptions(
				charts.WithLabelOpts(opts.Label{
					Show: true,
				}),
			)
		
	f, _ := os.Create("./output/linchaart.html")
			_ = line.Render(f)
}