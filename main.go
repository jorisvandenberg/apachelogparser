package main

import (
	"fmt"
)

func main() {
	args := getargs()
	db := createdb(args.dbpath)
	defer db.Close()
	tx := initialisedb(db)
	

	err := tx.Commit()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("%+v", db)
}