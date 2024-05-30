package tools

import (
	"fmt"
	"time"
)


// Creates a json time based from now plus hours and minutes.
func GetRelativeJsonTime(hours uint64, minutes uint64) string {
	timeIn := time.Now().
		Add(time.Hour * time.Duration(hours) + time.Minute * time.Duration(minutes))
	return timeIn.Format(time.RFC3339)
}


// Creates a json time based fixed point in time.
func GetFixedJsonTime(on string, at string) (string, error) {
	timestring := fmt.Sprintf("%v %v:00", on, at)
	timeIn, err := time.Parse("2006-01-02 15:04:05", timestring)
	if err != nil {
		return "", err
	}

	// Weird reset of timezone info without reinterpreting the given time.
	timeIn = time.Date(timeIn.Year(), timeIn.Month(), timeIn.Day(), timeIn.Hour(), timeIn.Minute(), timeIn.Second(), timeIn.Nanosecond(), time.Local)

	return timeIn.Format(time.RFC3339), nil
}
