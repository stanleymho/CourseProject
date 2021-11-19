# CS 410 Final Project - Brand Sentiment on Twitter using Sentiment Analysis (Fall 2021)

_Sentiment analysis_ can capture the market or customer sentiment towards a brand. Companies can use this information to better understand their audiences’ reactions to the brand’s news or marketing campaigns, and to further enhance the brand. Investors or traders can also leverage this information to determine whether they should long or short their positions in the stock behind the brand.

This project is to perform _sentiment analysis_ on the Twitter tweets related to a given brand over a period of time, and create a sentiment trend graph which visualizes the sentiment towards the brand.

# Tools

There are several tools developed for this project:
1. _tweetscollect_ for collecting the tweets for a topic from _Twitter_ for the past 7 days into a dataset.
2. _sentimentalyze_ for performing _sentiment analysis_ on the dataset.
3. TBD for creating a sentiment trend graph based on the analyzed data in the dataset.

## Prerequisites

There are several prerequisites for building and running the tools:
1. You will need to install [Go 1.17](https://golang.org/doc/install).
2. If you don't have a _Twitter developer account_,  [apply one](https://developer.twitter.com/en/apply-for-access). Once you have the account, you will need to create a _Bearer Token_ for authentication. Please see [How to generate from the developer portal](https://developer.twitter.com/en/docs/authentication/oauth-2-0/bearer-tokens).
3. If you don't have an _AWS account_, [apply one](https://aws.amazon.com). Once you have the account, you will need to [create an access key](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys) for programmatic access. The _access key_ consists of an _access key ID_ and a _secret access key_.

## 1. tweetscollect

### 1.1. Description

_tweetscollect_ is a tool for collecting the tweets for a topic from _Twitter_ for the past 7 days into a dataset. The topic could be one or more words. If the word contains special characters, e.g. `$`, the character should be escaped with `\`, i.e. `\$`.

### 1.2. Implementation

_tweetscollect_ uses the [Twitter's standard search API](https://developer.twitter.com/en/docs/twitter-api/v1/tweets/search/api-reference/get-search-tweets to query against a mixture of the recent and popular tweets for the past 7 days for a given topic. Each tweet in the returned result is then reduced to the mininal, and it includes the date, text, language, favorite count, and retweeted count. Each API call returns limited number of tweets, and multiple paginations are involved in order to collect all the tweets across 7 days. After all the tweets are collected, they are written out to a file in json format.

### 1.3. Usage

To run _tweetscollect_, you must have the _bearer token_ from a _Twitter developer account_.

Notice that the _Twitter developer account_ has rate limit on the maximum number of requests allowed in a 15-minutes time window, and collecting the tweets for one topic alone might get very close to the limit. Hence, in order to use the tool successfully, please run the tool at most once in a 15-minutes time window.

```
# Build tweetscollect into an executable.
#
$ go build ./cmd/tweetscollect/...

# Run tweetscollect to collect tweets for a topic, by using the bearer token
# from your Twitter developer account, and write the collected tweets to file.
#
$ ./tweetscollect -b "<bearer-token>" -t "<topic>" -o <output-file>
Collecting tweets from Twitter on topic "<topic>" ...

{ "date": "2021-11-16T00:46:05Z", "text": "JUST IN: Ohio Attorney General sues Facebook (Meta) for securities fraud.", "lang": "en", "favorite": 1077, "retweet": 328 }
...

Writing collected tweets to <output-file> ...
Done.
```

## 2. sentimentalyze

### 2.1. Description

_sentimentalyze_ is a tool for performing _sentiment analysis_ over the dataset which _tweetscollect_ has collected. 

### 2.2. Implementation

_sentimentalyze_ first normalizes all the tweets from the dataset, as there are many retweets and dedupling the tweets could significantly reduce the unique number of tweets for sentiment analysis. Afterwards, _sentimentalyze_ sends the tweets to _Amazon Comprehend_ in batches to determine sentiment in the tweets. After all the tweets have been analyzed, the tweets and their sentiment information are written out to a file in json format.

### 2.3. Usage

To run _sentimentalyze_, you must have the access key ID and secret access key from an _AWS account_.

Please be aware that because _sentimentalyze_ will use _Amazon Comprehend_ from the _AWS account_, and the _AWS account_ will be charged for usage. On average, each run involves between 40,000 to 60,000 tweets, and that's approximately 15,000 to 25,000 unique tweets which costs $1.5 to $2.5 to perform a single run of _sentiment analysis_.

```
# Build sentimentalyze into an executable.
#
$ go build ./cmd/sentimentalyze/...

# Run sentimentalyze to perform sentiment analysis on the tweets in the
# input file, using the access key ID and secret access key from the
# AWS account.
#
$ ./sentimentalyze -i <input-file> -o <output-file> -a <access-key-id> -s <secret-access-key> [-r <region>]
Reading 40655 tweets from ../tweetscollect/tweets.json ...
Normalizing 40655 tweets into 13089 unique tweets ...
Performing sentiment analysis on the unique tweets ...
{ "text": "Insight into the big rebrand.. 😎\nvia \n#Facebook #AR #VR #meta #Metaverse #virtualworlds\n\nhttps://t.co/bMQ4hhKlTW", "sentiment": "NEUTRAL" }
{ "text": "Facebook owner Meta has opened up more about the amount of bullying and harassment on its platforms amid pressure to increase transparency\n\n✍️:", "sentiment": "NEUTRAL" }
...
Writing tweets with analyzed sentiment to tweets-sentiment.json ...
Done.
```

## 3. TBD

TBD for creating a sentiment trend graph.