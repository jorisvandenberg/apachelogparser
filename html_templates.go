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
		<h2>{{.Pagedescription}}</h2>
		{{range .Pagecontent}}
		<p>{{.}}</p>
		{{end}}
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
		<p>{{.Pagefooter}}</p>
	</body>
</html>`
	/*
	   	html_index := `
	   <html>
	   <head>
	   </head>
	   <body>
	   <h1>index page</h1>
	   {{range $key, $value := .}}
	   		<h2>{{ $key }}</h2>
	   		{{range $value}}
	   			<p>{{.Textpre}}<a href="{{.Url}}" target="_blank">{{.Title}}</a>{{.Textpost}}</p>
	   		{{end}}
	   	{{end}}
	   	</body>
	   	</html>`
	*/

	html_index := `
<!DOCTYPE html>
<html>
<head>
	<title>My Statistics</title>
	<style>
		body {
			margin: 0;
			padding: 0;
			display: flex;
			flex-direction: row;
			height: 100vh;
			font-family: Arial, sans-serif;
		}

		nav {
			display: flex;
			flex-direction: column;
			background-color: #f5f5f5;
			padding: 10px;
			width: 250px;
		}

		nav a {
			text-decoration: none;
			color: #333;
			font-size: 16px;
			margin-bottom: 5px;
			padding-left: 0px;
		}

		nav a:hover {
			color: #fff;
			background-color: #333;
		}

		nav h2 {
			margin-top: 0;
			margin-bottom: 0px;
			padding-left: 10px;
			font-size: 18px;
			font-weight: bold;
		}

		iframe {
			flex: 1;
			height: 100%;
			border: none;
		}
	</style>
</head>
<body>
	<nav>
	
	
	{{range $key, $value := .}}
		<h2>{{ $key }}</h2>
		<ul>
		{{range $value}}
			<li>{{.Textpre}}<a href="{{.Url}}" target="statzframe">{{.Title}}</a>{{.Textpost}}</li>
		{{end}}
		</ul>
	{{end}}
	</nav>
	<iframe src="https://jorisvandenberg.github.io/" name="statzframe"></iframe>
	<script>
		// When a link is clicked, scroll to the top of the iframe
		var iframe = document.querySelector('iframe');
		iframe.onload = function() {
			iframe.contentWindow.scrollTo(0, 0);
		}
	</script>
</body>
</html>


`
	templatedb["table_tmpl"] = table_tmpl
	templatedb["html_index"] = html_index
	templatedb["html_page"] = html_page
	templatedb["baseTpl"] = baseTpl
}
