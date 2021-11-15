# CS 410 Final Project - Brand Sentiment on Twitter using Sentiment Analysis (Fall 2021)

## Introduction

_Sentiment analysis_ can capture the market or customer sentiment towards a brand. Companies can use this information to better understand their audiences’ reactions to the brand’s news or marketing campaigns, and to further enhance the brand. Investors or traders can also leverage this information to determine whether they should long or short their positions in the stock behind the brand.

This project is to perform _sentiment analysis_ on the Twitter tweets related to a given brand over a period of time, and create a sentiment trend graph which visualizes the sentiment towards the brand.

There are several tools developed for this project:
1. _tweetscollect_ for collecting the tweets for a topic from _Twitter_ into a dataset.
2. _sentimentalyze_ for performing _sentiment analysis_ on the dataset.
3. TBD for creating a sentiment trend graph based on the analyzed data in the dataset.

## tweetscollect

_tweetscollect_ is a tool for collecting the tweets for a topic from _Twitter_ into a dataset. To use _tweetscollect_, you must have a [Twitter Developer Account](https://developer.twitter.com/en/apply-for-access) and a _Bearer Token_ for authentication. 

```
# Build tweetscollect into an executable.
#
$ go build ./cmd/tweetscollect/...

# Run tweetscollect to collect tweets for a topic.
#
$ ./tweetscollect -b <bearer-token> -t <topic>
```

## sentimentalyze

_sentimentalyze_ is a tool for performing _sentiment analysis_ over the dataset using _Amazon Comprehend_.

```
# Build sentimentalyze into an executable.
#
$ go build ./cmd/sentimentalyze/...

# Run _sentimentalyze_ to perform sentiment analysis.
#
$ ./sentimentalyze -a <access-key-id> -s <secret-access-key>
```

## TBD

TBD for creating a sentiment trend graph.