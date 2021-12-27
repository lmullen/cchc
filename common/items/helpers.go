package items

import (
	"database/sql"
	"strconv"
)

// year takes an ISO8601 date string and figures out the year
func year(date string) sql.NullInt32 {
	var year sql.NullInt32
	null := sql.NullInt32{}

	if len(date) <= 4 {
		return null
	}

	s := date[0:4]
	y, err := strconv.Atoi(s)
	if err != nil {
		return null
	}

	year.Scan(y)
	return year
}

// stripXML is a package level function to strip XML/HTML from text
// var stripXML *bluemonday.Policy = bluemonday.StrictPolicy()
