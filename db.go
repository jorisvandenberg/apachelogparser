package main

import (
	"database/sql"
	"os"
	"fmt"
	_ "github.com/mattn/go-sqlite3"

)


func createdb(dbnaam string) *sql.DB {
	db, err := sql.Open("sqlite3", dbnaam)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	return db
}

func initialisedb(db *sql.DB) *sql.Tx {
	var querylist []string
	querylist = append(querylist, "CREATE TABLE IF NOT EXISTS `user` (`id`    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`ip`    TEXT NOT NULL,`useragent`     TEXT);")
	querylist = append(querylist, "CREATE TABLE IF NOT EXISTS `request` (`id`    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`request`       TEXT NOT NULL);")
	querylist = append(querylist, "CREATE TABLE IF NOT EXISTS `referrer` (`id`    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`referrer`      TEXT NOT NULL);")
	querylist = append(querylist, "CREATE TABLE IF NOT EXISTS `alreadyloaded` (`id`    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`hash`      TEXT NOT NULL);")
	querylist = append(querylist, "CREATE TABLE IF NOT EXISTS `visit` ( `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,`referrer` INTEGER NOT NULL, `request` INTEGER NOT NULL, `visit_timestamp` INTEGER NOT NULL, `user`  INTEGER NOT NULL, `statuscode` INTEGER, `httpsize` INTEGER, FOREIGN KEY(`request`) REFERENCES `request`(`id`),  FOREIGN KEY(`referrer`) REFERENCES `referrer`(`id`),  FOREIGN KEY(`user`) REFERENCES `user`(`id`) 	);")
	querylist = append(querylist, "CREATE INDEX IF NOT EXISTS user_ip_agent on user(ip,useragent);")
	querylist = append(querylist, "CREATE INDEX IF NOT EXISTS request_request on request(request);")
	querylist = append(querylist, "CREATE INDEX IF NOT EXISTS referrer_referrer on referrer(referrer);")
	for _, query := range querylist {
		_, err := db.Exec(query)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			os.Exit(1)
		}
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
	return tx
}