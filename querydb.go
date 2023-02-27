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

	//statement to insert a new user
	query_insertuser := "insert into user(ip, useragent) values (?,?)"
	stmt_insertuser, err := tx.Prepare(query_insertuser)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_insertuser"
	myquery.sqlcode = query_insertuser
	myquery.stmt = stmt_insertuser
	myquerydb["stmt_insertuser"] = myquery

	//query to insert a new request
	query_insertrequest := "insert into request(request) values (?)"
	stmt_insertrequest, err := tx.Prepare(query_insertrequest)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_insertrequest"
	myquery.sqlcode = query_insertrequest
	myquery.stmt = stmt_insertrequest
	myquerydb["stmt_insertrequest"] = myquery

	//statement to count if a referrer already exists
	query_countreferrer := "select count(*) from referrer where referrer = ?"
	stmt_countreferrer, err := tx.Prepare(query_countreferrer)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_countreferrer"
	myquery.sqlcode = query_countreferrer
	myquery.stmt = stmt_countreferrer
	myquerydb["stmt_countreferrer"] = myquery

	//statement to get the id of an existing referrer
	query_selectreferrerid := "select id from referrer where referrer = ?"
	stmt_selectreferrerid, err := tx.Prepare(query_selectreferrerid)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_selectreferrerid"
	myquery.sqlcode = query_selectreferrerid
	myquery.stmt = stmt_selectreferrerid
	myquerydb["stmt_selectreferrerid"] = myquery

	//statement to insert a new referrer
	query_insertreferrer := "insert into referrer(referrer) values (?)"
	stmt_insertreferrer, err := tx.Prepare(query_insertreferrer)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_insertreferrer"
	myquery.sqlcode = query_insertreferrer
	myquery.stmt = stmt_insertreferrer
	myquerydb["stmt_insertreferrer"] = myquery

	//statement to get the id of a unique statement
	query_selectrequestid := "select id from request where request = ?"
	stmt_selectrequestid, err := tx.Prepare(query_selectrequestid)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_selectrequestid"
	myquery.sqlcode = query_selectrequestid
	myquery.stmt = stmt_selectrequestid
	myquerydb["stmt_selectrequestid"] = myquery

	//statement to count if a request already exists
	query_countrequest := "select count(*) from request where request = ?"
	stmt_countrequest, err := tx.Prepare(query_countrequest)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_countrequest"
	myquery.sqlcode = query_countrequest
	myquery.stmt = stmt_countrequest
	myquerydb["stmt_countrequest"] = myquery

	//statement to count unique users 
	query_countusers := "select count(*) from user where ip = ? and useragent = ?"
	stmt_countusers, err := tx.Prepare(query_countusers)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_countusers"
	myquery.sqlcode = query_countusers
	myquery.stmt = stmt_countusers
	myquerydb["stmt_countusers"] = myquery

	//statement to get the id of an existing user record
	query_selectuserid := "select id from user where ip = ? and useragent = ?"
	stmt_selectuserid, err := tx.Prepare(query_selectuserid)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_selectuserid"
	myquery.sqlcode = query_selectuserid
	myquery.stmt = stmt_selectuserid
	myquerydb["stmt_selectuserid"] = myquery

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

	//statement for fetching the last known timestamp from the db
	query_maxvisittimestamp := "select max(visit_timestamp) from visit"
	stmt_maxvisittimestamp, err := tx.Prepare(query_maxvisittimestamp)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	myquery.querynaam = "stmt_maxvisittimestamp"
	myquery.sqlcode = query_maxvisittimestamp
	myquery.stmt = stmt_maxvisittimestamp
	myquerydb["stmt_maxvisittimestamp"] = myquery

}