package main

import (
	"fmt"
	"strconv"
	"time"
)

func noaggregation_nbdaysdetailed_uniqueperday_2xx_3xx(args args) {
	stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed := myquerydb["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"].stmt
	mintimestamp := int(time.Now().Unix()) - (args.number_of_days_detailed * 86400)
	rows, err := stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed.Query(mintimestamp)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer rows.Close()

	MyHeaders := map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_5": "NB unique HITS",
	}
	myTable := Table{
		Pagetitle:       "Number of unique 2xx and 3xx hits per day over th last " + strconv.Itoa(args.number_of_days_detailed) + " days",
		Pagedescription: "Count of all unique succesfull hits (filtering out all 4xx and 5xx return codes).",
		Pagecontent:     []string{"Unique means: only counting the first hit of a certain user for that day. If a visitor generates a hit at 11:59 and a second hit at 12:01 he counts as a unique at 11:00 and 12:00. If he then creates a thirth hit at 12:05, said hit will not be counted", "We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you'll get stats untill 15:34 x days ago)."},
		Pagefooter:      "only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they'll never be shown!",
		Headers:         MyHeaders,
		Data:            []map[string]string{},
	}

	var XValues_linegraph []string
	YValues_linegraph := make(map[string][]int)
	for rows.Next() {
		var year, month, day, count int
		if err := rows.Scan(&year, &month, &day, &count); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		MyData := map[string]string{
			"Value_1": strconv.Itoa(year),
			"Value_2": strconv.Itoa(month),
			"Value_3": strconv.Itoa(day),
			"Value_5": strconv.Itoa(count),
		}
		myTable.Data = append(myTable.Data, MyData)
		XValues_linegraph = append(XValues_linegraph, strconv.Itoa(year)+"-"+strconv.Itoa(month)+"-"+strconv.Itoa(day))
		YValues_linegraph["unique hits"] = append(YValues_linegraph["unique hits"], count)
	}

	createtable(args, "noaggregation_nbdaysdetailed_uniqueperday_2xx_3xx_table.html", "table of the unique 2xx and 3xx per day over the last "+strconv.Itoa(args.number_of_days_detailed)+" days", myTable)
	PreChartText = ""
	PostChartText = ""
	createlinegraph(XValues_linegraph, YValues_linegraph, "line graph of the unique hits with status 2xx and 3xx per day", "Count of all unique succesfull hits (filtering out all 4xx and 5xx return codes).", args, "noaggregation_nbdaysdetailed_uniqueperday_2xx_3xx_linegraph.html", "hits", 1)

}
