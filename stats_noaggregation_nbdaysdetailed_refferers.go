package main

import (
	"fmt"
	"strconv"
	"time"
	"reflect"
)

func stats_noaggregation_nbdaysdetailed_refferers(args Args) {
	check_if_stats_is_slice := reflect.ValueOf(args).FieldByName("Stats")
	foundcurstat := false
	var mycurstat Statconfig
	if check_if_stats_is_slice.Kind() == reflect.Ptr && check_if_stats_is_slice.IsNil() {
		logger("i wanted to verify if i had to create stats with the raw referrers, but it seems like all Stats are disabled in the config")

	} else if check_if_stats_is_slice.Kind() == reflect.Slice {
		for _, curstat := range args.Stats {
			if curstat.Statname == "stat_perhour_referrers_raw_2xx_3xx" {
				foundcurstat = true
				mycurstat = curstat
			}

		}
	}
	if foundcurstat {
	logger("i'm goig to generate a table with a sum of raw hits per referrer for the last " + strconv.Itoa(args.Outputs.Number_of_days_detailed) + " days")
	stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx := myquerydb["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"].stmt
	mintimestamp := int(time.Now().Unix()) - (args.Outputs.Number_of_days_detailed * 86400)
	rows, err := stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx.Query(mintimestamp, args.Outputs.Numberofreferrers)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer rows.Close()
	MyHeaders := map[string]string{
		"Title_1": "NB of raw hits",
		"Title_2": "REFERRER",
	}
	myTable := Table{
		Pagetitle:       mycurstat.Tableinfo.Table_title,
		Pagedescription: mycurstat.Tableinfo.Table_description,
		Pagecontent:     mycurstat.Tableinfo.Table_pagecontent,
		Pagefooter:     mycurstat.Tableinfo.Table_pagefooter,
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
	createtable(args,mycurstat.Tableinfo.Table_filename, mycurstat.Tableinfo.Table_index_name, myTable, mycurstat.Tableinfo.Table_index_group, mycurstat.Tableinfo.Table_index_order)
	logger("finished generating a table with a sum of raw hits per referrer for the last " + strconv.Itoa(args.Outputs.Number_of_days_detailed) + " days")
	} else {
		logger("i could not find this stat stat_perhour_referrers_raw_2xx_3xx in the config. Is it disabled?")
	}
}
