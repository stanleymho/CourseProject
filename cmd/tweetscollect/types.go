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

// TwitterSearchResponse represents the Twitter's search response.
type TwitterSearchResponseV1 struct {
	Tweets         []TweetV1        `json:"statuses"`
	SearchMetadata SearchMetadataV1 `json:"search_metadata"`
}

// Tweet represents a tweet.
type TweetV1 struct {
	ID            string  `json:"id_str"`
	Text          string  `json:"full_text"`
	Lang          string  `json:"lang"`
	RetweetCount  int     `json:"retweet_count"`
	FavoriteCount int     `json:"favorite_count"`
	User          *UserV1 `json:"user"`
	CreatedAt     string  `json:"created_at"`
}

func (t *TweetV1) NormalizedText() string {
	s := t.Text
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// UserV1 represents a Twitter user.
type UserV1 struct {
	ID             string `json:"id_str"`
	Name           string `json:"name"`
	ScreenName     string `json:"screen_name"`
	FollowersCount int    `json:"followers_count"`
	FriendsCount   int    `json:"friends_count"`
	ListedCount    int    `json:"listed_count"`
	CreatedAt      string `json:"created_at"`
}

// SearchMetadataV1 represents the metadata of the search.
type SearchMetadataV1 struct {
	NextURLPath string `json:"next_results"`
	Query       string `json:"query"`
}
