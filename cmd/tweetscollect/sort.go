package main

import "strings"

type Tweets []TweetV1

func (t Tweets) Len() int {
	return len(t)
}

func (t Tweets) Less(i, j int) bool {
	date1, _ := parseDate(t[i].CreatedAt)
	date2, _ := parseDate(t[j].CreatedAt)
	return strings.Compare(date1, date2) > 0
}

func (t Tweets) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
