package main

import (
	"fmt"
	"strconv"
	"time"
)

func noaggregation_nbdaysdetailed_raw_2xx_3xx(args args) {
	stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed := myquerydb["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"].stmt
	mintimestamp := int(time.Now().Unix()) - (args.number_of_days_detailed * 86400)
	rows, err := stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed.Query(mintimestamp)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer rows.Close()

	MyHeaders := map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "HOUR",
		"Title_5": "NB RAW HITS",
	}
	myTable := Table{
		Pagetitle:       "Number of raw 2xx and 3xx hits per hour over th last " + strconv.Itoa(args.number_of_days_detailed) + " days",
		Pagedescription: "Count of all raw succesfull hits (filtering out all 4xx and 5xx return codes).",
		Pagecontent:     []string{"We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you'll get stats untill 15:34 x days ago)."},
		Pagefooter:      "only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they'll never be shown!",
		Headers:         MyHeaders,
		Data:            []map[string]string{},
	}

	var XValues_linegraph []string
	YValues_linegraph := make(map[string][]int)
	for rows.Next() {
		var year, month, day, hour, count int
		if err := rows.Scan(&year, &month, &day, &hour, &count); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		MyData := map[string]string{
			"Value_1": strconv.Itoa(year),
			"Value_2": strconv.Itoa(month),
			"Value_3": strconv.Itoa(day),
			"Value_4": strconv.Itoa(hour),
			"Value_5": strconv.Itoa(count),
		}
		myTable.Data = append(myTable.Data, MyData)
		XValues_linegraph = append(XValues_linegraph, strconv.Itoa(year)+"-"+strconv.Itoa(month)+"-"+strconv.Itoa(day)+":"+strconv.Itoa(hour))
		YValues_linegraph["raw hits"] = append(YValues_linegraph["raw hits"], count)
	}

	createtable(args, "noaggregation_nbdaysdetailed_raw_2xx_3xx_table.html", "table of the raw 2xx and 3xx per hour over the last "+strconv.Itoa(args.number_of_days_detailed)+" days", myTable)
	PreChartText = ""
	PostChartText = ""
	createlinegraph(XValues_linegraph, YValues_linegraph, "line graph of the raw hits with status 2xx and 3xx", "Count of all raw succesfull hits (filtering out all 4xx and 5xx return codes).", args, "noaggregation_nbdaysdetailed_raw_2xx_3xx_linegraph.html", "hits", 2)

}
