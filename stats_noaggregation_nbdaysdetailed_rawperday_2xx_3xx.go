package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func noaggregation_nbdaysdetailed_rawperday_2xx_3xx(args Args) {
	check_if_stats_is_slice := reflect.ValueOf(args).FieldByName("Stats")
	foundcurstat := false
	var mycurstat Statconfig
	if check_if_stats_is_slice.Kind() == reflect.Ptr && check_if_stats_is_slice.IsNil() {
		logger("i wanted to verify if i had to create stats with the daily raw hits, but it seems like all Stats are disabled in the config")

	} else if check_if_stats_is_slice.Kind() == reflect.Slice {
		for _, curstat := range args.Stats {
			if curstat.Statname == "stat_perday_hits_raw_2xx_3xx" {
				foundcurstat = true
				mycurstat = curstat
			}

		}
	}
	if foundcurstat {
		logger("i'm going to generate a table and/or a linechart (depending on the config) with the daily raw hits")
		stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed := myquerydb["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"].stmt
		mintimestamp := int(time.Now().Unix()) - (args.Outputs.Number_of_days_detailed * 86400)
		rows, err := stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed.Query(mintimestamp)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		defer rows.Close()

		MyHeaders := map[string]string{
			"Title_1": "YEAR",
			"Title_2": "MONTH",
			"Title_3": "DAY",
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
		XValues_linegraph_4weekcomp := []string{"today", "day-1", "day-2 ", "day-3", "day-4", "day-5", "day-6"}
		YValues_linegraph_4weekcomp := make(map[string][]int)
		weekcounter := 0
		daycounter := 0
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
			YValues_linegraph["raw hits"] = append(YValues_linegraph["raw hits"], count)
			daycounter++
			if daycounter == 8 {
				weekcounter++
				daycounter = 1
				fmt.Printf("week done... starting with day 0, weeknumber %d\n\n", weekcounter)
			}
			YValues_linegraph_4weekcomp["week -" + strconv.Itoa(weekcounter)] = append(YValues_linegraph_4weekcomp["week -" + strconv.Itoa(weekcounter)], count)
			fmt.Printf("week %d => day %d => value %d\n", weekcounter, daycounter, count)
		}
		if mycurstat.Tableinfo.Table_enabled {
			createtable(args, mycurstat.Tableinfo.Table_filename, mycurstat.Tableinfo.Table_index_name, myTable, mycurstat.Tableinfo.Table_index_group, mycurstat.Tableinfo.Table_index_order)
		}
		if mycurstat.Linegraphinfo.Linegraph_enabled {
			createlinegraph(XValues_linegraph, YValues_linegraph, mycurstat.Linegraphinfo.Linegraph_title, mycurstat.Linegraphinfo.Linegraph_description, args, mycurstat.Linegraphinfo.Linegraph_filename, mycurstat.Linegraphinfo.Linegraph_index_group, mycurstat.Linegraphinfo.Linegraph_index_order)
		}
		if mycurstat.Linegraph4weekinfo.Linegraph_compare4weeks_enabled {
			createlinegraph(XValues_linegraph_4weekcomp, YValues_linegraph_4weekcomp, mycurstat.Linegraph4weekinfo.Linegraph_compare4weeks_title, mycurstat.Linegraph4weekinfo.Linegraph_compare4weeks_description, args, mycurstat.Linegraph4weekinfo.Linegraph_compare4weeks_filename, mycurstat.Linegraph4weekinfo.Linegraph_compare4weeks_index_group, mycurstat.Linegraph4weekinfo.Linegraph_compare4weeks_index_order)
		}
		logger("finished generating a table and/or a linechart with the daily raw hits")
	} else {
		logger("i could not find this stat stat_perday_hits_raw_2xx_3xx in the config. Is it disabled?")
	}
}
