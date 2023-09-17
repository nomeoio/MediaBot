package mediabot

import (
	"fmt"
	"time"
)

var mBot MediaBot

var TweetLists = map[string]string{}

func StartAll(slackWebHookUrlHN, slackWebHookUrlTwitter, slackWebHookUrlCartoons, slackWebHookUrlHNClassics, logUrl string) {
	if !IsTestMode { // if this is not in test mode
		var i int = 0
		for {
			SC.SendPlainText(fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"), ":", "Auto retrieving new posts... "), logUrl)
			go mBot.AutoRetrieveHN(slackWebHookUrlHN)
			// go rc.AutoRetrieveNew()
			go mBot.AutoRetrieveTwitter(TweetLists, 2000, 1000, Params.TwitterBearerToken, slackWebHookUrlTwitter)
			if i%12 == 0 { // run every 12 hours
				go mBot.AutoRetrieveXkcd(slackWebHookUrlCartoons)
			}
			// if i%24 == 0 { // run every 24 hours
			// 	go mBot.AutoRetrieveHNClassic(slackWebHookUrlHNClassics)
			// 	i = 0
			// }
			i++
			time.Sleep(time.Hour)
		}
	}
}

// func TestAll(t *testing.T) {
// 	DB.Init("file:./data/ids.db") // creates db file or table if doesn't exist, doesn't do anything if exists.
// 	// var savedItems []SavedItem = DB.ReturnAllRecords("HackerNews")
// 	// var savedItems []SavedItem = DB.ReturnAllRecords("Twitter")
// 	// t.Log(savedItems)
// 	// for _, savedItem := range savedItems {
// 	// 	t.Log(savedItem)
// 	// }
// 	// testTwitter()
// 	mBot.AutoRetrieveXkcd(slackWebHookUrlTest)
// }

// func testHN() {
// 	if err := mBot.AutoRetrieveHN(slackWebHookUrlTest); err != nil {
// 		log.Fatalln(err)
// 	}
// }

// func testTwitter() {
// 	if err := mBot.AutoRetrieveTwitter(TweetLists, 2000, 1000, Params.TwitterBearerToken, slackWebHookUrlTest); err != nil {
// 		log.Panic(err)
// 	}
// }
