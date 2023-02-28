package main

func demobarchart(args args) {
	XValues := []string{"Januari", "Februari", "March", "April"}
	YValues := make(map[string][]int)
	YValues["This year"] = append(YValues["This year"], 5)
	YValues["This year"] = append(YValues["This year"], 15)
	YValues["This year"] = append(YValues["This year"], 7)
	YValues["This year"] = append(YValues["This year"], 9)
	YValues["Last year"] = append(YValues["Last year"], 19)
	YValues["Last year"] = append(YValues["Last year"], 25)
	YValues["Last year"] = append(YValues["Last year"], 17)
	YValues["Last year"] = append(YValues["Last year"], 14)
	createbarchart(XValues, YValues, "this is a demo", "yup, a demo", args, "demobarchart.html")
}
