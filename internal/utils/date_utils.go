package utils

import "time"

var AllMonthsInt = []time.Month{
	time.January, time.February, time.March,
	time.April, time.May, time.June,
	time.July, time.August, time.September,
	time.October, time.November, time.December,
}

var AllMontsStr = []string{
	"Jan", "Feb", "Mar",
	"Apr", "May", "Jun",
	"Jul", "Aug", "Sep",
	"Oct", "Nov", "Dec",
}

var AllMonthsMapping = map[string]time.Month{
	"Jan": time.January, "Feb": time.February, "Mar": time.March,
	"Apr": time.April, "May": time.May, "Jun": time.June,
	"Jul": time.July, "Aug": time.August, "Sep": time.September,
	"Oct": time.October, "Nov": time.November, "Dec": time.December,
}
