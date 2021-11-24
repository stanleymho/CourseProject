package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// collectTweets collects the tweets for a topic.
func collectTweets(bearerToken, topic, outputFile string) error {
	fmt.Printf("Collecting tweets from Twitter on topic \"%s\" ...\n\n", topic)

	// Constuct URL.
	var sb strings.Builder
	sb.WriteString("?lang=en")
	sb.WriteString("&result_type=mixed")
	sb.WriteString("&tweet_mode=extended")
	sb.WriteString("&include_entities=false")
	sb.WriteString("&count=100")
	sb.WriteString("&q=" + url.QueryEscape(topic))
	searchURLPath := sb.String()

	ctx := context.Background()
	collectedTweets := make([]TweetV1, 0, 100)
	for {
		url := fmt.Sprintf("%s%s", twitterSearchURL, searchURLPath)
		tweets, nextSearchURLPath, err := collectTweetsByURL(ctx, bearerToken, url)
		if err != nil {
			return err
		}

		collectedTweets = append(collectedTweets, tweets...)

		for _, tweet := range tweets {
			createdAt, _ := parseDate(tweet.CreatedAt)
			text := tweet.NormalizedText()
			fmt.Printf("{ \"date\": \"%v\", \"text\": \"%s\", \"lang\": \"%s\", \"favorite\": %d, \"retweet\": %d }\n",
				createdAt, text, tweet.Lang, tweet.FavoriteCount, tweet.RetweetCount)
		}

		if nextSearchURLPath == "" {
			break
		}

		// Ensure the search URL is with the proper tweet_mode etc.
		searchURLPath = nextSearchURLPath
		searchURLPath += "&tweet_mode=extended"
		searchURLPath += "&include_entities=false"
	}

	fmt.Printf("\nWriting collected tweets to %s ...\n", outputFile)

	// Create output file.
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Output all tweets to the file.
	f.WriteString("{\n")
	f.WriteString("\t\"data\": [\n")

	for i, tweet := range collectedTweets {
		createdAt, _ := parseDate(tweet.CreatedAt)
		text := tweet.NormalizedText()
		seperator := ","
		if i == len(collectedTweets)-1 {
			seperator = ""
		}
		f.WriteString(fmt.Sprintf("\t\t{ \"date\": \"%v\", \"text\": \"%s\", \"lang\": \"%s\", \"favorite\": %d, \"retweet\": %d }%s\n",
			createdAt, text, tweet.Lang, tweet.FavoriteCount, tweet.RetweetCount, seperator))
	}
	f.WriteString("\t]\n")
	f.WriteString("}\n")

	fmt.Printf("Done.\n")

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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// If we hit rate limit, simply return.
		return make([]TweetV1, 0, 0), "", nil
	} else if resp.StatusCode != http.StatusOK {
		// If we hit other unexpected status code, returns error.
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
