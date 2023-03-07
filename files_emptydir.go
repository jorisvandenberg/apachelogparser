package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	logger("starting the output directory cleanup")
	logger("Reading " + path)

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == extension {
				os.Remove(path + file.Name())
				logger("Deleted " + file.Name())
			}
		}
	}
}
