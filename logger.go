package main

import (
	"time"
)

func logger(logtext string) {
	t := time.Now()
	mylog = append(mylog, t.Format("2006-01-02 15:04:05")+" => "+logtext)
}
