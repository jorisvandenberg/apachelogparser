package main

import (
	"fmt"
)

func main() {
	args := getargs()
	db := createdb(args.dbpath)
	defer db.Close()
	tx := initialisedb(db)
	loadquerydb(tx)
	if (args.truncatealreadyloaded == true) {
		truncatealreadyloaded()
	}
	if (args.runtype == "all" || args.runtype == "onlylogparse") {
		parselogs(args)
	}
	err := tx.Commit()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	demobarchart(args) 
}