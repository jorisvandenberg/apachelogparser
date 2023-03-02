package main

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"os"
	
)

import "fmt"
import "io"
import "html/template"
import	chartrender "github.com/go-echarts/go-echarts/v2/render"

type Renderer interface {
	Render(w io.Writer) error
}

var baseTpl = `
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="width:{{ .Initialization.Width }};height:{{ .Initialization.Height }};"></div>
</div>
{{- range .JSAssets.Values }}
   <script src="{{ . }}"></script>
{{- end }}
<script type="text/javascript">
    "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}");
    let option_{{ .ChartID | safeJS }} = {{ .JSON }};
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});
    {{- range .JSFunctions.Fns }}
    {{ . | safeJS }}
    {{- end }}
</script>blahblahblah
`


type snippetRenderer struct {
	c      interface{}
	before []func()
}

func newSnippetRenderer(c interface{}, before ...func()) chartrender.Renderer {
	return &snippetRenderer{c: c, before: before}
}

func (r *snippetRenderer) Render(w io.Writer) error {
	const tplName = "chart"
	for _, fn := range r.before {
		fn()
	}

	tpl := template.
		Must(template.New(tplName).
			Funcs(template.FuncMap{
				"safeJS": func(s interface{}) template.JS {
					return template.JS(fmt.Sprint(s))
				},
			}).
			Parse(baseTpl + "ikkedikkeslaggevallen"),
		)

	err := tpl.ExecuteTemplate(w, tplName, r.c)
	return err
}

func createlinegraph_html(XValues []string, YValues map[string][]int, title string, subtitle string, args args, filename string, writehtml ...bool) *charts.Line {
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
		f, _ := os.Create(args.outputpath + filename)
		_ = line.Render(f)
		MyPageForIndex := page_forindex{
			Title: title,
			Url:   filename,
		}
		indexpages = append(indexpages, MyPageForIndex)
	}
	return line

}


func createlinegraph(XValues []string, YValues map[string][]int, title string, subtitle string, args args, filename string, writehtml ...bool) *charts.Line {
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
		f, _ := os.Create(args.outputpath + filename)
		_ = line.Render(f)
		MyPageForIndex := page_forindex{
			Title: title,
			Url:   filename,
		}
		indexpages = append(indexpages, MyPageForIndex)
	}
	return line

}
