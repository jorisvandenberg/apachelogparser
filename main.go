package main

import (
	"fmt"
)

var mylog []string

func main() {
	args := getargs()
	db := createdb(args.dbpath)
	defer db.Close()
	tx := initialisedb(db)
	loadquerydb(tx)
	filltemplatedb()
	if args.truncatealreadyloaded {
		truncatealreadyloaded()
	}
	if args.emptyoutputpath {
		emptydir(args.outputpath, ".html")
	}
	if args.runtype == "all" || args.runtype == "onlylogparse" {
		parselogs(args)
	} 
	if args.runtype == "all" || args.runtype == "onlystats" {
		generatestats(args)
	}
	err := tx.Commit()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	if args.demographs {
		writedemographs(args)
	}
	
	writelog(args)
	createindex(args)
	
	if args.zipoutput {
		ZipWriter(args.outputpath, args.zippath)
	}
}
