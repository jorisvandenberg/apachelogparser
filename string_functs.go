package main

import (
	"strings"
	"strconv"
)

func splice_number_of_days_detailed_in(original string, number_of_days_detailed int) string {
	returnstring := strings.Replace(original, "|number_of_days_detailed|", strconv.Itoa(number_of_days_detailed), -1)
	return returnstring
}