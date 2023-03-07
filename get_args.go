package main

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

type inputarg struct {
	logfilepath              string
	logfileregex             string
	parseregex               string
	parserfield_ip           int
	parserfield_datetime     int
	parserfield_method       int
	parserfield_request      int
	parserfield_httpversion  int
	parserfield_returncode   int
	parserfield_httpsize     int
	parserfield_referrer     int
	parserfield_useragent    int
}

type output struct {
	outputpath               string
	emptyoutputpath          bool
	number_of_days_detailed  int
}

type args struct {
	inputargs 				 inputarg
	outputs				 output
	runtype                  string
	dbpath                   string
	timeformat               string
	mydomain                 string
	assethost                string
	ignoredips               []string
	ignoredhostagents        []string
	ignoredreferrers         []string
	ignoredrequests          []string
	number_of_days_per_hour  int
	number_of_days_per_day   int
	number_of_days_per_week  int
	number_of_days_per_month int
	numberofreferrers        int
	truncatealreadyloaded    bool
	writelog                 bool
	demographs               bool
	zipoutput                bool
	zippath                  string
}

func getargs() args {
	var returndb args
	var inputargs inputarg
	var outputs output
	runtypePtr := flag.String("runtype", `all`, "options: all, onlylogparse, onlystats. Default: all")
	customconfigPtr := flag.String("config", `default`, "the full path to a custom configfile")
	truncatealreadyloadedPtr := flag.Bool("truncatealreadyloaded", false, "if set, the \"alreadyloaded\" table will be truncated if combined with runtype all or onlylogparse")
	demographsPtr := flag.Bool("demographs", false, "write a bunch of demographs to the output dir")
	flag.Parse()
	returndb.runtype = *runtypePtr
	returndb.truncatealreadyloaded = *truncatealreadyloadedPtr
	returndb.demographs = *demographsPtr

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
	var ignorevisitorips_list []string
	var ignorehostagents_list []string
	var ignoredreferrers_list []string
	var ignoredrequests_list []string
	for _, ignoredip := range cfg.Section("ignorevisitorips").Keys() {
		ignorevisitorips_list = append(ignorevisitorips_list, ignoredip.String())
	}
	returndb.ignoredips = ignorevisitorips_list

	for _, ignoredhostagent := range cfg.Section("ignorehostagents").Keys() {
		ignorehostagents_list = append(ignorehostagents_list, ignoredhostagent.String())
	}
	returndb.ignoredhostagents = ignorehostagents_list

	for _, ignoredreferrer := range cfg.Section("ignorereferrers").Keys() {
		ignoredreferrers_list = append(ignoredreferrers_list, ignoredreferrer.String())
	}
	returndb.ignoredreferrers = ignoredreferrers_list

	for _, ignoredrequest := range cfg.Section("ignoredrequests").Keys() {
		ignoredrequests_list = append(ignoredrequests_list, ignoredrequest.String())
	}
	returndb.ignoredrequests = ignoredrequests_list

	inputargs.logfilepath = cfg.Section("input").Key("logfilepath").String()
	inputargs.logfileregex = cfg.Section("input").Key("logfileregex").String()
	inputargs.parseregex = cfg.Section("input").Key("parseregex").String()
	switch inputargs.parseregex {
    case "clf":
      inputargs.parseregex   = `(?m)^(\S*).*\[(.*)\]\s"(\S*)\s(\S*)\s([^"]*)"\s(\S*)\s(\S*)\s"([^"]*)"\s"([^"]*)"$`
	//case "other":
    }
	
	inputargs.parserfield_ip, _ = cfg.Section("input").Key("parserfield_ip").Int()
	inputargs.parserfield_datetime, _ = cfg.Section("input").Key("parserfield_datetime").Int()
	inputargs.parserfield_method, _ = cfg.Section("input").Key("parserfield_method").Int()
	inputargs.parserfield_request, _ = cfg.Section("input").Key("parserfield_request").Int()
	inputargs.parserfield_httpversion, _ = cfg.Section("input").Key("parserfield_httpversion").Int()
	inputargs.parserfield_returncode, _ = cfg.Section("input").Key("parserfield_returncode").Int()
	inputargs.parserfield_httpsize, _ = cfg.Section("input").Key("parserfield_httpsize").Int()
	inputargs.parserfield_referrer, _ = cfg.Section("input").Key("parserfield_referrer").Int()
	inputargs.parserfield_useragent, _ = cfg.Section("input").Key("parserfield_useragent").Int()

	outputs.outputpath = cfg.Section("output").Key("outputpath").String()
	returndb.assethost = cfg.Section("output").Key("assethost").String()
	outputs.number_of_days_detailed, _ = cfg.Section("output").Key("number_of_days_detailed").Int()
	returndb.number_of_days_per_hour, _ = cfg.Section("output").Key("number_of_days_per_hour").Int()
	returndb.number_of_days_per_day, _ = cfg.Section("output").Key("number_of_days_per_day").Int()
	returndb.number_of_days_per_week, _ = cfg.Section("output").Key("number_of_days_per_week").Int()
	returndb.number_of_days_per_month, _ = cfg.Section("output").Key("number_of_days_per_month").Int()
	outputs.emptyoutputpath, _ = cfg.Section("output").Key("emptyoutputpath").Bool()
	returndb.zipoutput, _ = cfg.Section("output").Key("zipoutput").Bool()
	returndb.zippath = cfg.Section("output").Key("zippath").String()
	returndb.numberofreferrers, _ = cfg.Section("output").Key("numberofreferrers").Int()

	returndb.dbpath = cfg.Section("general").Key("dbpath").String()
	returndb.timeformat = cfg.Section("general").Key("timeformat").String()
	returndb.mydomain = cfg.Section("general").Key("mydomain").String()
	returndb.writelog, _ = cfg.Section("general").Key("writelog").Bool()
	
	returndb.inputargs = inputargs
	returndb.outputs = outputs
	logger(logconfig)
	
	
	return returndb
}
