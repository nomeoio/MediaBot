package mediabot

import (
	"fmt"
	"log"
	"os"

	"github.com/naughtymonsta/utilities"
)

var SC = utilities.SlackClient{
	IsTestMode: false,
}

type MediaBot struct{}

func (mBot MediaBot) AutoRetrieveHN(slackWebHookUrlHN string) (err error) {
	var totalNews int = 0
	for _, s := range []string{"top", "new", "best"} {
		var mbss []utilities.MessageBlocks
		if mbss, err = hn.RetrieveNew(s, 200); err != nil {
			return
		}
		totalNews += len(mbss)
		for _, mbs := range mbss {
			if err = SC.SendBlocks(mbs, slackWebHookUrlHN); err != nil { // send the new and not published stories to slack #hacker-news
				return
			}
		}
	}
	log.Println("retrieved", totalNews, "hacker news.")
	return
}

// func (mBot MediaBot) AutoRetrieveHNClassic(slackWebHookUrlHNClassics string) (err error) {
// 	var results HNAlgoliaSearchResults
// 	if results, err = hn.RetrieveHNClassic(); err != nil {
// 		return
// 	}
// 	var mrkdwnList []string = hn.ClassicsFormatData(results)

// 	for _, text := range mrkdwnList {
// 		var mbarr = []utilities.MessageBlock{
// 			{Type: "divider"},
// 			SC.CreateTextBlock(text, "mrkdwn", ""),
// 		}
// 		if err = SC.SendBlocks(utilities.MessageBlocks{Blocks: mbarr}, slackWebHookUrlHNClassics); err != nil { // send the new and not published stories to slack #hacker-news
// 			return
// 		}
// 	}

// 	// mbarr = append(mbarr, SC.CreateTextBlock(text, "mrkdwn", ""))
// 	// if err = sc.SendBlocks(MessageBlocks{Blocks: mbarr}, os.Getenv("SlackWebHookUrlHNClassics")); err != nil { // send the new and not published stories to slack #hacker-news
// 	// 	return
// 	// }
// 	return
// }

func (mBot MediaBot) AutoRetrieveTwitter(tweetList map[string]string, leastOriginalLikes, leastRetweetLikes int, twitterBearerToken, slackWebHookUrlTwitter string) (err error) {
	// :leastOriginalLikes: if it's retweet, this checks the original tweet's likes
	var totalTweets int
	for listName, listId := range tweetList {
		var mbList [][]utilities.MessageBlock
		mbList, err = tc.RetrieveTweets(listId, twitterBearerToken, leastRetweetLikes, leastOriginalLikes)
		if err != nil {
			return
		}
		totalTweets += len(mbList)
		var i int
		var mb []utilities.MessageBlock
		for i, mb = range mbList {
			var mbs utilities.MessageBlocks
			if i == 0 {
				mbs.Blocks = []utilities.MessageBlock{{
					Type: "header",
					Text: &utilities.ElementText{
						Type: "plain_text",
						Text: listName,
					},
				}}
			}
			mbs.Blocks = append(mbs.Blocks, mb...)
			err = SC.SendBlocks(mbs, slackWebHookUrlTwitter)
			if err != nil {
				return
			}
		}
		if IsTestMode { // if in test mode, only go through 1 loop, instead going through all tweet lists
			break
		}
	}
	var text string = fmt.Sprintf("Sent %d tweets.", totalTweets)
	log.Println(text)
	SC.SendPlainText(text, os.Getenv("SlackWebHookUrlTest"))
	return
}

func (mBot MediaBot) AutoRetrieveXkcd(slackWebHookUrlCartoons string) (err error) {
	var item SavedItem = DB.UpdateXkcd()
	var mbs utilities.MessageBlocks
	mbs, err = xk.GetStoryById(item.Id)
	if err != nil {
		return
	}
	err = SC.SendBlocks(mbs, slackWebHookUrlCartoons)
	if err != nil {
		return
	}
	return
}
