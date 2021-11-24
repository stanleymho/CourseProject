#!/usr/bin/python3
#
import matplotlib.pyplot as plt
import json
from datetime import datetime

def sentiment_value(s):
    if s == 'POSITIVE':
        return 1
    elif s == 'NEGATIVE':
        return -1
    elif s == 'NEUTRAL' or s == 'MIXED':
        return 0

def sentiment_marker(s):
    if s == 'POSITIVE':
        return '^'
    elif s == 'NEGATIVE':
        return 'v'
    elif s == 'NEUTRAL' or s == 'MIXED':
        return '.'

def sentiment_color(s):
    if s == 'POSITIVE':
        return 'blue'
    elif s == 'NEGATIVE':
        return 'red'
    elif s == 'NEUTRAL' or s == 'MIXED':
        return 'green'

plt.title('Sentiment Trend Graph')
print("Loading tweets data with sentiment, it will takes a minute or two ...")
dictionary = json.load(open('../cmd/sentimentalyze/tweets-sentiment.json', 'r'))
tweetList = dictionary['data']
counter = 0
for i in reversed(range(len(tweetList))):
    counter += 1
    if counter % 2000 == 1999:
        print("Loaded {} tweets data with sentiment ...".format(counter + 1))
    tweet = tweetList[i]
    # Parse date, e.g. "2021-11-17T09:06:13Z"
    dt = datetime.strptime(tweet['date'], "%Y-%m-%dT%H:%M:%SZ")
    sentiment = tweet['sentiment']
    plt.plot(dt, sentiment_value(sentiment), 
        sentiment_marker(sentiment),
        color=sentiment_color(sentiment))

print("Finished loading {} tweets data with sentiment!".format(len(tweetList)))
print("Plotting sentiment trend graph ...")
plt.style.use('seaborn-whitegrid')
plt.xlabel('Date')
plt.ylabel('Sentiment')
plt.xticks(rotation='60')
# plt.margins
plt.show()
