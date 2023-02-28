package main

func demolinegraph(args args) {
	XValues := []string{"Apple", "Banana", "Peach ", "Lemon"}
	YValues := make(map[string][]int)
	YValues["This year"] = append(YValues["This year"], 5)
	YValues["This year"] = append(YValues["This year"], 15)
	YValues["This year"] = append(YValues["This year"], 7)
	YValues["This year"] = append(YValues["This year"], 9)
	YValues["Last year"] = append(YValues["Last year"], 19)
	YValues["Last year"] = append(YValues["Last year"], 25)
	YValues["Last year"] = append(YValues["Last year"], 17)
	YValues["Last year"] = append(YValues["Last year"], 14)
	createlinegraph(XValues, YValues, "this is a demo", "yup, a demo", args, "demolinegraph.html") 
}

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


func demotable(args args) {
	MyHeaders := map[string]string{
		"Title_1": "kolom 1",
		"Title_2": "kolom 2",
	}

	myTable := Table{
		Pagetitle:       "tadaa",
		Pagedescription: "blahblah",
		Headers:         MyHeaders,
		Data:            []map[string]string{},
	}

	MyData := map[string]string{
		"Value_1": "record 1.1",
		"Value_2": "record 1.2",
	}
	myTable.Data = append(myTable.Data, MyData)
	MyData = map[string]string{
		"Value_1": "record 2.1",
		"Value_2": "record 2.2",
	}
	myTable.Data = append(myTable.Data, MyData)

	createtable(args, "demotable.html", "this is a title", myTable)
}