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
	uniqueTweets := make(map[string]types.SentimentType)
	for _, tweet := range data.Tweets {
		key := tweet.Key()
		if _, ok := uniqueTweets[key]; !ok {
			uniqueTweets[key] = types.SentimentTypeNeutral
		}
	}

	fmt.Printf("Normalizing %d tweets into %d unique tweets ...\n", len(data.Tweets), len(uniqueTweets))
	fmt.Printf("Performing sentiment analysis on the unique tweets ...\n")

	counter := 0
	for tweet, _ := range uniqueTweets {
		// sentiment, err := analyzeSentimentInText(ctx, region, accessKeyID, secretAccessKey, t.Lang, text)
		// if err != nil {
		//	  return err

		sentiment := types.SentimentTypeNeutral
		if counter%4 == 0 {
			sentiment = types.SentimentTypeNeutral
		} else if counter%4 == 1 {
			sentiment = types.SentimentTypePositive
		} else if counter%4 == 2 {
			sentiment = types.SentimentTypeNegative
		} else if counter%4 == 3 {
			sentiment = types.SentimentTypeMixed
		}

		uniqueTweets[tweet] = sentiment
		counter++
	}

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
		date := tweet.Date
		text := tweet.NormalizedText()
		sentiment := uniqueTweets[tweet.Key()]
		seperator := ","
		if i == len(data.Tweets)-1 {
			seperator = ""
		}

		f.WriteString(fmt.Sprintf("\t\t{ \"date\": \"%v\", \"text\": \"%s\", \"lang\": \"%s\", \"favorite\": %d, \"retweet\": %d, \"sentiment\": \"%s\" }%s\n",
			date, text, tweet.Lang, tweet.FavoriteCount, tweet.RetweetCount, sentiment, seperator))
	}
	f.WriteString("\t]\n")
	f.WriteString("}\n")

	fmt.Printf("Done.\n")

	return nil
}

// analyzeSentimentInText performs sentiment analysis on a text using Amazon Comprehend.
func analyzeSentimentInText(ctx context.Context, region, accessKeyID, secretAccessKey,
	lang, text string) (types.SentimentType, error) {

	cfg := &aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
	}

	comprehendClient := comprehend.NewFromConfig(*cfg)
	input := &comprehend.DetectSentimentInput{}
	input.LanguageCode = types.LanguageCode(lang)
	input.Text = &text
	output, err := comprehendClient.DetectSentiment(ctx, input)
	if err != nil {
		return types.SentimentTypeNeutral, err
	}

	/*
		fmt.Printf("*** %v\n", output)
		fmt.Printf("  *** Sentiment: %v\n", output.Sentiment)
		fmt.Printf("  *** Positive: %v\n", *output.SentimentScore.Positive)
		fmt.Printf("  *** Mixed: %v\n", *output.SentimentScore.Mixed)
		fmt.Printf("  *** Negative: %v\n", *output.SentimentScore.Negative)
		fmt.Printf("  *** Neutral: %v\n", *output.SentimentScore.Neutral)
	*/
	return output.Sentiment, nil
}
