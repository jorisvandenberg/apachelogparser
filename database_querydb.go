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

	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] = " SELECT"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   CASE"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "     WHEN instr(r.referrer, '?') > 0 THEN REPLACE(RTRIM(substr(r.referrer, 1, instr(r.referrer, '?') - 1), '/'), '//', '/')"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "     ELSE REPLACE(RTRIM(r.referrer, '/'), '//', '/')"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   END AS subreferrer,"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   count(*) as aantal"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += " FROM"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   referrer r, visit v"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += " WHERE"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   r.id = v.referrer"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   and v.visit_timestamp > ?"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   and v.statuscode > 199"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   and v.statuscode < 400"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += " GROUP BY"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   subreferrer"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += " ORDER BY"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += "   aantal desc"
	querymap["stmt_noaggregation_nbdaysdetailed_refferers_noparams_2xx_3xx"] += " LIMIT ?"

	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] = " SELECT"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   CASE"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "     WHEN instr(r.referrer, '??') > 0 THEN REPLACE(RTRIM(substr(r.referrer, 1, instr(r.referrer, '??') - 1), '/'), '//', '/')"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "     ELSE REPLACE(RTRIM(r.referrer, '/'), '//', '/')"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   END AS subreferrer,"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   count(distinct(v.user)) as aantal"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += " FROM"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   referrer r, visit v"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += " WHERE"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   r.id = v.referrer"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   and v.visit_timestamp > ?"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   and v.statuscode > 199"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   and v.statuscode < 400"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += " GROUP BY"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   subreferrer"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += " ORDER BY"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += "   aantal desc"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_2xx_3xx"] += " LIMIT ?"

	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] = " SELECT"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   CASE"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "     WHEN instr(r.referrer, '??') > 0 THEN REPLACE(RTRIM(substr(r.referrer, 1, instr(r.referrer, '??') - 1), '/'), '//', '/')"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "     ELSE REPLACE(RTRIM(r.referrer, '/'), '//', '/')"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   END AS subreferrer,"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   count(distinct(v.user)) as aantal"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += " FROM"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   referrer r, visit v"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += " WHERE"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   r.id = v.referrer"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   and v.visit_timestamp > ?"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   and v.statuscode > 199"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   and v.statuscode < 400"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   and r.referrer != \"-\""
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   and r.referrer not like '%' || ? || '%'"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   and v.statuscode < 400"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += " GROUP BY"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   subreferrer"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += " ORDER BY"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += "   aantal desc"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_2xx_3xx"] += " LIMIT ?"

	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] = " SELECT"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "     SUBSTR(SUBSTR(r.referrer, INSTR(r.referrer, '//') + 2), 0, INSTR(SUBSTR(r.referrer, INSTR(r.referrer, '//') + 2), '/')) AS subreferrer,"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   count(distinct(v.user)) as aantal"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += " FROM"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   referrer r, visit v"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += " WHERE"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   r.id = v.referrer"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   and v.visit_timestamp > ?"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   and v.statuscode > 199"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   and v.statuscode < 400"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   and r.referrer != \"-\""
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   and r.referrer not like '%' || ? || '%'"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   and v.statuscode < 400"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += " GROUP BY"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   subreferrer"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += " ORDER BY"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += "   aantal desc"
	querymap["stmt_noaggregation_nbdaysdetailed_unique_refferers_noparams_noemptyorown_tld_2xx_3xx"] += " LIMIT ?"

	querymap["stmt_count_nbhits_per_searchengine"] = "   SELECT"
	querymap["stmt_count_nbhits_per_searchengine"] += "    SUBSTR(SUBSTR(referrer.referrer, INSTR(referrer.referrer, '//') + 2), 0, INSTR(SUBSTR(referrer.referrer, INSTR(referrer.referrer, '//') + 2), '/')) AS searchengine,"
	querymap["stmt_count_nbhits_per_searchengine"] += "    count(*) as aantal"
	querymap["stmt_count_nbhits_per_searchengine"] += "  FROM"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer, visit"
	querymap["stmt_count_nbhits_per_searchengine"] += "  WHERE"
	querymap["stmt_count_nbhits_per_searchengine"] += "    visit.referrer = referrer.id and"
	querymap["stmt_count_nbhits_per_searchengine"] += "	visit.statuscode > 199 and"
	querymap["stmt_count_nbhits_per_searchengine"] += "	visit.statuscode < 400 and"
	querymap["stmt_count_nbhits_per_searchengine"] += "    visit.visit_timestamp > ? and ("
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%google%%\" or"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%bing%%\" or"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%yahoo%%\" or"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%yandex%%\" or"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%baidu%%\" or"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%duckduckgo%%\" or"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%ecosia%%\" or"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%naver%%\" or"
	querymap["stmt_count_nbhits_per_searchengine"] += "    referrer.referrer like \"%%aol%%\")"
	querymap["stmt_count_nbhits_per_searchengine"] += "  GROUP BY"
	querymap["stmt_count_nbhits_per_searchengine"] += "    searchengine"
	querymap["stmt_count_nbhits_per_searchengine"] += "  ORDER BY"
	querymap["stmt_count_nbhits_per_searchengine"] += "    aantal DESC"
	querymap["stmt_count_nbhits_per_searchengine"] += "  LIMIT"
	querymap["stmt_count_nbhits_per_searchengine"] += "    ?"

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
