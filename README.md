# CS 410 Final Project - Brand Sentiment on Twitter using Sentiment Analysis (Fall 2021)

## Introduction

_Sentiment analysis_ can capture the market or customer sentiment towards a brand. Companies can use this information to better understand their audiences’ reactions to the brand’s news or marketing campaigns, and to further enhance the brand. Investors or traders can also leverage this information to determine whether they should long or short their positions in the stock behind the brand.

This project is to perform _sentiment analysis_ on the Twitter tweets related to a given brand over a period of time, and create a sentiment trend graph which visualizes the sentiment towards the brand.

There are several tools developed for this project: _tweetscollect_, _sentimentalyze_, etc.

## tweetscollect

_tweetscollect_ is a tool for collecting the tweets for a topic from _Twitter_ into a dataset.

To build _tweetscollect_,
```
$ go build ./cmd/tweetscollect/...
```

To run _tweetscollect_,
```
$ ./tweetscollect -b <bearer-token> -t <topic>
```

## sentimentalyze

_sentimentalyze_ is a tool for performing _sentiment analysis_ over the dataset using _Amazon Comprehend_.

To build _sentimentalyze_,
```
$ go build ./cmd/sentimentalyze/...
```

To run _sentimentalyze_,
```
$ ./sentimentalyze -a <access-key-id> -s <secret-access-key>
```

## TBD

TBD for creating a sentiment trend graph.