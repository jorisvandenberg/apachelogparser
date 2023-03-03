package main

//credits: https://stackoverflow.com/questions/37869793/how-do-i-zip-a-directory-containing-sub-directories-or-files-in-golang

import (
	"os"
	"fmt"
	"archive/zip"
	"io/ioutil"
)

func ZipWriter(inputfolder string, outputfile string) {
    outFile, err := os.Create(outputfile)
    if err != nil {
        fmt.Println(err)
    }
    defer outFile.Close()

    // Create a new zip archive.
    w := zip.NewWriter(outFile)

    // Add some files to the archive.
    addFiles(w, inputfolder, "")

    if err != nil {
        fmt.Println(err)
    }

    // Make sure to check the error on Close.
    err = w.Close()
    if err != nil {
        fmt.Println(err)
    }
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
    // Open the Directory
    files, err := ioutil.ReadDir(basePath)
    if err != nil {
        fmt.Println(err)
    }

    for _, file := range files {
        fmt.Println(basePath + file.Name())
        if !file.IsDir() {
            dat, err := ioutil.ReadFile(basePath + file.Name())
            if err != nil {
                fmt.Println(err)
            }

            // Add some files to the archive.
            f, err := w.Create(baseInZip + file.Name())
            if err != nil {
                fmt.Println(err)
            }
            _, err = f.Write(dat)
            if err != nil {
                fmt.Println(err)
            }
        } else if file.IsDir() {

            // Recurse
            newBase := basePath + file.Name() + "/"
            fmt.Println("Recursing and Adding SubDir: " + file.Name())
            fmt.Println("Recursing and Adding SubDir: " + newBase)

            addFiles(w, newBase, baseInZip  + file.Name() + "/")
        }
    }
}