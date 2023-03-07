package main

func generatestats(args args) {
	logger("started the function to generate statistics")
	noaggregation_nbdaysdetailed_raw_2xx_3xx(args)
	noaggregation_nbdaysdetailed_unique_2xx_3xx(args)
	noaggregation_nbdaysdetailed_rawperday_2xx_3xx(args)
	noaggregation_nbdaysdetailed_uniqueperday_2xx_3xx(args)
	stats_noaggregation_nbdaysdetailed_refferers(args)
	logger("finished the function to generate statistics")
}
