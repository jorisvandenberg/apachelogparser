package main

var templatedb = make(map[string]string)

func filltemplatedb() {
	baseTpl := `
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
</script>
`
	html_page := `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Pagetitle}}</title>
	</head>
	<body>
		<h1>{{.Pagetitle}}</h1>
		<p>{{.Pagedescription}}</p>
		<div>

				{{range .Paragraphs}}
				<p>{{.}}</p>
				{{end}}
		</div>

	</body>
</html>`
	table_tmpl := `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Pagetitle}}</title>
		<!-- choose a theme file -->
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.31.3/css/theme.default.min.css">
		<!-- load jQuery and tablesorter scripts -->
		<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.3/jquery.min.js"></script>
		<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.31.3/js/jquery.tablesorter.js"></script>
		
		<!-- tablesorter widgets (optional) -->
		<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/jquery.tablesorter/2.31.3/js/jquery.tablesorter.widgets.js"></script>
	</head>
	<body>
		<h1>{{.Pagetitle}}</h1>
		<p>{{.Pagedescription}}</p>
		<p>
		<table  id="myTable" class="tablesorter" border = "1">
		<thead>
			<tr>
				{{range .Headers}}
				<th>{{.}}</th>
				{{end}}
			</tr>
			</thead>
			<tbody>
		{{range .Data}}
			<tr>
				{{range .}}
				<td>{{.}}</td>
				{{end}}
			</tr>
		{{end}}
		</tbody>
		</table>
		</p>
		<script>
		$(function() {
			$("#myTable").tablesorter();
		  });
		</script>
	</body>
</html>`
	html_index := `<!DOCTYPE html>
<html>
	<body>
				{{range .}}
				<p>{{.Textpre}}<a href="{{.Url}}">{{.Title}}</a>{{.Textpost}}</p>
				{{end}}
	</body>
</html>`

	templatedb["table_tmpl"] = table_tmpl
	templatedb["html_index"] = html_index
	templatedb["html_page"] = html_page
	templatedb["baseTpl"] = baseTpl
}
