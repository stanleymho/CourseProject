#!/usr/bin/python3
#
import matplotlib.pyplot as plt
import json
import sys
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

def main():
    if len(sys.argv) < 2:
        print("Usage:\n  {} <sentiment-file>\n".format(sys.argv[0]))
        sys.exit(1)

    plt.title('Sentiment Trend Graph')
    print("Loading tweets data with sentiment, it will take a minute or two ...")
    dictionary = json.load(open(sys.argv[1], 'r'))
    tweetList = dictionary['data']

    # Compute running average.
    runningAverageNeutral = 0.0
    runningAveragePositive = 0.0
    runningAverageNegative = 0.0
    for i in reversed(range(len(tweetList))):
        if i >= 100:
            break
        tweet = tweetList[i]
        s = tweet['sentiment']
        if s == 'POSITIVE':
            runningAveragePositive += 1
        elif s == 'NEGATIVE':
            runningAverageNegative += 1
        elif s == 'NEUTRAL' or s == 'MIXED':
            runningAverageNeutral += 1

    runningAveragePositive /= 100
    runningAverageNegative /= 100
    runningAverageNeutral /= 100

    counter = 0
    for i in reversed(range(len(tweetList))):
        counter += 1
        if counter % 2000 == 1999:
            print("Loaded {} tweets data with sentiment ...".format(counter + 1))
        tweet = tweetList[i]
        # Parse date, e.g. "2021-11-17T09:06:13Z"
        dt = datetime.strptime(tweet['date'], "%Y-%m-%dT%H:%M:%SZ")
        s = tweet['sentiment']
        plt.plot(dt, sentiment_value(s), 
            sentiment_marker(s),
            color=sentiment_color(s))

    print("Finished loading {} tweets data with sentiment!".format(len(tweetList)))
    print("Plotting sentiment trend graph ...")
    plt.style.use('seaborn-whitegrid')
    plt.xlabel('Date')
    plt.ylabel('Sentiment')
    plt.xticks(rotation='60')
    # plt.margins
    plt.show()

if __name__ == '__main__':
    main()
