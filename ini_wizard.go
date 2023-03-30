package main

import (
	"bufio"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
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
	fmt.Println("hidden feature: if you rename the file template_defaultconf.txt to defaultconf.txt and edit the values therein, this wizard will automatically take those values if you opt for the minimal configuration. You'll be able to create a customised configscript just by pressing enter when running this wizard!!!")
	noskipdb := []string{"dbpath", "mydomain", "logfilepath", "logfileregex", "outputpath"}

	var conf_dbpath, conf_mydomain, conf_logfilepath, conf_logfileregex, conf_outputpath string

	// Attempt to open the config file
	file, err := os.Open("defaultconf.txt")
	if err != nil {
		// If the file does not exist, assign default values to the variables
		conf_dbpath = ""
		conf_mydomain = ""
		conf_logfilepath = ""
		conf_logfileregex = ""
		conf_outputpath = ""
	} else {
		// If the file exists, read its contents and assign them to the variables
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, "=")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			switch key {
			case "dbpath":
				conf_dbpath = value
			case "mydomain":
				conf_mydomain = value
			case "logfilepath":
				conf_logfilepath = value
			case "logfileregex":
				conf_logfileregex = value
			case "outputpath":
				conf_outputpath = value
			}
		}
	}

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
				skipthisone := false
				if key.Name() == "dbpath" && conf_dbpath != "" {
					skipthisone = true
					newValue = conf_dbpath

				}
				if key.Name() == "mydomain" && conf_mydomain != "" {
					skipthisone = true
					newValue = conf_mydomain

				}
				if key.Name() == "logfilepath" && conf_logfilepath != "" {
					skipthisone = true
					newValue = conf_logfilepath

				}
				if key.Name() == "logfileregex" && conf_logfileregex != "" {
					skipthisone = true
					newValue = conf_logfileregex

				}
				if key.Name() == "outputpath" && conf_outputpath != "" {
					skipthisone = true
					newValue = conf_outputpath

				}
				if !skipthisone {
					fmt.Printf("%s.%s [%s]: ", section.Name(), key.Name(), key.Value())
					newValue, err = readStringFromUser()
					if err != nil {
						fmt.Printf("Error reading input: %v\n", err)
						os.Exit(1)
					}
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
