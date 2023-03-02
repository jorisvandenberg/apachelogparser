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
	if args.truncatealreadyloaded == true {
		truncatealreadyloaded()
	}
	if args.emptyoutputpath == true {
		emptydir(args.outputpath, ".html")
	}
	if args.runtype == "all" || args.runtype == "onlylogparse" {
		parselogs(args)
	}
	err := tx.Commit()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	demobarchart(args)
	demotable(args)
	demolinegraph(args)
	demoboxplot(args)
	demopiechart(args)
	demowritemulti(args)
	demowritehtmlpage(args)
	createindex(args)
	writelog(args)
}
