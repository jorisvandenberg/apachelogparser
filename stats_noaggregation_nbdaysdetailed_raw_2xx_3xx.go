package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func noaggregation_nbdaysdetailed_raw_2xx_3xx(args Args) {
	check_if_stats_is_slice := reflect.ValueOf(args).FieldByName("Stats")
	foundcurstat := false
	var mycurstat Statconfig
	if check_if_stats_is_slice.Kind() == reflect.Ptr && check_if_stats_is_slice.IsNil() {
		logger("i wanted to verify if i had to create stats with the hourly raw hits, but it seems like all Stats are disabled in the config")

	} else if check_if_stats_is_slice.Kind() == reflect.Slice {
		for _, curstat := range args.Stats {
			if (curstat.Statname == "stat_perhour_hits_raw_2xx_3xx") {
			foundcurstat = true
			mycurstat = curstat
			}
			
		}
	}
	//als foundcurstat is true mag ik statistiek maken, de config zit in mycurstat
	if foundcurstat {
		//fmt.Printf("%+v", mycurstat)

		logger("i'm going to generate a table and/or a linechart (depending on the config) with the hourly raw hits")
		stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed := myquerydb["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"].stmt
		mintimestamp := int(time.Now().Unix()) - (args.Outputs.Number_of_days_detailed * 86400)
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
			Pagetitle:       mycurstat.Tableinfo.Table_title,
			Pagedescription: mycurstat.Tableinfo.Table_description,
			Pagecontent:     mycurstat.Tableinfo.Table_pagecontent,
			Pagefooter:      mycurstat.Tableinfo.Table_pagefooter,
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
		if mycurstat.Tableinfo.Table_enabled {
			createtable(args, mycurstat.Tableinfo.Table_filename, mycurstat.Tableinfo.Table_index_name, myTable, mycurstat.Tableinfo.Table_index_group, mycurstat.Tableinfo.Table_index_order)
		}
		if mycurstat.Linegraphinfo.Linegraph_enabled {
			createlinegraph(XValues_linegraph, YValues_linegraph, mycurstat.Linegraphinfo.Linegraph_title, mycurstat.Linegraphinfo.Linegraph_description, args, mycurstat.Linegraphinfo.Linegraph_filename, mycurstat.Linegraphinfo.Linegraph_index_group, mycurstat.Linegraphinfo.Linegraph_index_order)
		}
		logger("finished generating a table and/or a linechart with the hourly raw hits")
	} else {
		logger("i could not find this stat stat_perhour_hits_raw_2xx_3xx in the config. Is it disabled?")
	}
}
