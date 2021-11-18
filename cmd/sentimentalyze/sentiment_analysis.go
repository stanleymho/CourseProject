package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
)

// analyzeSentiment
func analyzeSentiment(ctx context.Context, inputFile, outputFile, region, accessKeyID,
	secretAccessKey string) error {

	fileData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	data := TweetsData{}
	if err = json.Unmarshal([]byte(fileData), &data); err != nil {
		return err
	}

	fmt.Printf("Reading %d tweets from %s ...\n", len(data.Tweets), inputFile)

	// Process all the tweets to determine the normalized set of tweets to process.
	tweetSentimentMap := make(map[string]types.SentimentType)
	originalTweets := make(map[string]string)
	for _, tweet := range data.Tweets {
		key := hashKey(tweet.Text)
		if _, ok := tweetSentimentMap[key]; !ok {
			tweetSentimentMap[key] = types.SentimentTypeNeutral
			originalTweets[key] = tweet.Text
		}
	}

	fmt.Printf("Normalizing %d tweets into %d unique tweets ...\n", len(data.Tweets), len(tweetSentimentMap))
	fmt.Printf("Performing sentiment analysis on the unique tweets ...\n")

	// Create client for Amazon Comprehend.
	cfg := &aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
	}
	comprehendClient := comprehend.NewFromConfig(*cfg)

	counter := 0
	tweetsLimit := 5
	keyList := make([]string, 0, 25)
	tweetList := make([]string, 0, 25)
	for key, _ := range tweetSentimentMap {
		keyList = append(keyList, key)
		tweetList = append(tweetList, trimTweet(originalTweets[key]))

		if counter%25 == 24 || counter == len(tweetSentimentMap)-1 || counter == tweetsLimit {
			textSentimentMap, err := analyzeSentimentInTextList(ctx, comprehendClient, "en", tweetList)
			if err != nil {
				return err
			}

			for text, sentiment := range textSentimentMap {
				fmt.Printf("{ \"text\": \"%s\", \"sentiment\": \"%v\" }\n", displayText(text), sentiment)
				tweetSentimentMap[hashKey(text)] = sentiment
			}
			keyList = nil
			tweetList = nil
		}
		counter++

		// Remove for production.
		if counter > tweetsLimit {
			break
		}
	}

	/*
		for key, sentiment := range tweetSentimentMap {
			if sentiment == types.SentimentTypePositive || sentiment == types.SentimentTypeNegative {
				fmt.Printf("*** %s -> %v\n", displayText(key), sentiment)
			}
		}
	*/

	fmt.Printf("Writing tweets with analyzed sentiment to %s ...\n", outputFile)

	// Create output file.
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// Output all tweets to the file.
	f.WriteString("{\n")
	f.WriteString("\t\"data\": [\n")

	for i, tweet := range data.Tweets {
		sentiment := tweetSentimentMap[hashKey(tweet.Text)]
		seperator := ","
		if i == len(data.Tweets)-1 {
			seperator = ""
		}

		f.WriteString(fmt.Sprintf("\t\t{ \"date\": \"%v\", \"text\": \"%s\", \"lang\": \"%s\", \"favorite\": %d, \"retweet\": %d, \"sentiment\": \"%s\" }%s\n",
			tweet.Date, displayText(tweet.Text), tweet.Lang, tweet.FavoriteCount,
			tweet.RetweetCount, sentiment, seperator))
	}
	f.WriteString("\t]\n")
	f.WriteString("}\n")

	fmt.Printf("Done.\n")

	return nil
}

// analyzeSentimentInTextList performs sentiment analysis on a text list using Amazon Comprehend.
func analyzeSentimentInTextList(ctx context.Context, comprehendClient *comprehend.Client,
	lang string, textList []string) (map[string]types.SentimentType, error) {
	result := make(map[string]types.SentimentType)

	// Perform sentiment analysis with Amazon Comprehend.
	input := &comprehend.BatchDetectSentimentInput{}
	input.LanguageCode = types.LanguageCode(lang)
	for _, text := range textList {
		input.TextList = append(input.TextList, text)
	}
	output, err := comprehendClient.BatchDetectSentiment(ctx, input)
	if err != nil {
		return result, err
	}

	// Put the result together.
	for _, r := range output.ResultList {
		result[textList[*r.Index]] = r.Sentiment
	}
	return result, nil
}
