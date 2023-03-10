package main

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

func insertrow(ip string, datumtijd string, method string, request string, httpversion string, returncode string, httpsize string, referrer string, useragent string, maxtimestamp int, args Args) {
	longForm := args.Generals.Timeformat
	/*
		create user and return userid or return userid of existing user (userid)
	*/
	thetime, e := time.Parse(longForm, datumtijd)
	if e != nil {
		fmt.Printf("Can't parse time format")
	}
	epoch := thetime.Unix()
	if int(epoch) > maxtimestamp {
		stmt_countusers := myquerydb["stmt_countusers"].stmt
		var numberofusers int
		stmt_countusers.QueryRow(ip, useragent).Scan(&numberofusers)

		var userid int
		if numberofusers > 0 {
			//user already exists... get his id :)
			stmt_selectuserid := myquerydb["stmt_selectuserid"].stmt
			stmt_selectuserid.QueryRow(ip, useragent).Scan(&userid)
		} else {
			//user does not exist... create the bugger
			stmt_insertuser := myquerydb["stmt_insertuser"].stmt
			stmt_insertuser_result, err := stmt_insertuser.Exec(ip, useragent)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}
			var id64 int64
			id64, err = stmt_insertuser_result.LastInsertId()
			userid = int(id64)
		}

		/*
			create request and return requestid or return requestid of existing request (requestid)
		*/
		stmt_countrequest := myquerydb["stmt_countrequest"].stmt
		var numberofrequests int
		stmt_countrequest.QueryRow(request).Scan(&numberofrequests)
		var requestid int
		if numberofrequests > 0 {
			stmt_selectrequestid := myquerydb["stmt_selectrequestid"].stmt
			stmt_selectrequestid.QueryRow(request).Scan(&requestid)
		} else {
			stmt_insertrequest := myquerydb["stmt_insertrequest"].stmt

			stmt_insertrequest_result, err := stmt_insertrequest.Exec(request)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}
			var id64 int64
			id64, err = stmt_insertrequest_result.LastInsertId()
			requestid = int(id64)
		}

		/*
			create referrer and return referrerid or return referrerid of existing referrer (referrerid)
		*/
		stmt_countreferrer := myquerydb["stmt_countreferrer"].stmt
		var numberofreferrers int
		stmt_countreferrer.QueryRow(referrer).Scan(&numberofreferrers)
		var referrerid int
		if numberofreferrers > 0 {
			stmt_selectreferrerid := myquerydb["stmt_selectreferrerid"].stmt
			stmt_selectreferrerid.QueryRow(referrer).Scan(&referrerid)
		} else {
			stmt_insertreferrer := myquerydb["stmt_insertreferrer"].stmt
			stmt_insertreferrer_result, err := stmt_insertreferrer.Exec(referrer)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				os.Exit(1)
			}
			var id64 int64
			id64, err = stmt_insertreferrer_result.LastInsertId()
			referrerid = int(id64)
		}
		/*
			get max timestamp of current db and insert newer records
		*/
		stmt_insertvisit := myquerydb["stmt_insertvisit"].stmt
		_, err := stmt_insertvisit.Exec(referrerid, requestid, int(epoch), userid, returncode, httpsize)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}
}

func truncatealreadyloaded() {
	stmt_truncatealreadyloaded := myquerydb["stmt_truncatealreadyloaded"].stmt

	_, err := stmt_truncatealreadyloaded.Exec()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

func InsertParsedFileHashIntoDb(filename string, filepath string) {

	filehandle, err := os.Open(filepath + filename)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer filehandle.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, filehandle); err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	filehash := hex.EncodeToString(hash.Sum(nil))
	stmt_insertalreadyloaded := myquerydb["stmt_insertalreadyloaded"].stmt

	_, err = stmt_insertalreadyloaded.Exec(filehash)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

}

func getfiles(regex string, pathS string) []string {
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(regex, f.Name())
			if err == nil && r {
				filehandle, err := os.Open(pathS + f.Name())
				if err != nil {
					fmt.Printf("%s\n", err.Error())
				}
				defer filehandle.Close()

				hash := sha256.New()
				if _, err := io.Copy(hash, filehandle); err != nil {
					fmt.Printf("%s\n", err.Error())
				}
				filehash := hex.EncodeToString(hash.Sum(nil))
				stmt_countalreadyloaded := myquerydb["stmt_countalreadyloaded"].stmt
				var countalreadyloaded int
				stmt_countalreadyloaded.QueryRow(filehash).Scan(&countalreadyloaded)
				if countalreadyloaded == 0 {
					files = append(files, f.Name())
					logger(f.Name() + " was added to the todo list")

				} else {
					logger(f.Name() + " was already parsed in the past... skipping")
				}

			}
		}
		return nil
	})
	return files
}

func parseme(line string, maxvisittimestamp int, args Args) {
	re := regexp.MustCompile(args.Inputargs.Parseregex)
	match := re.FindStringSubmatch(line)

	if len(match) == 10 {
		ip := match[args.Inputargs.Parserfield_ip]
		datumtijd := match[args.Inputargs.Parserfield_datetime]
		method := match[args.Inputargs.Parserfield_method]
		request := match[args.Inputargs.Parserfield_request]
		httpversion := match[args.Inputargs.Parserfield_httpversion]
		returncode := match[args.Inputargs.Parserfield_returncode]
		httpsize := match[args.Inputargs.Parserfield_httpsize]
		referrer := match[args.Inputargs.Parserfield_referrer]
		useragent := match[args.Inputargs.Parserfield_useragent]
		ignore := false
		for _, ignoredhostagent := range args.Ignoredhostagents {
			r, err := regexp.MatchString(ignoredhostagent, useragent)
			if err == nil && r {
				ignore = true
			}
		}
		for _, ignoredip := range args.Ignoredips {
			r, err := regexp.MatchString(ignoredip, ip)
			if err == nil && r {
				ignore = true
			}
		}
		for _, ignoredreferrer := range args.Ignoredreferrers {
			r, err := regexp.MatchString(ignoredreferrer, referrer)
			if err == nil && r {
				ignore = true
			}
		}
		for _, ignoredrequest := range args.Ignoredrequests {
			r, err := regexp.MatchString(ignoredrequest, request)
			if err == nil && r {
				ignore = true
			}
		}
		if ignore == false {
			insertrow(ip, datumtijd, method, request, httpversion, returncode, httpsize, referrer, useragent, maxvisittimestamp, args)
			//fmt.Printf("%+v\n\n\n", ip, datumtijd, method, request, httpversion, returncode, httpsize, referrer, useragent, maxvisittimestamp, args)
		}

	} else {
		fmt.Printf("unable to parse line %d %s", len(match), line)
	}
}

func getmaxvisittimestamp() int {
	stmt_maxvisittimestamp := myquerydb["stmt_maxvisittimestamp"].stmt
	var output int
	stmt_maxvisittimestamp.QueryRow().Scan(&output)
	logger("the current latest log record had a unix timestamp of " + strconv.Itoa(output) + ", i'm not going to load records that have a timestamp older than this one!!!")
	return output
}

func parselogs(args Args) {
	maxvisittimestamp := getmaxvisittimestamp()
	toparselist := getfiles(args.Inputargs.Logfileregex, args.Inputargs.Logfilepath)
	//fmt.Printf("%+v\n", toparselist)
	var scanner *bufio.Scanner
	for _, filename := range toparselist {
		file, err := os.Open(args.Inputargs.Logfilepath + filename)
		defer file.Close()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		r, err := regexp.MatchString(`.*\.gz`, filename)
		if err == nil && r {
			gz, err := gzip.NewReader(file)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			defer gz.Close()
			scanner = bufio.NewScanner(gz)
		} else {
			scanner = bufio.NewScanner(file)
		}
		for scanner.Scan() {
			currentline := scanner.Text()
			parseme(currentline, maxvisittimestamp, args)
		}
		InsertParsedFileHashIntoDb(filename, args.Inputargs.Logfilepath)
	}
}
