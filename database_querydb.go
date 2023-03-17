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
	querymap["stmt_insertreferrer"] = "insert into referrer (referrer) values (?)"                                                           //insert a new referrer
	querymap["stmt_selectrequestid"] = "select id from request where request = ?"                                                            //return the id of a unique request
	querymap["stmt_countrequest"] = "select count(*) from request where request = ?"                                                         //count wether or not a certain request already exists
	querymap["stmt_countusers"] = "select count(*) from user where ip = ? and useragent = ?"                                                 //count wether or not a certain user already exists
	querymap["stmt_selectuserid"] = "select id from user where ip = ? and useragent = ?"                                                     //return the id of a unique user
	querymap["stmt_insertvisit"] = "insert into visit(referrer, request,  visit_timestamp, user, statuscode, httpsize) values (?,?,?,?,?,?)" //insert a new visit record into the database
	querymap["stmt_countalreadyloaded"] = "select count(*) from alreadyloaded where hash = ?"                                                //count wether or not a file was already succesfully parsed in the past
	querymap["stmt_insertalreadyloaded"] = "insert into alreadyloaded(hash) values (?)"                                                      //insert a new sucesfully loaded file's hash into the database
	querymap["stmt_truncatealreadyloaded"] = "DELETE FROM alreadyloaded"                                                                     //truncate the alreadyloaded table so the system doesn't know wether a file was already loaded in the past
	querymap["stmt_maxvisittimestamp"] = "select max(visit_timestamp) from visit"                                                            //select the latest succesfully added record's timestamp to skip older records when loading
	querymap["stmt_raw_PerHour_hits"] = " SELECT"                                                                      //select
	querymap["stmt_raw_PerHour_hits"] += "   strftime('%Y', datetime(visit_timestamp, 'unixepoch')) as year,"
	querymap["stmt_raw_PerHour_hits"] += "   strftime('%m', datetime(visit_timestamp, 'unixepoch')) as month,"
	querymap["stmt_raw_PerHour_hits"] += "   strftime('%d', datetime(visit_timestamp, 'unixepoch')) as day,"
	querymap["stmt_raw_PerHour_hits"] += "   strftime('%H', datetime(visit_timestamp, 'unixepoch')) as hour,"
	querymap["stmt_raw_PerHour_hits"] += "   COUNT(*) as count"
	querymap["stmt_raw_PerHour_hits"] += " FROM"
	querymap["stmt_raw_PerHour_hits"] += "   visit"
	querymap["stmt_raw_PerHour_hits"] += " WHERE"
	querymap["stmt_raw_PerHour_hits"] += "   visit_timestamp > ?"
	querymap["stmt_raw_PerHour_hits"] += "  AND"
	querymap["stmt_raw_PerHour_hits"] += "   statuscode > 199"
	querymap["stmt_raw_PerHour_hits"] += "  AND"
	querymap["stmt_raw_PerHour_hits"] += "   statuscode < 400"
	querymap["stmt_raw_PerHour_hits"] += " GROUP BY"
	querymap["stmt_raw_PerHour_hits"] += "   year, month, day, hour"
	querymap["stmt_raw_PerHour_hits"] += " ORDER BY"
	querymap["stmt_raw_PerHour_hits"] += "   year asc, month asc, day asc, hour asc;"

	querymap["stmt_unique_PerHour_hits"] = " SELECT" //select
	querymap["stmt_unique_PerHour_hits"] += "   strftime('%Y', datetime(visit_timestamp, 'unixepoch')) as year,"
	querymap["stmt_unique_PerHour_hits"] += "   strftime('%m', datetime(visit_timestamp, 'unixepoch')) as month,"
	querymap["stmt_unique_PerHour_hits"] += "   strftime('%d', datetime(visit_timestamp, 'unixepoch')) as day,"
	querymap["stmt_unique_PerHour_hits"] += "   strftime('%H', datetime(visit_timestamp, 'unixepoch')) as hour,"
	querymap["stmt_unique_PerHour_hits"] += "   COUNT(distinct(user)) as count"
	querymap["stmt_unique_PerHour_hits"] += " FROM"
	querymap["stmt_unique_PerHour_hits"] += "   visit"
	querymap["stmt_unique_PerHour_hits"] += " WHERE"
	querymap["stmt_unique_PerHour_hits"] += "   visit_timestamp > ?"
	querymap["stmt_unique_PerHour_hits"] += "  AND"
	querymap["stmt_unique_PerHour_hits"] += "   statuscode > 199"
	querymap["stmt_unique_PerHour_hits"] += "  AND"
	querymap["stmt_unique_PerHour_hits"] += "   statuscode < 400"
	querymap["stmt_unique_PerHour_hits"] += " GROUP BY"
	querymap["stmt_unique_PerHour_hits"] += "   year, month, day, hour"
	querymap["stmt_unique_PerHour_hits"] += " ORDER BY"
	querymap["stmt_unique_PerHour_hits"] += "   year asc, month asc, day asc, hour asc;"

	querymap["stmt_raw_PerDay_hits"] = " SELECT" //select
	querymap["stmt_raw_PerDay_hits"] += "   strftime('%Y', datetime(visit_timestamp, 'unixepoch')) as year,"
	querymap["stmt_raw_PerDay_hits"] += "   strftime('%m', datetime(visit_timestamp, 'unixepoch')) as month,"
	querymap["stmt_raw_PerDay_hits"] += "   strftime('%d', datetime(visit_timestamp, 'unixepoch')) as day,"
	querymap["stmt_raw_PerDay_hits"] += "   COUNT(*) as count"
	querymap["stmt_raw_PerDay_hits"] += " FROM"
	querymap["stmt_raw_PerDay_hits"] += "   visit"
	querymap["stmt_raw_PerDay_hits"] += " WHERE"
	querymap["stmt_raw_PerDay_hits"] += "   visit_timestamp > ?"
	querymap["stmt_raw_PerDay_hits"] += "  AND"
	querymap["stmt_raw_PerDay_hits"] += "   statuscode > 199"
	querymap["stmt_raw_PerDay_hits"] += "  AND"
	querymap["stmt_raw_PerDay_hits"] += "   statuscode < 400"
	querymap["stmt_raw_PerDay_hits"] += " GROUP BY"
	querymap["stmt_raw_PerDay_hits"] += "   year, month, day"
	querymap["stmt_raw_PerDay_hits"] += " ORDER BY"
	querymap["stmt_raw_PerDay_hits"] += "   year asc, month asc, day asc;"

	querymap["stmt_unique_PerDay_hits"] = " SELECT" //select
	querymap["stmt_unique_PerDay_hits"] += "   strftime('%Y', datetime(visit_timestamp, 'unixepoch')) as year,"
	querymap["stmt_unique_PerDay_hits"] += "   strftime('%m', datetime(visit_timestamp, 'unixepoch')) as month,"
	querymap["stmt_unique_PerDay_hits"] += "   strftime('%d', datetime(visit_timestamp, 'unixepoch')) as day,"
	querymap["stmt_unique_PerDay_hits"] += "   COUNT(distinct(user)) as count"
	querymap["stmt_unique_PerDay_hits"] += " FROM"
	querymap["stmt_unique_PerDay_hits"] += "   visit"
	querymap["stmt_unique_PerDay_hits"] += " WHERE"
	querymap["stmt_unique_PerDay_hits"] += "   visit_timestamp > ?"
	querymap["stmt_unique_PerDay_hits"] += "  AND"
	querymap["stmt_unique_PerDay_hits"] += "   statuscode > 199"
	querymap["stmt_unique_PerDay_hits"] += "  AND"
	querymap["stmt_unique_PerDay_hits"] += "   statuscode < 400"
	querymap["stmt_unique_PerDay_hits"] += " GROUP BY"
	querymap["stmt_unique_PerDay_hits"] += "   year, month, day"
	querymap["stmt_unique_PerDay_hits"] += " ORDER BY"
	querymap["stmt_unique_PerDay_hits"] += "   year asc, month asc, day asc;"

	querymap["stmt_raw_PerHour_ReferringUrls"] = " SELECT"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   CASE"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "     WHEN instr(r.referrer, '?') > 0 THEN REPLACE(RTRIM(substr(r.referrer, 1, instr(r.referrer, '?') - 1), '/'), '//', '/')"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "     ELSE REPLACE(RTRIM(r.referrer, '/'), '//', '/')"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   END AS subreferrer,"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   count(*) as aantal"
	querymap["stmt_raw_PerHour_ReferringUrls"] += " FROM"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   referrer r, visit v"
	querymap["stmt_raw_PerHour_ReferringUrls"] += " WHERE"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   r.id = v.referrer"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   and v.visit_timestamp > ?"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   and v.statuscode > 199"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   and v.statuscode < 400"
	querymap["stmt_raw_PerHour_ReferringUrls"] += " GROUP BY"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   subreferrer"
	querymap["stmt_raw_PerHour_ReferringUrls"] += " ORDER BY"
	querymap["stmt_raw_PerHour_ReferringUrls"] += "   aantal desc"
	querymap["stmt_raw_PerHour_ReferringUrls"] += " LIMIT ?"

	querymap["stmt_unique_PerHour_ReferringUrls"] = " SELECT"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   CASE"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "     WHEN instr(r.referrer, '??') > 0 THEN REPLACE(RTRIM(substr(r.referrer, 1, instr(r.referrer, '??') - 1), '/'), '//', '/')"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "     ELSE REPLACE(RTRIM(r.referrer, '/'), '//', '/')"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   END AS subreferrer,"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   count(distinct(v.user)) as aantal"
	querymap["stmt_unique_PerHour_ReferringUrls"] += " FROM"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   referrer r, visit v"
	querymap["stmt_unique_PerHour_ReferringUrls"] += " WHERE"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   r.id = v.referrer"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   and v.visit_timestamp > ?"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   and v.statuscode > 199"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   and v.statuscode < 400"
	querymap["stmt_unique_PerHour_ReferringUrls"] += " GROUP BY"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   subreferrer"
	querymap["stmt_unique_PerHour_ReferringUrls"] += " ORDER BY"
	querymap["stmt_unique_PerHour_ReferringUrls"] += "   aantal desc"
	querymap["stmt_unique_PerHour_ReferringUrls"] += " LIMIT ?"

	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] = " SELECT"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   CASE"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "     WHEN instr(r.referrer, '??') > 0 THEN REPLACE(RTRIM(substr(r.referrer, 1, instr(r.referrer, '??') - 1), '/'), '//', '/')"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "     ELSE REPLACE(RTRIM(r.referrer, '/'), '//', '/')"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   END AS subreferrer,"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   count(distinct(v.user)) as aantal"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += " FROM"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   referrer r, visit v"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += " WHERE"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   r.id = v.referrer"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   and v.visit_timestamp > ?"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   and v.statuscode > 199"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   and v.statuscode < 400"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   and r.referrer != \"-\""
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   and r.referrer not like '%' || ? || '%'"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   and v.statuscode < 400"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += " GROUP BY"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   subreferrer"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += " ORDER BY"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += "   aantal desc"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelf"] += " LIMIT ?"

	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] = " SELECT"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "     SUBSTR(SUBSTR(r.referrer, INSTR(r.referrer, '//') + 2), 0, INSTR(SUBSTR(r.referrer, INSTR(r.referrer, '//') + 2), '/')) AS subreferrer,"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   count(distinct(v.user)) as aantal"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += " FROM"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   referrer r, visit v"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += " WHERE"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   r.id = v.referrer"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   and v.visit_timestamp > ?"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   and v.statuscode > 199"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   and v.statuscode < 400"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   and r.referrer != \"-\""
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   and r.referrer not like '%' || ? || '%'"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   and v.statuscode < 400"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += " GROUP BY"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   subreferrer"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += " ORDER BY"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += "   aantal desc"
	querymap["stmt_unique_PerHour_RefferingUrlsNoEmptyOrSelfOnlyTld"] += " LIMIT ?"

	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] = "   SELECT"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    SUBSTR(SUBSTR(referrer.referrer, INSTR(referrer.referrer, '//') + 2), 0, INSTR(SUBSTR(referrer.referrer, INSTR(referrer.referrer, '//') + 2), '/')) AS searchengine,"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    count(*) as aantal"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "  FROM"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer, visit"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "  WHERE"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    visit.referrer = referrer.id and"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "	visit.statuscode > 199 and"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "	visit.statuscode < 400 and"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    visit.visit_timestamp > ? and ("
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%google%%\" or"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%bing%%\" or"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%yahoo%%\" or"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%yandex%%\" or"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%baidu%%\" or"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%duckduckgo%%\" or"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%ecosia%%\" or"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%naver%%\" or"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    referrer.referrer like \"%%aol%%\")"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "  GROUP BY"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    searchengine"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "  ORDER BY"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    aantal DESC"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "  LIMIT"
	querymap["stmt_stat_raw_XDaysTotal_HitsFromSearchEngines"] += "    ?"

	querymap["stmt_count_unique_nbhits_per_searchengine"] = "   SELECT"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    SUBSTR(SUBSTR(referrer.referrer, INSTR(referrer.referrer, '//') + 2), 0, INSTR(SUBSTR(referrer.referrer, INSTR(referrer.referrer, '//') + 2), '/')) AS searchengine,"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    count(distinct(visit.user)) as aantal"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "  FROM"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer, visit"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "  WHERE"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    visit.referrer = referrer.id and"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "	visit.statuscode > 199 and"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "	visit.statuscode < 400 and"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    visit.visit_timestamp > ? and ("
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%google%%\" or"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%bing%%\" or"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%yahoo%%\" or"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%yandex%%\" or"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%baidu%%\" or"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%duckduckgo%%\" or"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%ecosia%%\" or"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%naver%%\" or"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    referrer.referrer like \"%%aol%%\")"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "  GROUP BY"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    searchengine"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "  ORDER BY"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    aantal DESC"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "  LIMIT"
	querymap["stmt_count_unique_nbhits_per_searchengine"] += "    ?"

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
