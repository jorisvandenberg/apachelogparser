package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"crypto/sha256"
	"io"
	"encoding/hex"
)

func truncatealreadyloaded () {
	stmt_truncatealreadyloaded := myquerydb["stmt_truncatealreadyloaded"].stmt

	_, err := stmt_truncatealreadyloaded.Exec()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}

func InsertParsedFileHashIntoDb(filename string, filepath string) {

	filehandle, err := os.Open(filepath + filename)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer filehandle.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, filehandle); err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	filehash := hex.EncodeToString(hash.Sum(nil))
	stmt_insertalreadyloaded := myquerydb["stmt_insertalreadyloaded"].stmt

	_, err = stmt_insertalreadyloaded.Exec(filehash)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}

}

func getfiles(regex string, pathS string) []string {
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(regex, f.Name())
			if err == nil && r {
				filehandle, err := os.Open(pathS + f.Name())
				if err != nil {
					fmt.Printf("%s\n", err.Error())
				}
				defer filehandle.Close()
				
				hash := sha256.New()
				if _, err := io.Copy(hash, filehandle); err != nil {
					fmt.Printf("%s\n", err.Error())
				}
				filehash := hex.EncodeToString(hash.Sum(nil))
				stmt_countalreadyloaded := myquerydb["stmt_countalreadyloaded"].stmt
				var countalreadyloaded int
				stmt_countalreadyloaded.QueryRow(filehash).Scan(&countalreadyloaded)
				
				if countalreadyloaded == 0 {
					files = append(files, f.Name())
					fmt.Printf("%s added to the todo list\n", f.Name())
				} else {
					fmt.Printf("%s was already parsed in the past... skipping\n", f.Name())
				}

			}
		}
		return nil
	})
	return files
}

func parselogs(args args) {
	toparselist := getfiles(args.logfileregex, args.logfilepath)
	fmt.Printf("%+v\n", toparselist)
	for _, filename := range toparselist {
		InsertParsedFileHashIntoDb(filename, args.logfilepath)
	}
}