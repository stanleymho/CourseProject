package main

import "time"

// parseDate parses date and return time.Time. e.g. "Fri Nov 12 20:38:53 +0000 2021".
func parseDate(d string) (time.Time, error) {
	return time.Parse(time.RubyDate, d)
}
