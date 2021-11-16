# CS 410 Final Project - Brand Sentiment on Twitter using Sentiment Analysis (Fall 2021)

## Introduction

_Sentiment analysis_ can capture the market or customer sentiment towards a brand. Companies can use this information to better understand their audiences’ reactions to the brand’s news or marketing campaigns, and to further enhance the brand. Investors or traders can also leverage this information to determine whether they should long or short their positions in the stock behind the brand.

This project is to perform _sentiment analysis_ on the Twitter tweets related to a given brand over a period of time, and create a sentiment trend graph which visualizes the sentiment towards the brand.

## Tools

There are several tools developed for this project:
1. _tweetscollect_ for collecting the tweets for a topic from _Twitter_ into a dataset.
2. _sentimentalyze_ for performing _sentiment analysis_ on the dataset.
3. TBD for creating a sentiment trend graph based on the analyzed data in the dataset.

### Prerequisites

To build and run the tools, there are several prerequisites:
1. You will need to install [Go 1.17](https://golang.org/doc/install).
2. If you don't have a _Twitter developer account_,  [apply one](https://developer.twitter.com/en/apply-for-access). Once you have the account, you will need to create a _Bearer Token_ for authentication. Please see [How to generate from the developer portal](https://developer.twitter.com/en/docs/authentication/oauth-2-0/bearer-tokens).
3. If you don't have an _AWS account_, [apply one](https://aws.amazon.com). Once you have the account, you will need to [create an access key](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys) for programmatic access. The _access key_ consists of an _access key ID_ and a _secret access key_.

### tweetscollect

_tweetscollect_ is a tool for collecting the tweets for a topic from _Twitter_ into a dataset.
```
# Build tweetscollect into an executable.
#
$ go build ./cmd/tweetscollect/...

# Run tweetscollect to collect tweets for a topic, by
# using the bearer token from your Twitter developer
# account, and write the collected tweets to a file.
#
$ tweetscollect -b "<bearer-token>" -t "<topic>" -o <output-file>
Collecting tweets from Twitter on topic "<topic>" ...

{ "date": "2021-11-15 17:56:14 +0000 +0000", "text": "JUST IN: Ohio Attorney General sues Facebook (Meta) for securities fraud.", "lang": "en", "favorite": 1077, "retweet": 328 }
...
Writing collected tweets to <output-file> ...
Done.
```

### sentimentalyze

_sentimentalyze_ is a tool for performing _sentiment analysis_ over the dataset using _Amazon Comprehend_.

```
# Build sentimentalyze into an executable.
#
$ go build ./cmd/sentimentalyze/...

# Run _sentimentalyze_ to perform sentiment analysis, by
# using the access key ID and secret access key from your
# AWS account.
#
$ ./sentimentalyze -a <access-key-id> -s <secret-access-key>
```

### TBD

TBD for creating a sentiment trend graph.