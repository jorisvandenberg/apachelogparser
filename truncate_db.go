package main

import (
    "fmt"
    "github.com/araddon/dateparse"
    "os"
)

func truncate_from(timestring string) {

	 parsedTime, err := dateparse.ParseAny(timestring)
    if err != nil {
        fmt.Printf("Failed to parse date: %v\n", err)
        os.Exit(1)
    }

    unixTimestamp := parsedTime.Unix()
    //fmt.Printf("User input '%s' converted to unix timestamp: %d\n", timestring, unixTimestamp)
   stmt_truncatevisit := myquerydb["stmt_truncatevisit"].stmt
   
		_, err = stmt_truncatevisit.Exec(unixTimestamp)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
}