package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ini_wizard(inpath string, outpath string) {
	noskipdb := []string{"dbpath", "mydomain", "logfilepath", "logfileregex", "outputpath"}
	// Open the INI file
	cfg, err := ini.Load(inpath)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Do you want to do a minimal config (all stats enabled, default output values. Only works with apache common log format). y/n [y]")
	myChoice, err := readStringFromUser()
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}
	skipstd := false
	if myChoice == "n" || myChoice == "N" {
		skipstd = false
	} else {
		skipstd = true
	}
	// Prompt the user to change the settings
	fmt.Println("Enter new values or leave blank to keep current settings.")
	for _, section := range cfg.Sections() {
		if section.Name() == ini.DEFAULT_SECTION {
			continue
		}

		for _, key := range section.Keys() {
			newValue := ""
			if (contains(noskipdb, key.Name()) && skipstd) || !skipstd {
				fmt.Printf("%s.%s [%s]: ", section.Name(), key.Name(), key.Value())
				newValue, err = readStringFromUser()
				if err != nil {
					fmt.Printf("Error reading input: %v\n", err)
					os.Exit(1)
				}

			}
			if newValue != "" {
				key.SetValue(newValue)
			}

		}
	}

	// Write the updated settings to a new INI file
	err = cfg.SaveTo(outpath)
	if err != nil {
		fmt.Printf("Failed to write config file: %v\n", err)
		os.Exit(1)
	}
}

func readStringFromUser() (string, error) {
	var buf strings.Builder
	for {
		var char rune
		_, err := fmt.Scanf("%c", &char)
		if err != nil {
			return "", err
		}

		if char == '\n' {
			break
		}

		buf.WriteRune(char)
	}

	return buf.String(), nil
}
