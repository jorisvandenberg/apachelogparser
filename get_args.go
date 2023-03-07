package main

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"time"
)

type args struct {
	runtype                  string
	logfilepath              string
	outputpath               string
	logfileregex             string
	dbpath                   string
	timeformat               string
	mydomain                 string
	parseregex               string
	assethost                string
	ignoredips               []string
	ignoredhostagents        []string
	ignoredreferrers         []string
	ignoredrequests          []string
	number_of_days_detailed  int
	number_of_days_per_hour  int
	number_of_days_per_day   int
	number_of_days_per_week  int
	number_of_days_per_month int
	parserfield_ip           int
	parserfield_datetime     int
	parserfield_method       int
	parserfield_request      int
	parserfield_httpversion  int
	parserfield_returncode   int
	parserfield_httpsize     int
	parserfield_referrer     int
	parserfield_useragent    int
	numberofreferrers	     int
	truncatealreadyloaded    bool
	emptyoutputpath          bool
	writelog                 bool
	demographs               bool
	zipoutput                bool
	zippath                  string
}

func getargs() args {
	var output args
	runtypePtr := flag.String("runtype", `all`, "options: all, onlylogparse, onlystats. Default: all")
	customconfigPtr := flag.String("config", `default`, "the full path to a custom configfile")
	truncatealreadyloadedPtr := flag.Bool("truncatealreadyloaded", false, "if set, the \"alreadyloaded\" table will be truncated if combined with runtype all or onlylogparse")
	demographsPtr := flag.Bool("demographs", false, "write a bunch of demographs to the output dir")
	flag.Parse()
	output.runtype = *runtypePtr
	output.truncatealreadyloaded = *truncatealreadyloadedPtr
	output.demographs = *demographsPtr

	var configfilepath string
	var logconfig string
	t := time.Now()
	if *customconfigPtr != `default` {
		configfilepath = *customconfigPtr
		logconfig = t.Format("2006-01-02 15:04:05") + " => path for config file was added as a parameter, using that one: " + configfilepath
	} else {
		if _, err := os.Stat("config.ini"); err == nil {
			logconfig = t.Format("2006-01-02 15:04:05") + " => found a config.ini file in the current path... using that one"
			configfilepath = "config.ini"
		} else if _, err := os.Stat("/etc/apachelogparser/config.ini"); err == nil {
			logconfig = t.Format("2006-01-02 15:04:05") + " => found a config.ini file: /etc/apachelogparser/config.ini... using that one"
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
	output.ignoredips = ignorevisitorips_list

	for _, ignoredhostagent := range cfg.Section("ignorehostagents").Keys() {
		ignorehostagents_list = append(ignorehostagents_list, ignoredhostagent.String())
	}
	output.ignoredhostagents = ignorehostagents_list

	for _, ignoredreferrer := range cfg.Section("ignorereferrers").Keys() {
		ignoredreferrers_list = append(ignoredreferrers_list, ignoredreferrer.String())
	}
	output.ignoredreferrers = ignoredreferrers_list

	for _, ignoredrequest := range cfg.Section("ignoredrequests").Keys() {
		ignoredrequests_list = append(ignoredrequests_list, ignoredrequest.String())
	}
	output.ignoredrequests = ignoredrequests_list

	output.logfilepath = cfg.Section("input").Key("logfilepath").String()
	output.logfileregex = cfg.Section("input").Key("logfileregex").String()
	output.parseregex = cfg.Section("input").Key("parseregex").String()
	output.parserfield_ip, _ = cfg.Section("input").Key("parserfield_ip").Int()
	output.parserfield_datetime, _ = cfg.Section("input").Key("parserfield_datetime").Int()
	output.parserfield_method, _ = cfg.Section("input").Key("parserfield_method").Int()
	output.parserfield_request, _ = cfg.Section("input").Key("parserfield_request").Int()
	output.parserfield_httpversion, _ = cfg.Section("input").Key("parserfield_httpversion").Int()
	output.parserfield_returncode, _ = cfg.Section("input").Key("parserfield_returncode").Int()
	output.parserfield_httpsize, _ = cfg.Section("input").Key("parserfield_httpsize").Int()
	output.parserfield_referrer, _ = cfg.Section("input").Key("parserfield_referrer").Int()
	output.parserfield_useragent, _ = cfg.Section("input").Key("parserfield_useragent").Int()

	output.outputpath = cfg.Section("output").Key("outputpath").String()
	output.assethost = cfg.Section("output").Key("assethost").String()
	output.number_of_days_detailed, _ = cfg.Section("output").Key("number_of_days_detailed").Int()
	output.number_of_days_per_hour, _ = cfg.Section("output").Key("number_of_days_per_hour").Int()
	output.number_of_days_per_day, _ = cfg.Section("output").Key("number_of_days_per_day").Int()
	output.number_of_days_per_week, _ = cfg.Section("output").Key("number_of_days_per_week").Int()
	output.number_of_days_per_month, _ = cfg.Section("output").Key("number_of_days_per_month").Int()
	output.emptyoutputpath, _ = cfg.Section("output").Key("emptyoutputpath").Bool()
	output.zipoutput, _ = cfg.Section("output").Key("zipoutput").Bool()
	output.zippath = cfg.Section("output").Key("zippath").String()
	output.numberofreferrers = cfg.Section("output").Key("numberofreferrers").String()

	output.dbpath = cfg.Section("general").Key("dbpath").String()
	output.timeformat = cfg.Section("general").Key("timeformat").String()
	output.mydomain = cfg.Section("general").Key("mydomain").String()
	output.writelog, _ = cfg.Section("general").Key("writelog").Bool()

	if output.writelog {
		mylog = append(mylog, logconfig)
	}
	return output
}
