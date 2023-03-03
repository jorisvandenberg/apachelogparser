package main

import (
	"fmt"
	chartrender "github.com/go-echarts/go-echarts/v2/render"
	"html/template"
	"io"
)

var PreChartText string
var PostChartText string

type Renderer interface {
	Render(w io.Writer) error
}

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
			Parse(PreChartText + templatedb["baseTpl"] + PostChartText),
		)

	err := tpl.ExecuteTemplate(w, tplName, r.c)
	return err
}
