package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type compare_x_days_weeks_months_config_struct struct {
	What_hours_days_weeks_months                       string
	Number_of_days_weeks_months_compare                int
	Number_of_days_weeks_months_compare_nbitems_inloop int
	Number_of_days_weeks_months_compare_legenda        string
	Linegraph_compare_x_days_weeks_months_title        string
	Linegraph_compare_x_days_weeks_months_description  string
	Linegraph_compare_x_days_weeks_months_filename     string
	Linegraph_compare_x_days_weeks_months_index_group  string
	Linegraph_compare_x_days_weeks_months_index_order  int
}

func genstats(args Args, string_for_log string, statname_from_conf string, querydb_key string, parameters []interface{}, tableheaders map[string]string, xaxisfields []int, valuefield int, legende string) bool {
	//what_hours_days_weeks_months : usually hour or day
	//number_of_days_weeks_months_compare : usually 4
	//number_of_days_weeks_months_compare_nbitems_inloop : 8 when your query is grouping by day and you want to compare weeks, 32 if you want to compare months, 25 if you group by hour and want to compare days,...
	//number_of_days_weeks_months_compare_legenda: usually week if you group by day and have 8 items in loop
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

	var compare_x_days_weeks_months_parameters_parts []string

	//als mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_enabled == False what_hours_days_weeks_months,number_of_days_weeks_months_compare,number_of_days_weeks_months_compare_nbitems_inloop,number_of_days_weeks_months_compare_legenda vullen met standaard value, zoniet vullen met echte values en er dan door loopen
	compare_x_days_weeks_months_config_struct_slice := []compare_x_days_weeks_months_config_struct{}
	if !mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_enabled {
		Linegraph_compare_x_days_weeks_months_index_order_int, _ := strconv.Atoi(mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_index_order)
		compare_x_days_weeks_months_config_struct_slice = append(compare_x_days_weeks_months_config_struct_slice, compare_x_days_weeks_months_config_struct{"day", 4, 8, "week", mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_title, mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_description, mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_filename, mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_index_group, Linegraph_compare_x_days_weeks_months_index_order_int})
	}
	if mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_enabled {

		linegraph_compare_x_days_weeks_months_parameters_slice := strings.Split(mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_parameters, "&")
		linegraph_compare_x_days_weeks_months_title_slice := strings.Split(mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_title, "&")
		Linegraph_compare_x_days_weeks_months_description_slice := strings.Split(mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_description, "&")
		Linegraph_compare_x_days_weeks_months_filename_slice := strings.Split(mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_filename, "&")
		Linegraph_compare_x_days_weeks_months_index_group_slice := strings.Split(mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_index_group, "&")
		Linegraph_compare_x_days_weeks_months_index_order_slice := strings.Split(mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_index_order, "&")
		indexnumber := 0
		for _, linegraph_compare_x_days_weeks_months_parameters_slice_current := range linegraph_compare_x_days_weeks_months_parameters_slice {

			compare_x_days_weeks_months_parameters_parts = strings.Split(linegraph_compare_x_days_weeks_months_parameters_slice_current, ",")
			//fmt.Printf("parameters: %+v\n", compare_x_days_weeks_months_parameters_parts)

			compare_x_days_weeks_months_parameters_parts_1, _ := strconv.Atoi(compare_x_days_weeks_months_parameters_parts[1])
			compare_x_days_weeks_months_parameters_parts_2, _ := strconv.Atoi(compare_x_days_weeks_months_parameters_parts[2])
			Linegraph_compare_x_days_weeks_months_index_order_slice_string, _ := strconv.Atoi(Linegraph_compare_x_days_weeks_months_index_order_slice[indexnumber])

			compare_x_days_weeks_months_config_struct_slice = append(compare_x_days_weeks_months_config_struct_slice, compare_x_days_weeks_months_config_struct{compare_x_days_weeks_months_parameters_parts[0], compare_x_days_weeks_months_parameters_parts_1, compare_x_days_weeks_months_parameters_parts_2, compare_x_days_weeks_months_parameters_parts[3], linegraph_compare_x_days_weeks_months_title_slice[indexnumber], Linegraph_compare_x_days_weeks_months_description_slice[indexnumber], Linegraph_compare_x_days_weeks_months_filename_slice[indexnumber], Linegraph_compare_x_days_weeks_months_index_group_slice[indexnumber], Linegraph_compare_x_days_weeks_months_index_order_slice_string})
			indexnumber += 1
		}
	}

	for _, current_x_days_weeks_months := range compare_x_days_weeks_months_config_struct_slice {
		logger("start " + string_for_log)
		myQuery := myquerydb[querydb_key].stmt
		rows, err := myQuery.Query(parameters...)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			fmt.Printf("%s\n", querydb_key)
			fmt.Printf("%+v\n", myQuery)
			os.Exit(1)
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
		XValues_linegraph_4weekcomp := []string{"now", current_x_days_weeks_months.What_hours_days_weeks_months + "-1", current_x_days_weeks_months.What_hours_days_weeks_months + "-2 ", current_x_days_weeks_months.What_hours_days_weeks_months + "-3", current_x_days_weeks_months.What_hours_days_weeks_months + "-4", current_x_days_weeks_months.What_hours_days_weeks_months + "-5", current_x_days_weeks_months.What_hours_days_weeks_months + "-6"}
		YValues_linegraph_4weekcomp := make(map[string][]int)
		weekcounter := 0
		daycounter := 0
		columns, err := rows.Columns()
		if err != nil {
			return false
		}

		var result []interface{}
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		for rows.Next() {

			err := rows.Scan(valuePtrs...)
			if err != nil {
				return false
			}
			result = append(result, values)

			MyData := make(map[string]string)

			counter := 0
			for _, value := range values {

				switch v := value.(type) {
				case string:
					// value is a string, so we can add it directly to MyData
					MyData["Value_"+strconv.Itoa(counter)] = v
				case int64:
					// value is an int64, so we need to convert it to a string first
					MyData["Value_"+strconv.Itoa(counter)] = strconv.FormatInt(v, 10)
				default:
					// value is neither a string nor an int64, so handle the error case
					panic("unsupported type")
				}
				counter++

			}

			titel := ""

			for _, xaxisfield := range xaxisfields {
				titel += values[xaxisfield].(string)
			}
			myTable.Data = append(myTable.Data, MyData)

			XValues_linegraph = append(XValues_linegraph, titel)
			YValues_linegraph[legende] = append(YValues_linegraph[legende], int(values[valuefield].(int64)))
			daycounter++
			if daycounter == current_x_days_weeks_months.Number_of_days_weeks_months_compare_nbitems_inloop {
				weekcounter++
				daycounter = 1
			}
			if weekcounter < current_x_days_weeks_months.Number_of_days_weeks_months_compare {
				YValues_linegraph_4weekcomp[current_x_days_weeks_months.Number_of_days_weeks_months_compare_legenda+" -"+strconv.Itoa(weekcounter)] = append(YValues_linegraph_4weekcomp[current_x_days_weeks_months.Number_of_days_weeks_months_compare_legenda+" -"+strconv.Itoa(weekcounter)], int(values[valuefield].(int64)))
			}

		}
		err = rows.Err()
		if err != nil {
			return false
		}

		if mycurstat.Tableinfo.Table_enabled {
			createtable(args, mycurstat.Tableinfo.Table_filename, mycurstat.Tableinfo.Table_index_name, myTable, mycurstat.Tableinfo.Table_index_group, mycurstat.Tableinfo.Table_index_order)
		}

		if mycurstat.Linegraphinfo.Linegraph_enabled {
			createlinegraph(XValues_linegraph, YValues_linegraph, mycurstat.Linegraphinfo.Linegraph_title, mycurstat.Linegraphinfo.Linegraph_description, args, mycurstat.Linegraphinfo.Linegraph_filename, mycurstat.Linegraphinfo.Linegraph_index_group, mycurstat.Linegraphinfo.Linegraph_index_order)
		}

		if mycurstat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_enabled {
			createlinegraph(XValues_linegraph_4weekcomp, YValues_linegraph_4weekcomp, current_x_days_weeks_months.Linegraph_compare_x_days_weeks_months_title, current_x_days_weeks_months.Linegraph_compare_x_days_weeks_months_description, args, current_x_days_weeks_months.Linegraph_compare_x_days_weeks_months_filename, current_x_days_weeks_months.Linegraph_compare_x_days_weeks_months_index_group, current_x_days_weeks_months.Linegraph_compare_x_days_weeks_months_index_order)
		}

		logger("stopped " + string_for_log)
	}
	return true
}
