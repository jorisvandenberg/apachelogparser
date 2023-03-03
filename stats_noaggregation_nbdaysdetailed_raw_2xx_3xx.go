package main

import (
	"time"
	"fmt"
)
func noaggregation_nbdaysdetailed_raw_2xx_3xx(args args) {
	stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed := myquerydb["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"].stmt
	mintimestamp := int(time.Now().Unix()) - (args.number_of_days_detailed * 86400)
	rows, err := stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed.Query(mintimestamp)
	if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	defer rows.Close()
	for rows.Next() {
		var year, month, day, hour, count int
		if err := rows.Scan(&year, &month, &day, &hour, &count); err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		fmt.Printf("%d/%d/%d => %d", year, month, day, count)
	}
}