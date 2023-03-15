package main

import (
	"fmt"
	"reflect"
	"strconv"
	"os"
)

func genstats(args Args, string_for_log string, statname_from_conf string, querydb_key string, parameters []interface{}, tableheaders map[string]string, xaxisfields []int, valuefield int, legende string ) bool {
	/*
			tableheaders := map[string]string{
					"Title_1": "YEAR",
					"Title_2": "MONTH",
					"Title_3": "DAY",
					"Title_4": "HOUR",
					"Title_5": "NB RAW HITS",
				}
				parameters := []interface{}{"value1", "value2", "value3"}
				tableheaders := map[string]string{
					"Title_1": "YEAR",
					"Title_2": "MONTH",
					"Title_3": "DAY",
					"Title_4": "HOUR",
					"Title_5": "NB RAW HITS",
				}
				sqlreturnvalues := []interface{}{
		        []interface{}{"year", "int"},
		        []interface{}{"month", "int"},
		        []interface{}{"day", "int"},
		        []interface{}{"hour", "int"},
		        []interface{}{"nm raw hits", "int"},
		    }
			year := sqlreturnvalues[0].([]interface{})[1].(int)
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
			//arr := value.(string)
			//MyData["Value_" + strconv.Itoa(counter)] = strconv.FormatInt(value.(int64), 10)
			
			switch v := value.(type) {
			case string:
				// value is a string, so we can add it directly to MyData
				MyData["Value_" + strconv.Itoa(counter)]  = v
			case int64:
				// value is an int64, so we need to convert it to a string first
				MyData["Value_" + strconv.Itoa(counter)] = strconv.FormatInt(v, 10)
			default:
				// value is neither a string nor an int64, so handle the error case
				panic("unsupported type")
			}
			
			
			
			//fmt.Printf("%d => %T => %s\n", counter, value, value)
			counter++
			
		}
		
		titel := ""
		
		for _, xaxisfield := range xaxisfields {
			titel += values[xaxisfield].(string)
		}
		fmt.Printf("%s\n", values[0])
			myTable.Data = append(myTable.Data, MyData)
		
			
			XValues_linegraph = append(XValues_linegraph, titel)
			YValues_linegraph[legende] = append(YValues_linegraph[legende], int(values[valuefield].(int64)))
		
	}
	err = rows.Err()
	if err != nil {
		return false
	}
	fmt.Printf("%+v\n", columns)
	fmt.Printf("%+v\n", result)
	fmt.Printf("%+v\n", myTable)
	
		if mycurstat.Tableinfo.Table_enabled {
			createtable(args, mycurstat.Tableinfo.Table_filename, mycurstat.Tableinfo.Table_index_name, myTable, mycurstat.Tableinfo.Table_index_group, mycurstat.Tableinfo.Table_index_order)
		}
		
		if mycurstat.Linegraphinfo.Linegraph_enabled {
			createlinegraph(XValues_linegraph, YValues_linegraph, mycurstat.Linegraphinfo.Linegraph_title, mycurstat.Linegraphinfo.Linegraph_description, args, mycurstat.Linegraphinfo.Linegraph_filename, mycurstat.Linegraphinfo.Linegraph_index_group, mycurstat.Linegraphinfo.Linegraph_index_order)
		}
	
	logger("stopped " + string_for_log)
	return true
}
