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
	stmt_truncatevisit_clean_referrer := myquerydb["stmt_truncatevisit_clean_referrer"].stmt

	_, err = stmt_truncatevisit_clean_referrer.Exec()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	stmt_truncatevisit_clean_request := myquerydb["stmt_truncatevisit_clean_request"].stmt

	_, err = stmt_truncatevisit_clean_request.Exec()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	stmt_truncatevisit_clean_user := myquerydb["stmt_truncatevisit_clean_user"].stmt

	_, err = stmt_truncatevisit_clean_user.Exec()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	stmt_truncatevisit_clean_user_ip := myquerydb["stmt_truncatevisit_clean_user_ip"].stmt

	_, err = stmt_truncatevisit_clean_user_ip.Exec()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	stmt_truncatevisit_clean_user_useragent := myquerydb["stmt_truncatevisit_clean_user_useragent"].stmt

	_, err = stmt_truncatevisit_clean_user_useragent.Exec()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	stmt_truncate_alreadyloaded := myquerydb["stmt_truncate_alreadyloaded"].stmt

	_, err = stmt_truncate_alreadyloaded.Exec()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

}
