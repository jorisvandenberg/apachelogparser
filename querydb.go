package main

import(
	"database/sql"
	"os"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type querydb struct {
	querynaam string
	sqlcode string
	stmt *sql.Stmt
}

var myquerydb = make(map[string]querydb)

func loadquerydb(tx *sql.Tx) {
	var myquery querydb

	//statement for inserting new visit 
	query_insertvisit := "insert into visit(referrer, request,  visit_timestamp, user, statuscode, httpsize) values (?,?,?,?,?,?)"
	stmt_insertvisit, err := tx.Prepare(query_insertvisit)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_insertvisit"
	myquery.sqlcode = query_insertvisit
	myquery.stmt = stmt_insertvisit
	myquerydb["stmt_insertvisit"] = myquery

	//statement to count if a file was loaded previously
	query_countalreadyloaded := "select count(*) from alreadyloaded where hash = ?"
	stmt_countalreadyloaded, err := tx.Prepare(query_countalreadyloaded)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_countalreadyloaded"
	myquery.sqlcode = query_countalreadyloaded
	myquery.stmt = stmt_countalreadyloaded
	myquerydb["stmt_countalreadyloaded"] = myquery

	//statement to insert the hash of a succesfully loaded file
	query_insertalreadyloaded := "insert into alreadyloaded(hash) values (?)"
	stmt_insertalreadyloaded, err := tx.Prepare(query_insertalreadyloaded)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_insertalreadyloaded"
	myquery.sqlcode = query_insertalreadyloaded
	myquery.stmt = stmt_insertalreadyloaded
	myquerydb["stmt_insertalreadyloaded"] = myquery

	//debug statement to clean up de alreadyloaded table
	query_truncatealreadyloaded := "DELETE FROM alreadyloaded"
	stmt_truncatealreadyloaded, err := tx.Prepare(query_truncatealreadyloaded)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_truncatealreadyloaded"
	myquery.sqlcode = query_truncatealreadyloaded
	myquery.stmt = stmt_truncatealreadyloaded
	myquerydb["stmt_truncatealreadyloaded"] = myquery

}