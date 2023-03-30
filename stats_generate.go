package main

import (
	"strconv"
	"time"
	"strings"
)

func getmaxdaysfromargs(args Args, configname string) int {
	max_days_from_args := args.Outputs.Number_of_days_detailed
	if args.Outputs.Max_number_of_days > max_days_from_args {
		max_days_from_args = args.Outputs.Max_number_of_days
	}

	for _, curStat := range args.Stats {
		if curStat.Statname == configname {
			if curStat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_enabled {
				linegraph_compare_x_days_weeks_months_parameters_slice := strings.Split(curStat.Linegraph4weekinfo.Linegraph_compare_x_days_weeks_months_parameters, "&")
				for _, linegraph_compare_x_days_weeks_months_parameters_slice_current := range linegraph_compare_x_days_weeks_months_parameters_slice {
					compare_x_days_weeks_months_parameters_parts := strings.Split(linegraph_compare_x_days_weeks_months_parameters_slice_current, ",")
					compare_x_days_weeks_months_parameters_parts_1, _ := strconv.Atoi(compare_x_days_weeks_months_parameters_parts[1])
					compare_x_days_weeks_months_parameters_parts_2, _ := strconv.Atoi(compare_x_days_weeks_months_parameters_parts[2])
					product_groupsize_nbgroups := compare_x_days_weeks_months_parameters_parts_1 * compare_x_days_weeks_months_parameters_parts_2
					if product_groupsize_nbgroups > max_days_from_args {
						max_days_from_args = product_groupsize_nbgroups
					}

				}
			}
			
		}
		
	}
	
	return max_days_from_args
}

func generatestats(args Args) {
	logger("started the function to generate statistics")

	//mintimestamp_Max_number_of_days := int(time.Now().Unix()) - (args.Outputs.Max_number_of_days * 86400)
	mintimestamp_Number_of_days_detailed := int(time.Now().Unix()) - (args.Outputs.Number_of_days_detailed * 86400)

	/*
		stat: unique 2xx and 3xx hits over the last 31 (default) days
		expecting 3 htmls:
		unique_PerDay_hits_table.html
		unique_PerDay_hits_linegraph.html
		unique_PerHour_hits_4WeeksLinegraph.html
	*/
		//parameters := []interface{}{mintimestamp_Max_number_of_days}
	parameters := []interface{}{getmaxdaysfromargs(args, "conf_stat_unique_PerDay_hits")}
	tableheaders := map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "NB UNIQUE HITS",
	}
	xaxisfields := []int{0, 1, 2}
	valuefield := 3
	genstats(args, "valid unique hits per day over the last "+strconv.Itoa(args.Outputs.Max_number_of_days)+" days", "conf_stat_unique_PerDay_hits", "stmt_unique_PerDay_hits", parameters, tableheaders, xaxisfields, valuefield, "raw hits", true)

	/*
		stat: raw 2xx and 3xx hits over the last 31 (default) days
		expecting 3 htmls:
		raw_PerDay_hits_table.html
		raw_PerDay_hits_linegraph.html
		raw_PerDay_hits_4WeeksLinegraph.html
	*/
	parameters = []interface{}{getmaxdaysfromargs(args, "conf_stat_raw_PerDay_hits")}
	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "NB RAW HITS",
	}
	genstats(args, "valid raw hits per day over the last "+strconv.Itoa(args.Outputs.Max_number_of_days)+" days", "conf_stat_raw_PerDay_hits", "stmt_raw_PerDay_hits", parameters, tableheaders, xaxisfields, valuefield, "unique hits", true)

	/*
		stat: unique 2xx and 3xx hits per hour over the last 31 (default) days
		expecting 2 htmls:
		unique_PerHour_hits_table.html
		unique_PerHour_hits_linegraph.html
	*/
	parameters = []interface{}{getmaxdaysfromargs(args, "conf_stat_unique_PerHour_hits")}
	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "HOUR",
		"Title_5": "NB UNIQUE HITS",
	}
	xaxisfields = []int{0, 1, 2, 3}
	valuefield = 4
	genstats(args, "valid unique hits per hour over the last "+strconv.Itoa(args.Outputs.Max_number_of_days)+" days", "conf_stat_unique_PerHour_hits", "stmt_unique_PerHour_hits", parameters, tableheaders, xaxisfields, valuefield, "unique hits", true)

	/*
		stat: raw 2xx and 3xx hits per hour over the last 31 (default) days
		expecting 2 htmls:
		raw_PerHour_hits_table.html
		raw_PerHour_hits_linegraph.html
	*/
	parameters = []interface{}{getmaxdaysfromargs(args, "conf_stat_raw_PerHour_hits")}
	tableheaders = map[string]string{
		"Title_1": "YEAR",
		"Title_2": "MONTH",
		"Title_3": "DAY",
		"Title_4": "HOUR",
		"Title_5": "NB RAW HITS",
	}
	genstats(args, "valid raw hits per hour over the last "+strconv.Itoa(args.Outputs.Max_number_of_days)+" days", "conf_stat_raw_PerHour_hits", "stmt_raw_PerHour_hits", parameters, tableheaders, xaxisfields, valuefield, "raw hits", true)

	/*
		stat: sum of raw hit per referrer over the last 31 (default) days
		expecting 1 htmls:
		raw_PerHour_ReferringUrls_table.html
	*/
	parameters = []interface{}{mintimestamp_Number_of_days_detailed, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of raw hits",
	}
	xaxisfields = []int{0}
	valuefield = 1
	genstats(args, "sum of raw hits per referrer over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_raw_PerHour_ReferringUrls", "stmt_raw_PerHour_ReferringUrls", parameters, tableheaders, xaxisfields, valuefield, "", false)

	/*
		stat: sum of unique hits per referrer over the last 31 (default) days
		expecting 1 htmls:
		unique_PerHour_ReferringUrls_table.html
	*/
	parameters = []interface{}{mintimestamp_Number_of_days_detailed, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of unique hits",
	}
	genstats(args, "sum of unique hits per referrer over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_PerHour_ReferringUrls", "stmt_unique_PerHour_ReferringUrls", parameters, tableheaders, xaxisfields, valuefield, "", false)

	/*
		stat: sum of unique hits per referrer, non empty non self, over the last 31 (default) days
		expecting 1 htmls:
		nunique_PerHour_RefferingUrlsNoEmptyOrSelf_table.html
	*/
	parameters = []interface{}{mintimestamp_Number_of_days_detailed, args.Generals.Mydomain, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of unique non self, non empty hits",
	}
	genstats(args, "sum of unique hits per referrer, non self non empty, over last "+strconv.Itoa(args.Outputs.Max_number_of_days)+" days", "conf_stat_unique_PerHour_RefferingUrlsNoEmptyOrSelf", "stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf", parameters, tableheaders, xaxisfields, valuefield, "", false)

	/*
		stat: sum of unique hits per referrer, non empty, non self, only tld over the last 31 (default) days
		expecting 1 htmls:
		unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld_table.html
	*/
	parameters = []interface{}{mintimestamp_Number_of_days_detailed, args.Generals.Mydomain, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "REFERRER",
		"Title_2": "NB of unique non self, non empty hits. TLDs",
	}
	genstats(args, "sum of unique hits per referrer, non self non empty, only tlds, over last "+strconv.Itoa(args.Outputs.Max_number_of_days)+" days", "conf_stat_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld", "stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld", parameters, tableheaders, xaxisfields, valuefield, "", false)
	logger("finished the function to generate statistics")

	/*
		stat: sum of raw hits per search engine over the last 31 (default) days
		expecting 1 htmls:
		raw_XDaysTotal_HitsFromSearchEngines_table.html
	*/
	parameters = []interface{}{mintimestamp_Number_of_days_detailed, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "search egine",
		"Title_2": "NB of raw se hits",
	}
	genstats(args, "sum of raw hits per search egine, over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_raw_XDaysTotal_HitsFromSearchEngines", "stmt_stat_raw_XDaysTotal_HitsFromSearchEngines", parameters, tableheaders, xaxisfields, valuefield, "", false)
	logger("finished the function to generate statistics")

	/*
		stat: sum of unique hits per search engine over the last 31 (default) days
		expecting 1 htmls:
		unique_XDaysTotal_HitsFromSearchEngines_table.html
	*/
	parameters = []interface{}{mintimestamp_Number_of_days_detailed, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "search egine",
		"Title_2": "NB of unique se hits",
	}
	genstats(args, "sum of unique hits per search egine, over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_XDaysTotal_HitsFromSearchEngines", "stmt_unique_XDaysTotal_HitsFromSearchEngines", parameters, tableheaders, xaxisfields, valuefield, "", false)
	logger("finished the function to generate statistics")

	/*
		stat: count number of times a page was used as an entry page over the last 31 (default) days
		expecting 1 htmls:
		unique_XDaysTotal_Entrypages_table.html
	*/
	parameters = []interface{}{mintimestamp_Number_of_days_detailed, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "entry page",
		"Title_2": "count",
	}
	genstats(args, "sum of times page is used as entry page, over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_XDaysTotal_Entrypages", "stmt_unique_XDaysTotal_Entrypages", parameters, tableheaders, xaxisfields, valuefield, "", false)

	/*
		stat: count number of times a page was used as an exit page over the last 31 (default) days
		expecting 1 htmls:
		unique_XDaysTotal_Exitpages_table.html
	*/
	parameters = []interface{}{mintimestamp_Number_of_days_detailed, int(args.Outputs.Numberofreferrers)}
	tableheaders = map[string]string{
		"Title_1": "exit page",
		"Title_2": "count",
	}
	genstats(args, "sum of times page is used as exit page, over last "+strconv.Itoa(args.Outputs.Number_of_days_detailed)+" days", "conf_stat_unique_XDaysTotal_Exitpages", "stmt_unique_XDaysTotal_Exitpages", parameters, tableheaders, xaxisfields, valuefield, "", false)
	logger("finished the function to generate statistics")
}
