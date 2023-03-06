package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type querydb struct {
	querynaam string
	sqlcode   string
	stmt      *sql.Stmt
}

var myquerydb = make(map[string]querydb)

func loadquerydb(tx *sql.Tx) {
	var myquery querydb

	querymap := make(map[string]string)
	querymap["stmt_insertuser"] = "insert into user(ip, useragent) values (?,?)"                                                             //insert a new user
	querymap["stmt_insertrequest"] = "insert into request(request) values (?)"                                                               //insert a new request
	querymap["stmt_countreferrer"] = "select count(*) from referrer where referrer = ?"                                                      //count wether or not a certain referrer already exists
	querymap["stmt_selectreferrerid"] = "select id from referrer where referrer = ?"                                                         //return the id of a unique referrer
	querymap["stmt_insertreferrer"] = "select id from referrer where referrer = ?"                                                           //insert a new referrer
	querymap["stmt_selectrequestid"] = "select id from request where request = ?"                                                            //return the id of a unique request
	querymap["stmt_countrequest"] = "select count(*) from request where request = ?"                                                         //count wether or not a certain request already exists
	querymap["stmt_countusers"] = "select count(*) from user where ip = ? and useragent = ?"                                                 //count wether or not a certain user already exists
	querymap["stmt_selectuserid"] = "select id from user where ip = ? and useragent = ?"                                                     //return the id of a unique user
	querymap["stmt_insertvisit"] = "insert into visit(referrer, request,  visit_timestamp, user, statuscode, httpsize) values (?,?,?,?,?,?)" //insert a new visit record into the database
	querymap["stmt_countalreadyloaded"] = "select count(*) from alreadyloaded where hash = ?"                                                //count wether or not a file was already succesfully parsed in the past
	querymap["stmt_insertalreadyloaded"] = "insert into alreadyloaded(hash) values (?)"                                                      //insert a new sucesfully loaded file's hash into the database
	querymap["stmt_truncatealreadyloaded"] = "DELETE FROM alreadyloaded"                                                                     //truncate the alreadyloaded table so the system doesn't know wether a file was already loaded in the past
	querymap["stmt_maxvisittimestamp"] = "select max(visit_timestamp) from visit"                                                            //select the latest succesfully added record's timestamp to skip older records when loading
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] = " SELECT"                                                                      //select
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   strftime('%Y', datetime(visit_timestamp, 'unixepoch')) as year,"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   strftime('%m', datetime(visit_timestamp, 'unixepoch')) as month,"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   strftime('%d', datetime(visit_timestamp, 'unixepoch')) as day,"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   strftime('%H', datetime(visit_timestamp, 'unixepoch')) as hour,"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   COUNT(*) as count"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += " FROM"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   visit"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += " WHERE"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   visit_timestamp > ?"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "  AND"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   statuscode > 199"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "  AND"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   statuscode < 400"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += " GROUP BY"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   year, month, day, hour"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += " ORDER BY"
	querymap["stmt_raw_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   year asc, month asc, day asc, hour asc;"

	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] = " SELECT" //select
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   strftime('%Y', datetime(visit_timestamp, 'unixepoch')) as year,"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   strftime('%m', datetime(visit_timestamp, 'unixepoch')) as month,"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   strftime('%d', datetime(visit_timestamp, 'unixepoch')) as day,"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   strftime('%H', datetime(visit_timestamp, 'unixepoch')) as hour,"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   COUNT(distinct(user)) as count"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += " FROM"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   visit"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += " WHERE"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   visit_timestamp > ?"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "  AND"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   statuscode > 199"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "  AND"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   statuscode < 400"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += " GROUP BY"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   year, month, day, hour"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += " ORDER BY"
	querymap["stmt_unique_2xx_3xx_hourly_maxnbofdaysdetailed"] += "   year asc, month asc, day asc, hour asc;"

	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] = " SELECT"                                                                      //select
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   strftime('%Y', datetime(visit_timestamp, 'unixepoch')) as year,"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   strftime('%m', datetime(visit_timestamp, 'unixepoch')) as month,"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   strftime('%d', datetime(visit_timestamp, 'unixepoch')) as day,"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   COUNT(*) as count"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += " FROM"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   visit"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += " WHERE"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   visit_timestamp > ?"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "  AND"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   statuscode > 199"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "  AND"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   statuscode < 400"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += " GROUP BY"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   year, month, day"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += " ORDER BY"
	querymap["stmt_raw_2xx_3xx_daily_maxnbofdaysdetailed"] += "   year asc, month asc, day asc;"


	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] = " SELECT" //select
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   strftime('%Y', datetime(visit_timestamp, 'unixepoch')) as year,"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   strftime('%m', datetime(visit_timestamp, 'unixepoch')) as month,"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   strftime('%d', datetime(visit_timestamp, 'unixepoch')) as day,"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   COUNT(distinct(user)) as count"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += " FROM"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   visit"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += " WHERE"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   visit_timestamp > ?"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "  AND"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   statuscode > 199"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "  AND"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   statuscode < 400"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += " GROUP BY"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   year, month, day"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += " ORDER BY"
	querymap["stmt_unique_2xx_3xx_dayly_maxnbofdaysdetailed"] += "   year asc, month asc, day asc;"
	
	

	for naam, sql := range querymap {
		stmt, err := tx.Prepare(sql)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
		myquery.querynaam = naam
		myquery.sqlcode = sql
		myquery.stmt = stmt
		myquerydb[naam] = myquery
	}

}
