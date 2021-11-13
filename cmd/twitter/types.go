package main

import (
	"strings"
)

// For more details on the Twitter Search v1.1 API, please see
// https://developer.twitter.com/en/docs/twitter-api/v1/tweets/search/api-reference/get-search-tweets
//

const (
	twitterSearchURL = "https://api.twitter.com/1.1/search/tweets.json"
)

// TwitterSearchResponse represents the Twitter search response.
type TwitterSearchResponse struct {
	Tweets         []Tweet        `json:"statuses"`
	SearchMetadata SearchMetadata `json:"search_metadata"`
}

// Tweet represents a tweet.
type Tweet struct {
	ID            string `json:"id_str"`
	FullText      string `json:"full_text"`
	Truncated     bool   `json:"truncated"`
	RetweetCount  int    `json:"retweet_count"`
	FavoriteCount int    `json:"favorite_count"`
	Lang          string `json:"lang"`
	User          User   `json:"user"`
	CreatedAt     string `json:"created_at"`
}

func (t *Tweet) NormalizedFullText() string {
	s := t.FullText
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// User represents a Twitter user.
type User struct {
	ID             string `json:"id_str"`
	Name           string `json:"name"`
	ScreenName     string `json:"screen_name"`
	FollowersCount int    `json:"followers_count"`
	FriendsCount   int    `json:"friends_count"`
	ListedCount    int    `json:"listed_count"`
	CreatedAt      string `json:"created_at"`
}

// SearchMetadata represents the metadata of the search.
type SearchMetadata struct {
	NextURLPath string `json:"next_results"`
}
