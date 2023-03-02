package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func emptydir(path string, extension string) {
	d, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	t := time.Now()
	mylog = append(mylog, t.Format("2006-01-02 15:04:05") + " => starting the output directory cleanup")
	mylog = append(mylog, t.Format("2006-01-02 15:04:05") + " => Reading " + path)
	

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == extension {
				os.Remove(path + file.Name())
				t := time.Now()
				mylog = append(mylog, t.Format("2006-01-02 15:04:05") + " => Deleted " + file.Name())
			}
		}
	}
}
