package main

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func genstats(args Args, string_for_log string, statname_from_conf string, querydb_key string, parameters []string, tableheaders map[string]string ) bool {
	/*
	tableheaders := map[string]string{
			"Title_1": "YEAR",
			"Title_2": "MONTH",
			"Title_3": "DAY",
			"Title_4": "HOUR",
			"Title_5": "NB RAW HITS",
		}
		parameters := []interface{}{"value1", "value2", "value3"}
	*/
	check_if_stats_is_slice := reflect.ValueOf(args).FieldByName("Stats")
	foundcurstat := false
	var mycurstat Statconfig
	
	if check_if_stats_is_slice.Kind() == reflect.Ptr && check_if_stats_is_slice.IsNil() {
		logger("i wanted to run: " + statname_from_conf + ", but my argumentparser did not find any Stats defined in the config!!!")
		return false
	}
	if check_if_stats_is_slice.Kind() == reflect.Slice {
		for _, curstat := range args.Stats {
			if curstat.Statname == statname_from_conf {
				foundcurstat = true
				mycurstat = curstat
			}
		}

	} else {
		logger("i wanted to run: " + statname_from_conf + ", but my argumentparser did not find any Stats defined in the config!!!")
		return false
	}

	if !foundcurstat {
		logger("i wanted to run: " + statname_from_conf + ", but i could not find said section in the config.ini!!!")
		return false
	}

	logger("start " + string_for_log)
	myQuery := myquerydb[querydb_key].stmt
		mintimestamp := int(time.Now().Unix()) - (args.Outputs.Number_of_days_detailed * 86400)
		rows, err := myQuery.Query(parameters...)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		defer rows.Close()
		myTable := Table{
			Pagetitle:       mycurstat.Tableinfo.Table_title,
			Pagedescription: mycurstat.Tableinfo.Table_description,
			Pagecontent:     mycurstat.Tableinfo.Table_pagecontent,
			Pagefooter:      mycurstat.Tableinfo.Table_pagefooter,
			Headers:         tableheaders,
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
}