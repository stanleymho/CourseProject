package main

import (
	"regexp"
	"strings"
)

type SentimentType string

const (
	Sentiment_Neutral  = SentimentType("neutral")
	Sentiment_Positive = SentimentType("positive")
	Sentiment_Negative = SentimentType("negative")
	Sentiment_Mixed    = SentimentType("mixed")
)

// TweetsData represents the Twitter's dataset.
type TweetsData struct {
	Tweets []Tweet `json:"data"`
}

// Tweet represents a tweet.
type Tweet struct {
	Date          string `json:"date"`
	Text          string `json:"text"`
	Lang          string `json:"lang"`
	RetweetCount  int    `json:"retweet"`
	FavoriteCount int    `json:"favorite"`
}

func displayText(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func isRetweeted(tweetText string) bool {
	return strings.HasPrefix(tweetText, "RT ")
}

// trimTweet removes URLs, @<user> and other unused stuffs from tweet.
func trimTweet(tweetText string) string {
	result := tweetText

	if isRetweeted(result) {
		// For retweet, remove "RT @<user> :" prefix.
		if i := strings.Index(result, ":"); i > 0 {
			result = result[i+1:]
		}
		// Remove "…" suffix.
		if strings.HasSuffix(result, "…") {
			result = result[:len(result)-1]
			if i := strings.LastIndex(result, " "); i > 0 {
				result = result[:i]
			}
		}
		result = strings.TrimSpace(result)
	}

	// Remove "http://...", "https://...", or @<user> from the text.
	tokens := make([]string, 0, 1)
	for _, token := range regexp.MustCompile("[ \\n\\t\\r]+").Split(result, -1) {
		if !strings.HasPrefix(token, "http://") &&
			!strings.HasPrefix(token, "https://") &&
			!strings.HasPrefix(token, "@") &&
			strings.TrimSpace(token) != "" {
			tokens = append(tokens, token)
		}
	}
	result = strings.Join(tokens, " ")
	return strings.TrimSpace(result)
}

func hashKey(tweetText string) string {
	result := trimTweet(tweetText)
	if len(result) > 100 {
		result = result[:100]
	}
	if result = strings.TrimSpace(result); result != "" {
		result += "…"
	}
	return result
}
