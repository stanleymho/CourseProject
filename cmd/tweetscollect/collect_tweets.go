package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

// collectTweets collects the tweets for a topic.
func collectTweets(bearerToken, topic string) error {
	ctx := context.Background()

	// Constuct URL.
	var sb strings.Builder
	sb.WriteString("?lang=en")
	sb.WriteString("&result_type=mixed")
	sb.WriteString("&tweet_mode=extended")
	sb.WriteString("&include_entities=false")
	sb.WriteString("&count=100")
	sb.WriteString("&q=" + url.QueryEscape(topic))
	searchURLPath := sb.String()

	tweets := make([]TweetV1, 0, 100)
	for {
		url := fmt.Sprintf("%s%s", twitterSearchURL, searchURLPath)
		t, nextSearchURLPath, err := collectTweetsByURL(ctx, bearerToken, url)
		if err != nil {
			return err
		}

		for _, tweet := range t {
			createdAt, _ := parseDate(tweet.CreatedAt)
			text := tweet.NormalizedText()
			fmt.Printf("{ date=\"%v\", text=\"%s\", lang=\"%s\", favorite=%d, retweet=%d }\n",
				createdAt, text, tweet.Lang, tweet.FavoriteCount, tweet.RetweetCount)
		}

		tweets = append(tweets, t...)
		searchURLPath = nextSearchURLPath
		if nextSearchURLPath == "" {
			break
		}
		// Ensure the search URL is with the proper tweet_mode etc.
		searchURLPath += "&tweet_mode=extended"
		searchURLPath += "&include_entities=false"
	}

	return nil
}

// collectTweetsByURL collects the tweets based on URL path.
func collectTweetsByURL(ctx context.Context, bearerToken, url string) ([]TweetV1, string, error) {
	// Create HTTP request.
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	// Use retryable HTTP library with automatic retries and exponential back-off.
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	client := retryClient.StandardClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code (%d) from %s",
			resp.StatusCode, url)
	}

	// Decode response body as JSON.
	val := &TwitterSearchResponseV1{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(val); err != nil {
		return nil, "", err
	}

	return val.Tweets, val.SearchMetadata.NextURLPath, nil
}
