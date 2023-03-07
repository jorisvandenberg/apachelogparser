package main

import (
	"fmt"
)

var mylog []string

func main() {
	logger("started the program")
	logger("fetching the arguments from the ini file and the commandline")
	args := getargs()
	logger("finished the arguments from the ini file and the commandline")
	logger("started the db initialisation en query loading")
	db := createdb(args.dbpath)
	defer db.Close()
	tx := initialisedb(db)
	loadquerydb(tx)
	logger("finished the db initialisation en query loading")
	logger("fetching all the html templates")
	filltemplatedb()
	logger("finished fetching all the html templates")
	if args.truncatealreadyloaded {
		logger("received the command line parameter to truncate the already loaded table hashes...")
		truncatealreadyloaded()
		logger("finished truncating the already loaded table hashes")
	} else {
		logger("the program wasn't asked to truncate the already loaded table hashes, so previously loaded tables will be skipped")
	}
	if args.outputs.emptyoutputpath {
		logger("i was asked to remove the html files generated by the previous run from the output path")
		emptydir(args.outputs.outputpath, ".html")
		logger("finished removing the html files generated by the previous run from the output path")
	} else {
		logger("the option to clean the output path before each run was deactivated. Old html files will not be cleaned (they can be overwritten tough!)")
	}
	logger("the requested runtype (all/onlylogparse/onlystats) was " + args.runtype + ", acting accordingly")
	if args.runtype == "all" || args.runtype == "onlylogparse" {
		logger("starting the log parsing process")
		parselogs(args)
		logger("finishing the log parsing process")
	}
	if args.runtype == "all" || args.runtype == "onlystats" {
		logger("starting the stat generating process")
		generatestats(args)
		logger("finishing the stat generating process")
	}
	logger("commiting changes to the sqlite database")
	err := tx.Commit()
	logger("finished commiting changes to the sqlite database")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		logger("Whoops: tx.commit:" + err.Error())
	}

	if args.demographs {
		writedemographs(args)
	}
	logger("finished the program, no more logging from this point out, since the logs have to be written, included in the index and included in the zipfile (if requested), so i cannot log these functions (chicken or the egg problem)")
	writelog(args)
	createindex(args)

	if args.zipoutput {
		ZipWriter(args.outputs.outputpath, args.zippath)
	}

}
