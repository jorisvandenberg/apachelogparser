package main

import (
	"strconv"
	"time"
	)
	
func generatestats(args Args) {
	logger("started the function to generate statistics")
	
	mintimestamp := int(time.Now().Unix()) - (args.Outputs.Number_of_days_detailed * 86400)
	
	/*
	stat: unique 2xx and 3xx hits over the last 31 (default) days
	expecting 3 htmls: 
	noaggregation_nbdaysdetailed_uniqueperday_2xx_3xx_table.html
	noaggregation_nbdaysdetailed_uniqueperday_2xx_3xx_linegraph.html
	noaggregation_nbdaysdetailed_unique_2xx_3xx_linegraph_4weekcompare.html
	*/
	parameters := []interface{}{mintimestamp}
	tableheaders := map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "NB UNIQUE HITS",
	}
	xaxisfields := []int{0,1,2}
	valuefield := 3
	genstats(args, "valid unique hits per day over the last " + strconv.Itoa(args.Outputs.Number_of_days_detailed) + " days" , "stat_perday_hits_unique_2xx_3xx", "stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed", parameters, tableheaders,xaxisfields, valuefield, "raw hits" )
	
	/*
	stat: raw 2xx and 3xx hits over the last 31 (default) days
	expecting 3 htmls: 
	noaggregation_nbdaysdetailed_rawperday_2xx_3xx_table.html
	noaggregation_nbdaysdetailed_rawperday_2xx_3xx_linegraph.html
	noaggregation_nbdaysdetailed_raw_2xx_3xx_linegraph_4weekcompare.html
	*/
	//identical parameters, xaxisfields, valuefield as the stat above... not re-initialising!
	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "NB RAW HITS",
	}
	genstats(args, "valid raw hits per day over the last " + strconv.Itoa(args.Outputs.Number_of_days_detailed) + " days" , "stat_perday_hits_raw_2xx_3xx", "stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed", parameters, tableheaders,xaxisfields, valuefield, "unique hits" )
	
	
	/*
	stat: unique 2xx and 3xx hits per hour over the last 31 (default) days
	expecting 2 htmls: 
	noaggregation_nbdaysdetailed_unique_2xx_3xx_table.html
	noaggregation_nbdaysdetailed_unique_2xx_3xx_linegraph.html
	*/

	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "HOUR",
		"Title_5": "NB UNIQUE HITS",
	}
	xaxisfields = []int{0,1,2,3}
	valuefield = 4
	genstats(args, "valid unique hits per hour over the last " + strconv.Itoa(args.Outputs.Number_of_days_detailed) + " days" , "stat_perhour_hits_unique_2xx_3xx", "stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed", parameters, tableheaders,xaxisfields, valuefield, "unique hits" )
	
	/*
	stat: raw 2xx and 3xx hits per hour over the last 31 (default) days
	expecting 2 htmls: 
	noaggregation_nbdaysdetailed_raw_2xx_3xx_table.html
	noaggregation_nbdaysdetailed_raw_2xx_3xx_linegraph.html
	*/
	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "HOUR",
		"Title_5": "NB RAW HITS",
	}
	genstats(args, "valid raw hits per hour over the last " + strconv.Itoa(args.Outputs.Number_of_days_detailed) + " days" , "stat_perhour_hits_raw_2xx_3xx", "stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed", parameters, tableheaders,xaxisfields, valuefield, "raw hits" )
	
	
	/*
	stat: sum of raw hit per referrer over the last 31 (default) days
	expecting 1 htmls: 
	noaggregation_nbdaysdetailed_rawreferrers_2xx_3xx_table.html
	*/
	parameters = []interface{}{mintimestamp, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of raw hits",
			
	}
	xaxisfields = []int{0}
	valuefield = 1
	genstats(args, "sum of raw hits per referrer over last " + strconv.Itoa(args.Outputs.Number_of_days_detailed) + " days" , "stat_perhour_referrers_raw_2xx_3xx", "stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx", parameters, tableheaders,xaxisfields, valuefield, "" )
	
	
	/*
	stat: sum of unique hits per referrer over the last 31 (default) days
	expecting 1 htmls: 
	noaggregation_nbdaysdetailed_uniquereferrers_2xx_3xx_table.html
	*/
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of unique hits",
			
	}
	genstats(args, "sum of unique hits per referrer over last " + strconv.Itoa(args.Outputs.Number_of_days_detailed) + " days" , "stat_perhour_referrers_unique_2xx_3xx", "stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx", parameters, tableheaders,xaxisfields, valuefield, "" )
	
	//noaggregation_nbdaysdetailed_raw_2xx_3xx(args)
	//noaggregation_nbdaysdetailed_unique_2xx_3xx(args)
	//noaggregation_nbdaysdetailed_rawperday_2xx_3xx(args)
	//noaggregation_nbdaysdetailed_uniqueperday_2xx_3xx(args)
	//stats_noaggregation_nbdaysdetailed_refferers(args)
	//stats_noaggregation_nbdaysdetailed_unique_refferers(args)
	stats_noaggregation_nbdaysdetailed_unique_refferers_noemptyorself(args)
	stats_noaggregation_nbdaysdetailed_unique_refferers_noemptyorself_onlytld(args)
	logger("finished the function to generate statistics")
}
