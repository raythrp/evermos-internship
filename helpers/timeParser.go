package helpers

import "time"

func TimeParserToDate(inputDate string) (time.Time, error) {
	// Parse the date string to a time.Time object
	parsedDate, err := time.Parse("02/01/2006", inputDate)
	if err != nil {
		return time.Time{}, err
	}

	// Set the location to +07:00 timezone
	location := time.FixedZone("WIB", 7*60*60)
	adjustedDate := parsedDate.In(location)

	return adjustedDate, nil
}