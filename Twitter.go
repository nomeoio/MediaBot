package mediabot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/dghubble/oauth1"
	"github.com/naughtymonsta/utilities"
)

type TwitterClient struct{}

func (tc TwitterClient) LookUpTweets(ids []string, twitterBearerToken string) (tweets []Tweet, err error) {
	tweets = make([]Tweet, len(ids))
	var i int
	var id string
	wg := sync.WaitGroup{}
	for i, id = range ids {
		wg.Add(1)
		go func(i int, id string) {
			defer wg.Done()
			var tweet Tweet
			if tweet, err = tc.lookUpTweet(id, twitterBearerToken); err != nil {
				return
			}
			tweets[i] = tweet
		}(i, id)
	}
	wg.Wait()
	return
}

func (tc TwitterClient) lookUpTweet(id, twitterBearerToken string) (tweet Tweet, err error) {
	var url string = fmt.Sprintf(tweetLoopUpEndpoint, id)
	var respJson []byte
	if respJson, err = tc.SendHttpRequest(url, "v1", twitterBearerToken); err != nil {
		return
	}
	_ = json.Unmarshal(respJson, &tweet)
	utilities.WriteFile(respJson, "data-samples/tweet.json")
	return
}

func (tc TwitterClient) LookUpTwitterUsers(ids []string, idType string, twitterBearerToken string) (respJson []byte, err error) {
	var idsString string = ids[0]
	for _, id := range ids[1:] {
		idsString = idsString + "," + id
	}
	if idType == "id" {
		respJson, err = tc.SendHttpRequest(fmt.Sprintf(usersLookupEndpoint, fmt.Sprintf("?ids=%s", idsString)), "v2", twitterBearerToken)
	} else if idType == "username" {
		respJson, err = tc.SendHttpRequest(fmt.Sprintf(usersLookupEndpoint, fmt.Sprintf("/by?usernames=%s", idsString)), "v2", twitterBearerToken)
	} else {
		err = fmt.Errorf("id type: %s is wrong", idType)
		return
	}
	return
}

func (tc TwitterClient) SendHttpRequest(url, version, twitterBearerToken string) (body []byte, err error) {
	if version == "v1" {
		body = tc.oauth1Request(url)
	} else if version == "v2" {
		var headers = [][]string{{"Authorization", fmt.Sprintf("Bearer %s", twitterBearerToken)}}
		if body, err = utilities.HttpRequest("GET", nil, url, headers); err != nil {
			log.Panic(err)
		}
	}
	return
}

func (tc TwitterClient) oauth1Request(url string) (body []byte) {
	config := oauth1.NewConfig(Params.TwitterApiKey, Params.TwitterApiKeySecret)
	token := oauth1.NewToken(Params.TwitterAccessToken, Params.TwitterAccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}
	return
}

func (tc TwitterClient) RetrieveByCommand(cmdTxt, twitterBearerToken string) (mbs utilities.MessageBlocks, err error) { // /twt listname limit
	var leastLikes int
	var listId, numStr string
	var fields []string = strings.Fields(cmdTxt)
	listId = fields[0]
	numStr = fields[1]
	leastLikes, _ = strconv.Atoi(numStr)
	var mbList [][]utilities.MessageBlock
	mbList, err = tc.RetrieveTweets(listId, twitterBearerToken, 1000, leastLikes)
	if err != nil {
		return
	}
	var i int
	var mb []utilities.MessageBlock
	for i, mb = range mbList {
		if i == 0 {
			mbs.Blocks = []utilities.MessageBlock{{Type: "header", Text: &utilities.ElementText{Type: "plain_text", Text: listId}}}
		}
		mbs.Blocks = append(mbs.Blocks, mb...)
		if err != nil {
			return
		}
		if i == 1 {
			break
		}
	}
	return
}

func (tc TwitterClient) RetrieveTweets(listId, twitterBearerToken string, leastRetweetLikes, leastLikes int) (mbList [][]utilities.MessageBlock, err error) {
	// :leastRetweetLikes: the least number of likes the retweet should have
	var listTweets, qualifiedListTweets []Tweet
	if listTweets, err = tc.GetListContent(listId, twitterBearerToken); err != nil {
		return
	}

	// check if tweets are qualified
	var listTweet Tweet
	var qualifiedSavedItems []SavedItem
	for _, listTweet = range listTweets {
		var retweet *Tweet = listTweet.Retweeted_Status
		var quoted *Tweet = listTweet.Quoted_Status
		if listTweet.Favorite_Count >= leastLikes || (retweet != nil && retweet.Favorite_Count >= leastRetweetLikes) || (quoted != nil && quoted.Favorite_Count >= leastRetweetLikes) { // if tweet is qualified
			if DB.QueryRow(listTweet.Id_Str).Platform == "Twitter" { // if exists
				continue
			} else {
				qualifiedSavedItems = append(qualifiedSavedItems, SavedItem{Id: listTweet.Id_Str, Platform: "Twitter"})
				// db.InsertRow(SavedItem{Id: listTweet.Id_Str, Platform: "Twitter"})
				qualifiedListTweets = append(qualifiedListTweets, listTweet)
			}
		}
	}
	if len(qualifiedSavedItems) > 0 {
		DB.InsertRows(qualifiedSavedItems)
	}

	var mbarr []utilities.MessageBlock
	for _, listTweet = range qualifiedListTweets {
		if mbarr, err = tc.FormatTweet(listTweet); err != nil {
			return
		}
		mbList = append(mbList, mbarr)
	}

	return
}

// func (tc TwitterClient) TestFormatTweet(tweet Tweet) (err error) {
// 	var mbarr []MessageBlock
// 	// var mbList [][]MessageBlock
// 	if mbarr, err = tc.FormatTweet(tweet); err != nil {
// 		return
// 	}
// 	var mbs MessageBlocks
// 	mbs.Blocks = append(mbs.Blocks, mbarr...)
// 	if err = sc.SendBlocks(mbs, os.Getenv("SlackWebHookUrlTwitter")); err != nil {
// 		return
// 	}
// 	return
// }

func (tc TwitterClient) replaceTwitterUrls(txt string) string {
	var reg *regexp.Regexp = regexp.MustCompile(`https:\/\/t.co\/([A-Za-z0-9])\w+`) // match links like "https://t.co/se6Ys5aJ4x"
	var res []string = reg.FindAllString(txt, -1)
	var redirects = make(map[string]string)

	var wg sync.WaitGroup
	for _, url := range res {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			finalurl, _, err := utilities.CheckUrl(url)
			if err != nil {
				return
			}
			if strings.Contains(finalurl, "twitter.com") {
				redirects[url] = finalurl
			}
		}(url)
	}
	wg.Wait()

	for url := range redirects {
		txt = strings.ReplaceAll(txt, url, "")
	}
	return txt
}

func (tc TwitterClient) FormatTweet(tweet Tweet) (mbarr []utilities.MessageBlock, err error) {
	mbarr = append(mbarr, utilities.MessageBlock{Type: "divider"})
	var txt string
	var retweet *Tweet
	if tweet.Retweeted_Status == nil {
		retweet = tweet.Quoted_Status
	} else {
		retweet = tweet.Retweeted_Status
	}

	tweet.Full_Text = tc.replaceTwitterUrls(tweet.Full_Text)
	mbarr = append(mbarr, tc.Addthumbnail(tweet.User.Profile_image_url_https, tweet.User.Screen_Name))
	if retweet != nil { // if it's a retweet
		txt = " RT"
		if len(tweet.Full_Text) > 4 && tweet.Full_Text[:4] != "RT @" {
			txt = tweet.Full_Text + txt
		}
		retweet.Full_Text = tc.replaceTwitterUrls(retweet.Full_Text)
		mbarr = append(mbarr, SC.CreateTextBlock(txt, "mrkdwn", ""))
		var tweetMedia []TweetMedia = tweet.Extended_Entities.Media
		var retweetMedia []TweetMedia = retweet.Extended_Entities.Media
		var tweetRetweetMediaDifferent bool = tweetMedia != nil && retweetMedia != nil && tweetMedia[0].Media_Url_Https != retweetMedia[0].Media_Url_Https
		var tweetHasMediaRetweetDont bool = tweetMedia != nil && retweetMedia == nil
		if tweetRetweetMediaDifferent || tweetHasMediaRetweetDont { // avoid duplicated media
			mbarr = append(mbarr, tc.loopMediaList(tweet.Extended_Entities.Media)...)
		}
		if tweet.Favorite_Count > 0 {
			txt = fmt.Sprintf(`[<https://twitter.com/%s/status/%s|tweet>] retweets: *%d*, likes: *%d*`, tweet.User.Screen_Name, tweet.Id_Str, tweet.Retweet_Count, tweet.Favorite_Count)
			mbarr = append(mbarr, SC.CreateTextBlock(txt, "mrkdwn", ""))
		}
		mbarr = append(mbarr, tc.Addthumbnail(retweet.User.Profile_image_url_https, retweet.User.Screen_Name))
		mbarr = append(mbarr, SC.CreateTextBlock(retweet.Full_Text, "mrkdwn", ""))
		mbarr = append(mbarr, tc.loopMediaList(retweet.Extended_Entities.Media)...)
		txt = fmt.Sprintf(`[<https://twitter.com/%s/status/%s|tweet>] retweets: *%d*, likes: *%d*`, tweet.User.Screen_Name, tweet.Id_Str, retweet.Retweet_Count, retweet.Favorite_Count)
		mbarr = append(mbarr, SC.CreateTextBlock(txt, "mrkdwn", ""))
	} else {
		mbarr = append(mbarr, SC.CreateTextBlock(tweet.Full_Text, "mrkdwn", ""))
		mbarr = append(mbarr, tc.loopMediaList(tweet.Extended_Entities.Media)...)
		txt = fmt.Sprintf(`[<https://twitter.com/%s/status/%s|tweet>] retweets: *%d*, likes: *%d*`, tweet.User.Screen_Name, tweet.Id_Str, tweet.Retweet_Count, tweet.Favorite_Count)
		mbarr = append(mbarr, SC.CreateTextBlock(txt, "mrkdwn", ""))
	}
	return
}

func (tc TwitterClient) Addthumbnail(thumbnailUrl, username string) utilities.MessageBlock {
	return utilities.MessageBlock{
		Type: "context",
		Elements: []*utilities.Element{{
			Type:      "image",
			Image_Url: thumbnailUrl,
			AltText:   "profile",
		}, {
			Type: "mrkdwn",
			Text: fmt.Sprintf(`<https://twitter.com/%s|@%s>`, username, username),
		}},
	}
}

func (tc TwitterClient) loopMediaList(mediaList []TweetMedia) (mbarr []utilities.MessageBlock) {
	var media TweetMedia
	for _, media = range mediaList {
		var mb utilities.MessageBlock = SC.CreateImageBlock(media.Media_Url_Https, "ok")
		mbarr = append(mbarr, mb)
		mbarr = append(mbarr, utilities.MessageBlock{
			Type: "context",
			Elements: []*utilities.Element{{
				Type:  "plain_text",
				Text:  media.Type,
				Emoji: true,
			}},
		})
	}
	return
}

func (tc TwitterClient) GetListContent(listId, twitterBearerToken string) (tweets []Tweet, err error) {
	var url string = fmt.Sprintf(listEndpoint, listId)
	var respJson []byte
	if respJson, err = tc.SendHttpRequest(url, "v1", twitterBearerToken); err != nil {
		return
	}
	_ = json.Unmarshal(respJson, &tweets)
	return
}

func (tc TwitterClient) GetThread(tweetID, userID, twitterBearerToken, slackUrl string) (err error) {
	var tweets []Tweet
	if tweets, err = tc.getThreadTweets(tweetID, userID, twitterBearerToken); err != nil {
		return
	}
	if tc.sendThread(tweets, slackUrl); err != nil {
		return
	}
	return
}

func (tc TwitterClient) sendThread(threadList []Tweet, slackUrl string) (err error) {
	var mbs utilities.MessageBlocks
	for _, tweet := range threadList {
		var mbarr []utilities.MessageBlock
		if mbarr, err = tc.FormatTweet(tweet); err != nil {
			return
		}
		mbs.Blocks = append(mbs.Blocks, mbarr...)
	}
	if err = SC.SendBlocks(mbs, slackUrl); err != nil {
		return
	}
	return
}

func (tc TwitterClient) getThreadTweets(convoID, userID, twitterBearerToken string) (tweets []Tweet, err error) {
	var url string = fmt.Sprintf(convoEndpoint, convoID, userID, userID)
	url = strings.ReplaceAll(url, " ", "%20")

	var respJson []byte
	var thread Thread
	if respJson, err = tc.SendHttpRequest(url, "v2", twitterBearerToken); err != nil {
		return
	}
	if err = json.Unmarshal(respJson, &thread); err != nil {
		return
	}
	if tweets, err = tc.sortThreadTweets(thread, twitterBearerToken); err != nil {
		return
	}
	return
}

func (tc TwitterClient) sortThreadTweets(thread Thread, twitterBearerToken string) (tweets []Tweet, err error) {
	var sortedThreadTweets []ThreadTweetInfo
	var threadTweet, tweet ThreadTweetInfo
	var threadTweets []ThreadTweetInfo
	for _, tweet = range append(thread.Data, thread.Includes.Tweets...) {
		var exist bool = false
		for _, threadTweet = range threadTweets {
			if tweet.Id == threadTweet.Id {
				exist = true
				break
			}
		}
		if !exist {
			threadTweets = append(threadTweets, tweet)
			if tweet.Referenced_tweets == nil {
				sortedThreadTweets = append(sortedThreadTweets, tweet)
			}
		}
	}
	for i := 0; i < len(threadTweets); i++ {
		var id string = sortedThreadTweets[i].Id
		for _, threadTweet = range threadTweets {
			var referencedTweet ThreadReferencedTweet
			for _, t := range threadTweet.Referenced_tweets {
				if t.Type == "replied_to" {
					referencedTweet = t
					break
				}
			}
			if referencedTweet.Id == id {
				sortedThreadTweets = append(sortedThreadTweets, threadTweet)
				break
			}
		}
	}
	var ids []string
	for _, tt := range sortedThreadTweets {
		ids = append(ids, tt.Id)
	}

	if tweets, err = tc.LookUpTweets(ids, twitterBearerToken); err != nil {
		return
	}
	return
}
