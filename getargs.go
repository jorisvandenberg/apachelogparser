package main

import (
	"flag"
	"os"
	"fmt"
	"gopkg.in/ini.v1"
)

type args struct {
	runtype string
	logfilepath string
	outputpath string
	logfileregex string
	dbpath string
	timeformat string 
	mydomain          string
	ignoredips        []string
	ignoredhostagents []string
	ignoredreferrers  []string
	ignoredrequests   []string
	number_of_days_detailed  int
	number_of_days_per_hour  int
	number_of_days_per_day   int
	number_of_days_per_week  int
	number_of_days_per_month int
	truncatealreadyloaded bool
}

func getargs() args {
	var output args
	runtypePtr := flag.String("runtype", `all`, "options: all, onlylogparse, onlystats. Default: all")
	truncatealreadyloadedPtr := flag.Bool("truncatealreadyloaded", false, "if set, the \"alreadyloaded\" table will be truncated if combined with runtype all or onlylogparse")
	flag.Parse()
	output.runtype = *runtypePtr
	output.truncatealreadyloaded = *truncatealreadyloadedPtr
	var configfilepath string
	if _, err := os.Stat("config.ini"); err == nil {
		fmt.Printf("found a config.ini file in the current path... using that one\n")
		configfilepath = "config.ini"
	} else if _, err := os.Stat("/etc/apachelogparser/config.ini"); err == nil {
		fmt.Printf("found a config.ini file: /etc/apachelogparser/config.ini... using that one\n")
		configfilepath = "/etc/apachelogparser/config.ini"
	} else {
		os.Exit(1)
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

	output.outputpath = cfg.Section("output").Key("outputpath").String()
	output.number_of_days_detailed, _ = cfg.Section("output").Key("number_of_days_detailed").Int()
	output.number_of_days_per_hour, _ = cfg.Section("output").Key("number_of_days_per_hour").Int()
	output.number_of_days_per_day, _ = cfg.Section("output").Key("number_of_days_per_day").Int()
	output.number_of_days_per_week, _ = cfg.Section("output").Key("number_of_days_per_week").Int()
	output.number_of_days_per_month, _ = cfg.Section("output").Key("number_of_days_per_month").Int()

	output.dbpath = cfg.Section("general").Key("dbpath").String()
	output.timeformat = cfg.Section("general").Key("timeformat").String()
	output.mydomain = cfg.Section("general").Key("mydomain").String()
	return output
}