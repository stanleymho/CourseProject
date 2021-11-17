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
	cleanedTweets := make(map[string]string)
	for _, tweet := range data.Tweets {
		key := tweet.Key()
		if _, ok := uniqueTweets[key]; !ok {
			uniqueTweets[key] = types.SentimentTypeNeutral
			cleanedTweets[key] = tweet.cleanText()
		}
	}

	fmt.Printf("Normalizing %d tweets into %d unique tweets ...\n", len(data.Tweets), len(uniqueTweets))
	fmt.Printf("Performing sentiment analysis on the unique tweets ...\n")

	counter := 0
	textList := make([]string, 0, 25)
	for tweet, _ := range uniqueTweets {
		cleanText := cleanedTweets[tweet]
		textList = append(textList, cleanText)

		if counter%25 == 24 || counter == len(uniqueTweets)-1 || counter == 27 {
			sentimentMap, err := analyzeSentimentInText(ctx, region, accessKeyID, secretAccessKey, "en", textList)
			if err != nil {
				return err
			}

			for text, sentiment := range sentimentMap {
				fmt.Printf("{ \"text\": \"%s\", \"sentiment\": \"%v\" }\n", displayText(text), sentiment)
				uniqueTweets[tweet] = sentiment
			}
			textList = nil
		}
		counter++

		// Remove for production.
		if counter > 27 {
			break
		}
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
		sentiment := uniqueTweets[tweet.Key()]
		seperator := ","
		if i == len(data.Tweets)-1 {
			seperator = ""
		}

		f.WriteString(fmt.Sprintf("\t\t{ \"date\": \"%v\", \"text\": \"%s\", \"lang\": \"%s\", \"favorite\": %d, \"retweet\": %d, \"sentiment\": \"%s\" }%s\n",
			date, displayText(tweet.Text), tweet.Lang, tweet.FavoriteCount, tweet.RetweetCount, sentiment, seperator))
	}
	f.WriteString("\t]\n")
	f.WriteString("}\n")

	fmt.Printf("Done.\n")

	return nil
}

// analyzeSentimentInText performs sentiment analysis on a text using Amazon Comprehend.
func analyzeSentimentInText(ctx context.Context, region, accessKeyID, secretAccessKey,
	lang string, textList []string) (map[string]types.SentimentType, error) {

	result := make(map[string]types.SentimentType)

	cfg := &aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
	}

	comprehendClient := comprehend.NewFromConfig(*cfg)
	input := &comprehend.BatchDetectSentimentInput{}
	//	input := &comprehend.DetectSentimentInput{}
	input.LanguageCode = types.LanguageCode(lang)
	for _, text := range textList {
		input.TextList = append(input.TextList, text)
	}
	output, err := comprehendClient.BatchDetectSentiment(ctx, input)
	if err != nil {
		return result, err
	}

	/*
		fmt.Printf("*** %v\n", output)
		fmt.Printf("  *** Sentiment: %v\n", output.Sentiment)
		fmt.Printf("  *** Positive: %v\n", *output.SentimentScore.Positive)
		fmt.Printf("  *** Mixed: %v\n", *output.SentimentScore.Mixed)
		fmt.Printf("  *** Negative: %v\n", *output.SentimentScore.Negative)
		fmt.Printf("  *** Neutral: %v\n", *output.SentimentScore.Neutral)
	*/
	for _, r := range output.ResultList {
		result[textList[*r.Index]] = r.Sentiment
	}

	return result, nil
}
