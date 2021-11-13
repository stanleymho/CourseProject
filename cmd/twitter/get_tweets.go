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

// getTweets gets the tweets for a topic.
func getTweets(bearerToken, topic string) error {
	ctx := context.Background()

	// Constuct URL.
	var sb strings.Builder
	sb.WriteString("?lang=en")
	sb.WriteString("&result_type=popular")
	sb.WriteString("&tweet_mode=extended")
	sb.WriteString("&include_entities=false")
	sb.WriteString("&q=" + url.QueryEscape(topic))
	sb.WriteString("&count=100")
	searchURLPath := sb.String()

	tweets := make([]Tweet, 0, 100)
	for {
		url := fmt.Sprintf("%s%s", twitterSearchURL, searchURLPath)
		t, nextSearchURLPath, err := getTweetsByURL(ctx, bearerToken, url)
		if err != nil {
			return err
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

// getTweetsByURL gets the tweets based on URL path.
func getTweetsByURL(ctx context.Context, bearerToken, url string) ([]Tweet, string, error) {
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
	val := &TwitterSearchResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(val); err != nil {
		return nil, "", err
	}

	for _, tweet := range val.Tweets {
		fmt.Printf("*** %s\n", tweet.NormalizedFullText())
		t, _ := parseDate(tweet.CreatedAt)
		fmt.Printf("---> date=%v, favorite=%d, retweet=%d\n", t, tweet.FavoriteCount, tweet.RetweetCount)
	}

	return val.Tweets, val.SearchMetadata.NextURLPath, nil
}
