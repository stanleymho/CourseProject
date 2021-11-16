package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
)

// analyzeSentiment
func analyzeSentiment(ctx context.Context, inputFile, outputFile, region, accessKeyID,
	secretAccessKey string) error {

	f, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return err
	}

	data := TweetsData{}
	if err = json.Unmarshal([]byte(f), &data); err != nil {
		return err
	}

	fmt.Printf("Reading %d tweets from %s ...\n", len(data.Tweets), inputFile)

	// Process regular tweets first.
	sentimentMap := make(map[string]types.SentimentType)
	counter := 0
	for _, tweet := range data.Tweets {
		// Skip retweets.
		if tweet.IsRetweeted() {
			continue
		}
		text := tweet.NormalizedText()
		//		sentiment, err := analyzeSentimentInText(ctx, region, accessKeyID, secretAccessKey, t.Lang, text)
		//if err != nil {
		//	return err
		//}
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

		fmt.Printf("%v: \"%s\"\n", sentiment, text)
		counter++
		key := text
		if len(text) > 100 {
			key = text[:100]
		}
		sentimentMap[key] = sentiment
	}

	// Process retweets first.
	for _, tweet := range data.Tweets {
		// Skip regular tweets.
		if !tweet.IsRetweeted() {
			continue
		}

		text := tweet.NormalizedText()
		key := text
		if len(text) > 100 {
			key = text[:100]
		}

		if _, ok := sentimentMap[key]; !ok {
			fmt.Printf("Sentiment not found for retweet: %s\n", text)
		}
	}

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
