package main

func stats_noaggregation_nbdaysdetailed_refferers (args args) {
	
	stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx := myquerydb["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"].stmt
	mintimestamp := int(time.Now().Unix()) - (args.number_of_days_detailed * 86400)
	rows, err := stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed.Query(mintimestamp, args.numberofreferrers)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer rows.Close()
	MyHeaders := map[string]string{
		"Title_1": "NB of raw hits",
		"Title_2": "REFERRER",
	}
	myTable := Table{
		Pagetitle:       "Number of raw 2xx and 3xx hits coming from a referrer over the last " + strconv.Itoa(args.number_of_days_detailed) + " days",
		Pagedescription: "Count of all raw succesfull hits (filtering out all 4xx and 5xx return codes).",
		Pagecontent:     []string{"We limit the output to the number of days that were defined in your config.ini file with a sliding window (so if you run this tool at 15:34 you'll get stats untill 15:34 x days ago)."},
		Pagefooter:      "only hits that were actually loaded are shown, so if you filtered out certain lines in your config.ini they'll never be shown!",
		Headers:         MyHeaders,
		Data:            []map[string]string{},
	}
	var XValues_linegraph []string
	YValues_linegraph := make(map[string][]int)
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
		XValues_linegraph = append(XValues_linegraph, strconv.Itoa(year)+"-"+strconv.Itoa(month)+"-"+strconv.Itoa(day))
		YValues_linegraph["raw hits from referrer"] = append(YValues_linegraph["raw hits from referrer"], count)
		createtable(args, "stats_noaggregation_nbdaysdetailed_refferers_table.html", "table of the raw 2xx and 3xx per referrer over the last "+strconv.Itoa(args.number_of_days_detailed)+" days", myTable, "referrer", 3)
	PreChartText = ""
	PostChartText = ""
	createlinegraph(XValues_linegraph, YValues_linegraph, "line graph of the raw hits with status 2xx and 3xx per referrer", "Count of all raw succesfull hits (filtering out all 4xx and 5xx return codes).", args, "stats_noaggregation_nbdaysdetailed_refferers_linegraph.html", "referrer", 3)
	}
}