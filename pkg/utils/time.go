package utils

import "time"

func AfterNow(date, now time.Time) bool {
	return date.After(now)
}

const DateFormat = "20060102"
