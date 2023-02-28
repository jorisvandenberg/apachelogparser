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

      fmt.Println("Reading "+ path)

      for _, file := range files {
          if file.Mode().IsRegular() {
              if filepath.Ext(file.Name()) == extension {
                os.Remove(path + file.Name())
                fmt.Println("Deleted ", file.Name())
              }
          }
      }
}