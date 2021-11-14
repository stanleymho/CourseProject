package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
)

// runSentimentAnalysis performs sentiment analysis using Amazon Comprehend.
func runSentimentAnalysis(ctx context.Context, region, accessKeyID, secretAccessKey, lang, text string) error {
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
		return err
	}

	fmt.Printf("*** %v\n", output)
	fmt.Printf("  *** Sentiment: %v\n", output.Sentiment)
	fmt.Printf("  *** Positive: %v\n", *output.SentimentScore.Positive)
	fmt.Printf("  *** Mixed: %v\n", *output.SentimentScore.Mixed)
	fmt.Printf("  *** Negative: %v\n", *output.SentimentScore.Negative)
	fmt.Printf("  *** Neutral: %v\n", *output.SentimentScore.Neutral)
	return nil
}
