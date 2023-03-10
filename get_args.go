package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

type Tableconfig struct {
	Table_enabled     bool
	Table_title       string
	Table_description string
	Table_pagecontent []string
	Table_pagefooter  string
	Table_filename    string
	Table_index_name  string
	Table_index_group string
	Table_index_order int
}

type Linegraphconfig struct {
	Linegraph_enabled     bool
	Linegraph_title       string
	Linegraph_description string
	Linegraph_filename    string
	Linegraph_index_group string
	Linegraph_index_order int
}

type Linegraph4weekconfig struct {
	Linegraph_compare4weeks_enabled     bool
	Linegraph_compare4weeks_title       string
	Linegraph_compare4weeks_description string
	Linegraph_compare4weeks_filename    string
	Linegraph_compare4weeks_index_group string
	Linegraph_compare4weeks_index_order int
}

type Statconfig struct {
	Statname           string
	Tableinfo          Tableconfig
	Linegraphinfo      Linegraphconfig
	Linegraph4weekinfo Linegraph4weekconfig
}

type Inputarg struct {
	Logfilepath             string
	Logfileregex            string
	Parseregex              string
	Parserfield_ip          int
	Parserfield_datetime    int
	Parserfield_method      int
	Parserfield_request     int
	Parserfield_httpversion int
	Parserfield_returncode  int
	Parserfield_httpsize    int
	Parserfield_referrer    int
	Parserfield_useragent   int
}

type Output struct {
	Outputpath              string
	Emptyoutputpath         bool
	Number_of_days_detailed int
	Assethost               string
	Zipoutput               bool
	Zippath                 string
	Numberofreferrers       int
}

type General struct {
	Dbpath     string
	Timeformat string
	Mydomain   string
	Writelog   bool
}

type Commandline struct {
	Runtype               string
	Truncatealreadyloaded bool
	Demographs            bool
	Debug                 bool
}

type Args struct {
	Inputargs         Inputarg
	Outputs           Output
	Generals          General
	Commandlines      Commandline
	Ignoredips        []string
	Ignoredhostagents []string
	Ignoredreferrers  []string
	Ignoredrequests   []string
	Stats             []Statconfig
}

func getargs() Args {
	var returndb Args
	var inputargs Inputarg
	var outputs Output
	var generals General
	var commandlines Commandline
	var mystats []Statconfig
	/*
		start command line flags input
	*/
	runtypePtr := flag.String("runtype", `all`, "options: all, onlylogparse, onlystats. Default: all")
	customconfigPtr := flag.String("config", `default`, "the full path to a custom configfile")
	truncatealreadyloadedPtr := flag.Bool("truncatealreadyloaded", false, "if set, the \"alreadyloaded\" table will be truncated if combined with runtype all or onlylogparse")
	demographsPtr := flag.Bool("demographs", false, "write a bunch of demographs to the output dir")
	debugPtr := flag.Bool("debug", false, "enable or disable debug (verbose)")
	flag.Parse()
	commandlines.Runtype = *runtypePtr
	commandlines.Truncatealreadyloaded = *truncatealreadyloadedPtr
	commandlines.Demographs = *demographsPtr
	commandlines.Debug = *debugPtr
	/*
		end command line flags input
	*/

	/*
		start config file selection
	*/
	var configfilepath string
	var logconfig string
	if *customconfigPtr != `default` {
		configfilepath = *customconfigPtr
		logconfig = "path for config file was added as a parameter, using that one: " + configfilepath
	} else {
		if _, err := os.Stat("config.ini"); err == nil {
			logconfig = "found a config.ini file in the current path... using that one"
			configfilepath = "config.ini"
		} else if _, err := os.Stat("/etc/apachelogparser/config.ini"); err == nil {
			logconfig = "found a config.ini file: /etc/apachelogparser/config.ini... using that one"
			configfilepath = "/etc/apachelogparser/config.ini"
		} else {
			os.Exit(1)
		}
	}
	cfg, err := ini.Load(configfilepath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	/*
		end config file selection
	*/

	/*
		start ignore list creation
	*/
	var ignorevisitorips_list []string
	var ignorehostagents_list []string
	var ignoredreferrers_list []string
	var ignoredrequests_list []string
	for _, ignoredip := range cfg.Section("ignorevisitorips").Keys() {
		ignorevisitorips_list = append(ignorevisitorips_list, ignoredip.String())
	}
	returndb.Ignoredips = ignorevisitorips_list

	for _, ignoredhostagent := range cfg.Section("ignorehostagents").Keys() {
		ignorehostagents_list = append(ignorehostagents_list, ignoredhostagent.String())
	}
	returndb.Ignoredhostagents = ignorehostagents_list

	for _, ignoredreferrer := range cfg.Section("ignorereferrers").Keys() {
		ignoredreferrers_list = append(ignoredreferrers_list, ignoredreferrer.String())
	}
	returndb.Ignoredreferrers = ignoredreferrers_list

	for _, ignoredrequest := range cfg.Section("ignoredrequests").Keys() {
		ignoredrequests_list = append(ignoredrequests_list, ignoredrequest.String())
	}
	returndb.Ignoredrequests = ignoredrequests_list

	/*
		end ignore list creation
	*/

	/*
		start input gathering
	*/
	inputargs.Logfilepath = cfg.Section("input").Key("logfilepath").String()
	inputargs.Logfileregex = cfg.Section("input").Key("logfileregex").String()
	inputargs.Parseregex = cfg.Section("input").Key("parseregex").String()
	switch inputargs.Parseregex {
	case "clf":
		inputargs.Parseregex = `(?m)^(\S*).*\[(.*)\]\s"(\S*)\s(\S*)\s([^"]*)"\s(\S*)\s(\S*)\s"([^"]*)"\s"([^"]*)"$`
		//case "other":
	}
	inputargs.Parserfield_ip, _ = cfg.Section("input").Key("parserfield_ip").Int()
	inputargs.Parserfield_datetime, _ = cfg.Section("input").Key("parserfield_datetime").Int()
	inputargs.Parserfield_method, _ = cfg.Section("input").Key("parserfield_method").Int()
	inputargs.Parserfield_request, _ = cfg.Section("input").Key("parserfield_request").Int()
	inputargs.Parserfield_httpversion, _ = cfg.Section("input").Key("parserfield_httpversion").Int()
	inputargs.Parserfield_returncode, _ = cfg.Section("input").Key("parserfield_returncode").Int()
	inputargs.Parserfield_httpsize, _ = cfg.Section("input").Key("parserfield_httpsize").Int()
	inputargs.Parserfield_referrer, _ = cfg.Section("input").Key("parserfield_referrer").Int()
	inputargs.Parserfield_useragent, _ = cfg.Section("input").Key("parserfield_useragent").Int()
	/*
		end input gathering
	*/

	/*
		start output gathering
	*/
	outputs.Outputpath = cfg.Section("output").Key("outputpath").String()
	outputs.Assethost = cfg.Section("output").Key("assethost").String()
	outputs.Number_of_days_detailed, _ = cfg.Section("output").Key("number_of_days_detailed").Int()
	outputs.Emptyoutputpath, _ = cfg.Section("output").Key("emptyoutputpath").Bool()
	outputs.Zipoutput, _ = cfg.Section("output").Key("zipoutput").Bool()
	outputs.Zippath = cfg.Section("output").Key("zippath").String()
	outputs.Numberofreferrers, _ = cfg.Section("output").Key("numberofreferrers").Int()
	/*
		end output gathering
	*/

	/*
		start general config gathering
	*/
	generals.Dbpath = cfg.Section("general").Key("dbpath").String()
	generals.Timeformat = cfg.Section("general").Key("timeformat").String()
	generals.Mydomain = cfg.Section("general").Key("mydomain").String()
	generals.Writelog, _ = cfg.Section("general").Key("writelog").Bool()
	/*
		end general config gathering
	*/

	/*
		start stats config secion stat_perhour_hits_raw_2xx_3xx
	*/
	stat_enabled, _ := cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("enabled").Bool()
	table_enabled, _ := cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_enabled").Bool()
	linegraph_enabled, _ := cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("linegraph_enabled").Bool()

	if stat_enabled && (table_enabled || linegraph_enabled) {
		var mystatconfig Statconfig
		var mytableconfig Tableconfig
		var mylinegraphconfig Linegraphconfig

		mystatconfig.Statname = "stat_perhour_hits_raw_2xx_3xx"
		if table_enabled {
			mytableconfig.Table_enabled = true
			mytableconfig.Table_title = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_title").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_description = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_description").String(), outputs.Number_of_days_detailed)
			tablecontent_unsplitstring := splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_pagecontent").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_pagecontent = strings.Split(tablecontent_unsplitstring, "|")
			mytableconfig.Table_pagefooter = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_pagefooter").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_filename = cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_filename").String()
			mytableconfig.Table_index_name = cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_index_name").String()
			mytableconfig.Table_index_group = cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_index_group").String()
			mytableconfig.Table_index_order, _ = cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("table_index_order").Int()

		} else {
			mytableconfig.Table_enabled = false
		}

		if linegraph_enabled {
			mylinegraphconfig.Linegraph_enabled = true
			mylinegraphconfig.Linegraph_title = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("linegraph_title").String(), outputs.Number_of_days_detailed)
			mylinegraphconfig.Linegraph_description = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("linegraph_description").String(), outputs.Number_of_days_detailed)
			mylinegraphconfig.Linegraph_filename = cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("linegraph_filename").String()
			mylinegraphconfig.Linegraph_index_group = cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("linegraph_index_group").String()
			mylinegraphconfig.Linegraph_index_order, _ = cfg.Section("stat_perhour_hits_raw_2xx_3xx").Key("linegraph_index_order").Int()

		} else {
			mylinegraphconfig.Linegraph_enabled = false
		}

		mystatconfig.Tableinfo = mytableconfig
		mystatconfig.Linegraphinfo = mylinegraphconfig
		mystats = append(mystats, mystatconfig)
	}

	/*
		end stats config secion stat_perhour_hits_raw_2xx_3xx
	*/

	/*
		start stats config secion stat_perday_hits_raw_2xx_3xx
	*/
	stat_enabled, _ = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("enabled").Bool()
	table_enabled, _ = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_enabled").Bool()
	linegraph_enabled, _ = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_enabled").Bool()
	linegraph_compare4weeks_enabled, _ := cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_compare4weeks_enabled").Bool()
	if stat_enabled && (table_enabled || linegraph_enabled || linegraph_compare4weeks_enabled) {
		var mystatconfig Statconfig
		var mytableconfig Tableconfig
		var mylinegraphconfig Linegraphconfig
		var mylinegraph4weekconfig Linegraph4weekconfig
		mystatconfig.Statname = "stat_perday_hits_raw_2xx_3xx"
		if table_enabled {
			mytableconfig.Table_enabled = true
			mytableconfig.Table_title = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_title").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_description = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_description").String(), outputs.Number_of_days_detailed)
			tablecontent_unsplitstring := splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_pagecontent").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_pagecontent = strings.Split(tablecontent_unsplitstring, "|")
			mytableconfig.Table_pagefooter = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_pagefooter").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_filename = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_filename").String()
			mytableconfig.Table_index_name = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_index_name").String()
			mytableconfig.Table_index_group = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_index_group").String()
			mytableconfig.Table_index_order, _ = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("table_index_order").Int()

		} else {
			mytableconfig.Table_enabled = false
		}

		if linegraph_enabled {
			mylinegraphconfig.Linegraph_enabled = true
			mylinegraphconfig.Linegraph_title = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_title").String(), outputs.Number_of_days_detailed)
			mylinegraphconfig.Linegraph_description = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_description").String(), outputs.Number_of_days_detailed)
			mylinegraphconfig.Linegraph_filename = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_filename").String()
			mylinegraphconfig.Linegraph_index_group = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_index_group").String()
			mylinegraphconfig.Linegraph_index_order, _ = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_index_order").Int()

		} else {
			mylinegraphconfig.Linegraph_enabled = false
		}
		if linegraph_compare4weeks_enabled {
			mylinegraph4weekconfig.Linegraph_compare4weeks_enabled = true
			mylinegraph4weekconfig.Linegraph_compare4weeks_title = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_compare4weeks_title").String(), outputs.Number_of_days_detailed)
			mylinegraph4weekconfig.Linegraph_compare4weeks_description = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_compare4weeks_description").String(), outputs.Number_of_days_detailed)
			mylinegraph4weekconfig.Linegraph_compare4weeks_filename = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_compare4weeks_filename").String()
			mylinegraph4weekconfig.Linegraph_compare4weeks_index_group = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_compare4weeks_index_group").String()
			mylinegraph4weekconfig.Linegraph_compare4weeks_index_order, _ = cfg.Section("stat_perday_hits_raw_2xx_3xx").Key("linegraph_compare4weeks_index_order").Int()

		} else {
			mylinegraph4weekconfig.Linegraph_compare4weeks_enabled = false
		}
		mystatconfig.Tableinfo = mytableconfig
		mystatconfig.Linegraphinfo = mylinegraphconfig
		mystatconfig.Linegraph4weekinfo = mylinegraph4weekconfig
		mystats = append(mystats, mystatconfig)
	}

	/*
		end stats config secion stat_perday_hits_raw_2xx_3xx
	*/

	/*
		start stats config secion stat_perhour_hits_unique_2xx_3xx
	*/
	stat_enabled, _ = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("enabled").Bool()
	table_enabled, _ = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_enabled").Bool()
	linegraph_enabled, _ = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("linegraph_enabled").Bool()
	if stat_enabled && (table_enabled || linegraph_enabled) {
		var mystatconfig Statconfig
		var mytableconfig Tableconfig
		var mylinegraphconfig Linegraphconfig
		mystatconfig.Statname = "stat_perhour_hits_unique_2xx_3xx"
		if table_enabled {
			mytableconfig.Table_enabled = true
			mytableconfig.Table_title = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_title").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_description = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_description").String(), outputs.Number_of_days_detailed)
			tablecontent_unsplitstring := splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_pagecontent").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_pagecontent = strings.Split(tablecontent_unsplitstring, "|")
			mytableconfig.Table_pagefooter = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_pagefooter").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_filename = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_filename").String()
			mytableconfig.Table_index_name = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_index_name").String()
			mytableconfig.Table_index_group = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_index_group").String()
			mytableconfig.Table_index_order, _ = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("table_index_order").Int()

		} else {
			mytableconfig.Table_enabled = false
		}

		if linegraph_enabled {
			mylinegraphconfig.Linegraph_enabled = true
			mylinegraphconfig.Linegraph_title = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("linegraph_title").String(), outputs.Number_of_days_detailed)
			mylinegraphconfig.Linegraph_description = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("linegraph_description").String(), outputs.Number_of_days_detailed)
			mylinegraphconfig.Linegraph_filename = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("linegraph_filename").String()
			mylinegraphconfig.Linegraph_index_group = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("linegraph_index_group").String()
			mylinegraphconfig.Linegraph_index_order, _ = cfg.Section("stat_perhour_hits_unique_2xx_3xx").Key("linegraph_index_order").Int()

		} else {
			mylinegraphconfig.Linegraph_enabled = false
		}
		mystatconfig.Tableinfo = mytableconfig
		mystatconfig.Linegraphinfo = mylinegraphconfig
		mystats = append(mystats, mystatconfig)
	}

	/*
		end stats config secion stat_perhour_hits_unique_2xx_3xx
	*/

	/*
		start stats config secion stat_perday_unique_raw_2xx_3xx
	*/
	stat_enabled, _ = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("enabled").Bool()
	table_enabled, _ = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_enabled").Bool()
	linegraph_enabled, _ = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("linegraph_enabled").Bool()
	if stat_enabled && (table_enabled || linegraph_enabled) {
		var mystatconfig Statconfig
		var mytableconfig Tableconfig
		var mylinegraphconfig Linegraphconfig
		mystatconfig.Statname = "stat_perday_hits_unique_2xx_3xx"
		if table_enabled {
			mytableconfig.Table_enabled = true
			mytableconfig.Table_title = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_title").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_description = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_description").String(), outputs.Number_of_days_detailed)
			tablecontent_unsplitstring := splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_pagecontent").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_pagecontent = strings.Split(tablecontent_unsplitstring, "|")
			mytableconfig.Table_pagefooter = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_pagefooter").String(), outputs.Number_of_days_detailed)
			mytableconfig.Table_filename = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_filename").String()
			mytableconfig.Table_index_name = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_index_name").String()
			mytableconfig.Table_index_group = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_index_group").String()
			mytableconfig.Table_index_order, _ = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("table_index_order").Int()

		} else {
			mytableconfig.Table_enabled = false
		}

		if linegraph_enabled {
			mylinegraphconfig.Linegraph_enabled = true
			mylinegraphconfig.Linegraph_title = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("linegraph_title").String(), outputs.Number_of_days_detailed)
			mylinegraphconfig.Linegraph_description = splice_number_of_days_detailed_in(cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("linegraph_description").String(), outputs.Number_of_days_detailed)
			mylinegraphconfig.Linegraph_filename = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("linegraph_filename").String()
			mylinegraphconfig.Linegraph_index_group = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("linegraph_index_group").String()
			mylinegraphconfig.Linegraph_index_order, _ = cfg.Section("stat_perday_hits_unique_2xx_3xx").Key("linegraph_index_order").Int()

		} else {
			mylinegraphconfig.Linegraph_enabled = false
		}
		mystatconfig.Tableinfo = mytableconfig
		mystatconfig.Linegraphinfo = mylinegraphconfig
		mystats = append(mystats, mystatconfig)
	}

	/*
		end stats config secion stat_perday_hits_raw_2xx_3xx
	*/

	/*
		start stats config secion stat_perhour_referrers_raw_2xx_3xx
	*/
	stat_enabled, _ = cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("enabled").Bool()
	table_enabled, _ = cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_enabled").Bool()
	linegraph_enabled, _ = cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("linegraph_enabled").Bool()
	if stat_enabled && table_enabled {
		var mystatconfig Statconfig
		var mytableconfig Tableconfig
		var mylinegraphconfig Linegraphconfig
		mystatconfig.Statname = "stat_perhour_referrers_raw_2xx_3xx"
		mylinegraphconfig.Linegraph_enabled = false
		mytableconfig.Table_enabled = true
		mytableconfig.Table_title = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_title").String(), outputs.Number_of_days_detailed)
		mytableconfig.Table_description = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_description").String(), outputs.Number_of_days_detailed)
		tablecontent_unsplitstring := splice_number_of_days_detailed_in(cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_pagecontent").String(), outputs.Number_of_days_detailed)
		mytableconfig.Table_pagecontent = strings.Split(tablecontent_unsplitstring, "|")
		mytableconfig.Table_pagefooter = splice_number_of_days_detailed_in(cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_pagefooter").String(), outputs.Number_of_days_detailed)
		mytableconfig.Table_filename = cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_filename").String()
		mytableconfig.Table_index_name = cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_index_name").String()
		mytableconfig.Table_index_group = cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_index_group").String()
		mytableconfig.Table_index_order, _ = cfg.Section("stat_perhour_referrers_raw_2xx_3xx").Key("table_index_order").Int()

		mystatconfig.Tableinfo = mytableconfig
		mystatconfig.Linegraphinfo = mylinegraphconfig
		mystats = append(mystats, mystatconfig)
	}

	/*
		end stats config secion stat_perhour_referrers_raw_2xx_3xx
	*/
	/*
		start fill struct, log and return the args
	*/
	returndb.Inputargs = inputargs
	returndb.Outputs = outputs
	returndb.Generals = generals
	returndb.Commandlines = commandlines
	returndb.Stats = mystats
	logger(logconfig)

	returndb_json, err := json.MarshalIndent(returndb, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}
	if commandlines.Debug {
		fmt.Println(string(returndb_json))
	}

	logger("config (json): " + string(returndb_json))
	return returndb
	/*
		end fill struct, log and return the args
	*/
}
