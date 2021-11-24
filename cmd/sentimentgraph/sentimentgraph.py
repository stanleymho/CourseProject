#!/usr/bin/python3
#
import matplotlib.pyplot as plt
import json
import sys
from datetime import datetime
import warnings
warnings.filterwarnings("ignore")

def main():
    if len(sys.argv) < 2:
        print("Usage:\n  {} <sentiment-file>\n".format(sys.argv[0]))
        sys.exit(1)

    print("Loading sentiment data from {}, it will take a minute or two ...".format(sys.argv[1]))
    dictionary = json.load(open(sys.argv[1], 'r'))
    tweetList = dictionary['data']
    plt.title("Sentiment Trend Graph: {}".format(dictionary['query']))

    totalCount = 0
    positiveCount = 0
    negativeCount = 0
    mixedCount = 0
    neutralCount = 0
    dateTimePoints = []
    positivePoints = []
    negativePoints = []
    mixedPoints = []
    neutralPoints = []
    minDateTime = datetime.now()
    maxDateTime = datetime(1900, 1, 1)
    for i in reversed(range(len(tweetList))):
        totalCount += 1
        if totalCount % 2000 == 0:
            print("Loaded {} sentiment data ...".format(totalCount))
        tweet = tweetList[i]

        # Parse date, e.g. "2021-11-17T09:06:13Z"
        dt = datetime.strptime(tweet['date'], "%Y-%m-%dT%H:%M:%SZ")
        if minDateTime > dt:
            minDateTime = dt
        if maxDateTime < dt:
            maxDateTime = dt

        s = tweet['sentiment']
        if s == "POSITIVE":
            positiveCount += 1
        elif s == "NEGATIVE":
            negativeCount += 1
        elif s == "MIXED":
            mixedCount += 1
        else:
            neutralCount +=1


        dateTimePoints.append(dt)
        positivePoints.append(positiveCount * 1.0 / totalCount)
        negativePoints.append(negativeCount * 1.0 / totalCount)
        mixedPoints.append(mixedCount * 1.0 / totalCount)
        neutralPoints.append(neutralCount * 1.0 / totalCount)

    print("Finished loading {} sentiment data!".format(len(tweetList)))
    print("Plotting sentiment trend graph ...")
    plt.grid(True)
    plt.xlabel('Date')
    plt.ylabel('Sentiment')
    plt.xlim([minDateTime, maxDateTime])
    plt.xticks(rotation='60')
    current_values = plt.gca().get_yticks()
    plt.gca().set_yticklabels(['{:,.0%}'.format(x) for x in current_values])
    plt.stackplot(dateTimePoints, positivePoints, mixedPoints, negativePoints, neutralPoints,
        labels=['Positive','Mixed', 'Negative', 'Neutral'],
        colors=['tab:blue', 'gold', 'darkorange', 'grey'])
    plt.legend(loc='upper right')
    plt.tight_layout()
    plt.show()

if __name__ == '__main__':
    main()
