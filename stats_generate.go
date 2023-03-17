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
		unique_PerDay_hits_table.html
		unique_PerDay_hits_linegraph.html
		unique_PerHour_hits_4WeeksLinegraph.html
	*/
	parameters := []interface{}{mintimestamp}
	tableheaders := map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "NB UNIQUE HITS",
	}
	xaxisfields := []int{0, 1, 2}
	valuefield := 3
	genstats(args, "valid unique hits per day over the last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_PerDay_hits", "stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed", parameters, tableheaders, xaxisfields, valuefield, "raw hits")

	/*
		stat: raw 2xx and 3xx hits over the last 31 (default) days
		expecting 3 htmls:
		raw_PerDay_hits_table.html
		raw_PerDay_hits_linegraph.html
		raw_PerDay_hits_4WeeksLinegraph.html
	*/
	//identical parameters, xaxisfields, valuefield as the stat above... not re-initialising!
	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "NB RAW HITS",
	}
	genstats(args, "valid raw hits per day over the last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_raw_PerDay_hits", "stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed", parameters, tableheaders, xaxisfields, valuefield, "unique hits")

	/*
		stat: unique 2xx and 3xx hits per hour over the last 31 (default) days
		expecting 2 htmls:
		unique_PerHour_hits_table.html
		unique_PerHour_hits_linegraph.html
	*/

	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "HOUR",
		"Title_5": "NB UNIQUE HITS",
	}
	xaxisfields = []int{0, 1, 2, 3}
	valuefield = 4
	genstats(args, "valid unique hits per hour over the last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_PerHour_hits", "stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed", parameters, tableheaders, xaxisfields, valuefield, "unique hits")

	/*
		stat: raw 2xx and 3xx hits per hour over the last 31 (default) days
		expecting 2 htmls:
		raw_PerHour_hits_table.html
		raw_PerHour_hits_linegraph.html
	*/
	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "HOUR",
		"Title_5": "NB RAW HITS",
	}
	genstats(args, "valid raw hits per hour over the last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_raw_PerHour_hits", "stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed", parameters, tableheaders, xaxisfields, valuefield, "raw hits")

	/*
		stat: sum of raw hit per referrer over the last 31 (default) days
		expecting 1 htmls:
		raw_PerHour_ReferringUrls_table.html
	*/
	parameters = []interface{}{mintimestamp, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of raw hits",
	}
	xaxisfields = []int{0}
	valuefield = 1
	genstats(args, "sum of raw hits per referrer over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_raw_PerHour_ReferringUrls", "stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx", parameters, tableheaders, xaxisfields, valuefield, "")

	/*
		stat: sum of unique hits per referrer over the last 31 (default) days
		expecting 1 htmls:
		unique_PerHour_ReferringUrls_table.html
	*/
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of unique hits",
	}
	genstats(args, "sum of unique hits per referrer over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_PerHour_ReferringUrls", "stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx", parameters, tableheaders, xaxisfields, valuefield, "")

	/*
		stat: sum of unique hits per referrer, non empty non self, over the last 31 (default) days
		expecting 1 htmls:
		nunique_PerHour_RefferingUrlsNoEmptyOrSelf_table.html
	*/
	parameters = []interface{}{mintimestamp, args.Generals.Mydomain, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of unique non self, non empty hits",
	}
	genstats(args, "sum of unique hits per referrer, non self non empty, over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_PerHour_RefferingUrlsNoEmptyOrSelf", "stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx", parameters, tableheaders, xaxisfields, valuefield, "")

	/*
		stat: sum of unique hits per referrer, non empty, non self, only tld over the last 31 (default) days
		expecting 1 htmls:
		unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld_table.html
	*/
	parameters = []interface{}{mintimestamp, args.Generals.Mydomain, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of unique non self, non empty hits. TLDs",
	}
	genstats(args, "sum of unique hits per referrer, non self non empty, only tlds, over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld", "stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx", parameters, tableheaders, xaxisfields, valuefield, "")
	logger("finished the function to generate statistics")

	/*
		stat: sum of raw hits per search engine over the last 31 (default) days
		expecting 1 htmls:
		raw_XDaysTotal_HitsFromSearchEngines_table.html
	*/
	parameters = []interface{}{mintimestamp, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "search egine",
		"Title_2": "NB of raw se hits",
	}
	genstats(args, "sum of raw hits per search egine, over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_raw_XDaysTotal_HitsFromSearchEngines", "stmt_count_nbhits_per_searchengine", parameters, tableheaders, xaxisfields, valuefield, "")
	logger("finished the function to generate statistics")

	/*
		stat: sum of unique hits per search engine over the last 31 (default) days
		expecting 1 htmls:
		unique_XDaysTotal_HitsFromSearchEngines_table.html
	*/
	parameters = []interface{}{mintimestamp, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "search egine",
		"Title_2": "NB of unique se hits",
	}
	genstats(args, "sum of unique hits per search egine, over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_XDaysTotal_HitsFromSearchEngines", "stmt_count_unique_nbhits_per_searchengine", parameters, tableheaders, xaxisfields, valuefield, "")
	logger("finished the function to generate statistics")
}
