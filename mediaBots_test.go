package mediabot

import (
	"fmt"
	"testing"
	"time"
)

var mBot MediaBot

var TweetLists = map[string]string{
	"Makers":        "1229215345526722560",
	"Entrepreneurs": "1229216130662723584",
	"Greats":        "1310225357019074562",
	"Investors":     "1237393320378118149",
	"Physicists":    "1394817230630572034",
	"YouTubers":     "1229243949950201856",
	"Writters":      "1286864227475447808",
}

func StartAll(slackWebHookUrlHN, slackWebHookUrlTwitter, slackWebHookUrlCartoons, slackWebHookUrlHNClassics, logUrl string) {
	if !IsTestMode { // if this is not in test mode
		var i int = 0
		for {
			SC.SendPlainText(fmt.Sprint(time.Now().Format("2006-01-02 15:04:05"), ":", "Auto retrieving new posts... "), logUrl)
			go mBot.AutoRetrieveHN(slackWebHookUrlHN)
			// go rc.AutoRetrieveNew()
			go mBot.AutoRetrieveTwitter(TweetLists, 2000, slackWebHookUrlTwitter)
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

func TestAll(t *testing.T) {
	mBot.AutoRetrieveHN("https://hooks.slack.com/services/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")
}
