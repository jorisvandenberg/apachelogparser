package main

import (
	"fmt"
	"strconv"
	"time"
)

func stats_noaggregation_nbdaysdetailed_refferers(args args) {
	logger("i'm goig to generate a table with a sum of raw hits per referrer for the last " + strconv.Itoa(args.outputs.Number_of_days_detailed) + " days")
	stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx := myquerydb["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"].stmt
	mintimestamp := int(time.Now().Unix()) - (args.outputs.Number_of_days_detailed * 86400)
	rows, err := stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx.Query(mintimestamp, args.outputs.Numberofreferrers)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer rows.Close()
	MyHeaders := map[string]string{
		"Title_1": "NB of raw hits",
		"Title_2": "REFERRER",
	}
	myTable := Table{
		Pagetitle:       "Number of raw 2xx and 3xx hits coming from a referrer over the last " + strconv.Itoa(args.outputs.Number_of_days_detailed) + " days",
		Pagedescription: "Count of all raw succesfull hits (filtering out all 4xx and 5xx return codes).",
		Pagecontent:     []string{"We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you'll get stats untill 15:34 x days ago)."},
		Pagefooter:      "only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they'll never be shown!",
		Headers:         MyHeaders,
		Data:            []map[string]string{},
	}
	for rows.Next() {
		var count int
		var referrer string
		if err := rows.Scan(&referrer, &count); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		MyData := map[string]string{
			"Value_1": referrer,
			"Value_2": strconv.Itoa(count),
		}
		myTable.Data = append(myTable.Data, MyData)
	}
	createtable(args, "stats_noaggregation_nbdaysdetailed_refferers_table.html", "table of the raw 2xx and 3xx per referrer over the last "+strconv.Itoa(args.outputs.Number_of_days_detailed)+" days", myTable, "referrer", 3)
	logger("finished generating a table with a sum of raw hits per referrer for the last " + strconv.Itoa(args.outputs.Number_of_days_detailed) + " days")
}
