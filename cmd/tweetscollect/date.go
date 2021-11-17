package main

import "time"

// parseDate parses date and return time.Time. e.g. "Fri Nov 12 20:38:53 +0000 2021".
func parseDate(d string) (string, error) {
	t, err := time.Parse(time.RubyDate, d)
	if err != nil {
		return "", err
	}
	return t.Format(time.RFC3339), nil
}
