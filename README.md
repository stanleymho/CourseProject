# CS 410 Final Project - Brand Sentiment on Twitter using Sentiment Analysis (Fall 2021)

_Sentiment analysis_ can capture the market or customer sentiment towards a brand. Companies can use this information to better understand their audiences’ reactions to the brand’s news or marketing campaigns, and to further enhance the brand. Investors or traders can also leverage this information to determine whether they should long or short their positions in the stock behind the brand.

This project is to perform _sentiment analysis_ on the Twitter tweets related to a given brand over a period of time, and create a sentiment trend graph to visualize the sentiment towards the brand.

# Tools

There are several tools developed for this project:
1. _tweetscollect_ for collecting the tweets for a topic from _Twitter_ for the past 7 days into a dataset.
2. _sentimentalyze_ for performing _sentiment analysis_ on the dataset.
3. _sentimentgraph_ for plotting a _Sentiment Trend Graph_ based on the results from the _sentiment analysis_.

These tools are designed to work together as follows:

![Workflow](/diagrams/Final-Project-Workflow.png)

1. First, _tweetscollect_ is used to collect the tweets for a given topic from Twitter and write the tweets into a file (i.e. tweets.json).
2. Next, _sentimentalyze_ takes the tweets returned from _tweetscollect_, and performs _sentiment analysis_ on these tweets using _Amazon Comprehend_ and output the result into another file (i.e. sentiment.json).
3. Finally, _sentimentgraph_ takes the result from _sentimentalyze_ to create a _Sentiment Trend Graph_ for visualization.

**ATTENTION TO REVIEWERS: To better understand how to use these software, please watch the [Software Usage Tutorial Presentation](https://drive.google.com/file/d/1uAjvmu3oai6wEMxDIIw5moyfi71hFeWb/view?usp=sharing) or read the [Final Project Documentation](/Final-Project-Documentation.pdf) for more details.**



## Prerequisites

There are several prerequisites for building and running the tools:
1. You will need to install [Go 1.17](https://golang.org/doc/install).
2. You will need to install [Python 3.10](https://www.python.org/downloads/).
3. If you don't have a _Twitter developer account_,  [apply one](https://developer.twitter.com/en/apply-for-access). Once you have the account, you will need to create a _Bearer Token_ for authentication. Please see [How to generate from the developer portal](https://developer.twitter.com/en/docs/authentication/oauth-2-0/bearer-tokens).
4. If you don't have an _AWS account_, [apply one](https://aws.amazon.com). Once you have the account, you will need to [create an access key](https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys) for programmatic access. The _access key_ consists of an _access key ID_ and a _secret access key_.

## 1. tweetscollect

### 1.1. Description

_tweetscollect_ is a tool for collecting the tweets for a topic from _Twitter_ for the past 7 days into a dataset. The topic consists of one or more words. If a word contains special characters, e.g. `$`, the character must be escaped, i.e. `\$`.

### 1.2. Implementation

_tweetscollect_ uses the [Twitter's standard search API](https://developer.twitter.com/en/docs/twitter-api/v1/tweets/search/api-reference/get-search-tweets) to query against a mixture of the recent and popular tweets for the past 7 days for a given topic. Each API call returns a limited number of tweets, and multiple paginations are involved in order to collect all the tweets across 7 days. After all the tweets are collected, the tweets are sorted and further reduced to include only the date, text, language, favorite count, and retweeted count. At the end, the resulting tweets are written out to a file in json format based on [schema](/schema/tweets-schema.json).

### 1.3. Usage

To run _tweetscollect_, you must have the _bearer token_ from a _Twitter developer account_. Collecting tweets for the past 7-days involves retrieving tens of thousands of tweets from Twitter, and it will take a few minutes for the tool to run to completion. Please be patient!

Notice that the _Twitter developer account_ has rate limit on the maximum number of requests allowed in a 15-minutes time window, and collecting the tweets for one topic alone might get very close to the limit. Hence, in order to use the tool successfully, please run the tool at most once in a 15-minutes time window.

```
# Build tweetscollect into an executable.
#
$ go build ./cmd/tweetscollect/...

# Display tweetscollect usage.
#
$ ./tweetscollect
Usage:
  tweetscollect -b <bearer-token> -t <topic> [-o <output-file>] [flags]

Flags:
  -b, --bearer-token string   Bearer token
  -h, --help                  help for tweetscollect
  -o, --output-file string    Output file (default "tweets.json")
  -t, --topic string          Topic, e.g. Facebook
  -v, --verbose count         Increase verbosity. May be given multiple times.

# Run tweetscollect to collect tweets for a topic, by using the bearer token
# from your Twitter developer account, and write the collected tweets to file.
#
# $ tweetscollect -b "<bearer-token>" -t "<topic>" -o <output-file>
#
$ ./tweetscollect -b "..." -t "Facebook Meta" -o tweets.json
Collecting tweets from Twitter on topic "Facebook Meta" ...

{ "date": "2021-11-16T00:46:05Z", "text": "JUST IN: Ohio Attorney General sues Facebook (Meta) for securities fraud.", "lang": "en", "favorite": 1077, "retweet": 328 }
...

Writing collected tweets to tweets.json ...
Done.
```

## 2. sentimentalyze

### 2.1. Description

_sentimentalyze_ is a tool for performing _sentiment analysis_ over the dataset with tweets which _tweetscollect_ has collected.

### 2.2. Implementation

_sentimentalyze_ first normalizes all the tweets from the dataset, as there are many retweets and dedupling the tweets could significantly reduce the unique number of tweets for sentiment analysis. Afterwards, _sentimentalyze_ sends the unique tweets to _Amazon Comprehend_ in batches to determine the sentiment. After all the unique tweets have been analyzed, _sentimentalyze_ would reprocess each of the original tweets from the dataset, identify its associated unique tweet and sentiment, and eventually write out the original tweets along with their sentiment to a file in json format based on [schema](/schema/sentiment-schema.json).

### 2.3. Usage

To run _sentimentalyze_, you must have the access key ID and secret access key from an _AWS account_. Performing sentiment analysis involves sending all the tweets to _Amazon Comprehend_ in multiple batches to process, and it will takes a few minutes for the tool to run to completion. Please be patient!

Please be aware that since _sentimentalyze_ will use _Amazon Comprehend_ from the _AWS account_, the _AWS account_ will be charged for usage. On average, each run involves between 40,000 to 60,000 tweets, and that's approximately 15,000 to 25,000 unique tweets which costs $1.5 to $2.5 to perform a single run of _sentiment analysis_.

```
# Build sentimentalyze into an executable.
#
$ go build ./cmd/sentimentalyze/...

# Display sentimentalyze usage.
#
$ ./sentimentalyze
Usage:
  sentimentalyze -i <input-file> -o <output-file> -a <access-key-id> -s <secret-access-key> [-r <region>] [-l <lang>] [flags]

Flags:
  -a, --access-key-id string       Access key ID
  -h, --help                       help for sentimentalyze
  -i, --input-file string          Input file (default "tweets.json")
  -o, --output-file string         Output file (default "sentiment.json")
  -r, --region string              Region (default "us-east-1")
  -s, --secret-access-key string   Secret access key
  -v, --verbose count              Increase verbosity. May be given multiple times.

# Run sentimentalyze to perform sentiment analysis on the tweets in the
# input file, using the access key ID and secret access key from the
# AWS account.
#
# $ sentimentalyze -i <input-file> -o <output-file> -a "<access-key-id>" -s "<secret-access-key>" [-r <region>]
#
$ ./sentimentalyze -i tweets.json -o sentiment.json -a "..." -s "..." -r "us-east-1"
Reading 40655 tweets from tweets.json ...
Normalizing 40655 tweets into 13089 unique tweets ...
Performing sentiment analysis on the unique tweets ...
{ "text": "Insight into the big rebrand.. 😎\nvia \n#Facebook #AR #VR #meta #Metaverse #virtualworlds\n\nhttps://t.co/bMQ4hhKlTW", "sentiment": "NEUTRAL" }
{ "text": "Facebook owner Meta has opened up more about the amount of bullying and harassment on its platforms amid pressure to increase transparency\n\n✍️:", "sentiment": "NEUTRAL" }
...
Writing tweets with analyzed sentiment to sentiment.json ...
Done.
```

## 3. sentimentgraph

### 3.1. Description

_sentimentgraph_ is a tool for plotting the _Sentiment Trend Graph_ based on the _sentiment analysis_ which _sentimentalyze_ has performed.

### 3.2. Implementation

_sentimentgraph_ simply reads the result of the _sentiment analysis_ from a file, and uses matplotlib.pyplot library to create a _Sentiment Trend Graph_. The sentiment returned from the _sentiment analysis_ for each tweet is either `POSITIVE`, `NEGATIVE`, `NEUTRAL`, or `MIXED`. To plot the graph, _sentimentgraph_ creates a _stacked area graph_ to show how the percentage of positive, negative, neutral, and mixed sentiment have changed over time.

### 3.3. Usage

Notice that running sentimentgraph involves plotting a Sentiment Trend Graph with tens of thousands of data points, and it will take a few minutes to display the visualization. Please be patient!

```
# Go to sentimentgraph directory.
#
$ cd ./cmd/sentimentgraph/

# Show sentimentgraph usage.
#
$ ./sentimentgraph.py
Usage:
  sentimentgraph.py <sentiment-file>

# Run sentimentgraph.py to create the Sentiment Trend Graph.
#
$ ./sentimentgraph.py sentiment.json
Loading sentiment data from sentiment.json, it will take a minute or two ...
Loaded 2000 sentiment data ...
Loaded 4000 sentiment data ...
...
Loaded 36000 sentiment data ...
Finished loading 37742 sentiment data!
Plotting sentiment trend graph ...
```

