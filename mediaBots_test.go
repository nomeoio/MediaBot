package mediabot

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var mBot MediaBot

var TweetLists = map[string]string{}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error loading .env file: ", err)
	}

}

func StartAll(slackWebHookUrlHN, slackWebHookUrlTwitter, slackWebHookUrlCartoons, slackWebHookUrlHNClassics, logUrl string) {
	if !IsTestMode { // if this is not in test mode
		var i int = 0
		for {
			SC.SendPlainText(fmt.Sprint(time.Now().Format(time.DateTime), ":", "Auto retrieving new posts... "), logUrl)
			// go mBot.AutoRetrieveHN(slackWebHookUrlHN, "")
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

func TestFormatTG(t *testing.T) {
	var i = HNItem{Id: 38533208,
		Deleted:     false,
		Type:        "story",
		By:          "gmays",
		Descendants: 6,
		//  Kids:[38534059 38534007 38534976 38535371 38535017] ,
		Parent: 0,
		Score:  15,
		Time:   1701794417,
		Title:  "Poor Charlie's Almanack: The Essential Wit and Wisdom of Charles T. Munger",
		Url:    "https://www.stripe.press/poor-charlies-almanack",
		Text:   "",
		Dead:   false,
		Poll:   0,
		// Parts:[],
	}
	var text string
	text, _ = formatTelegramMessage(i)
	TG.SendMessage(os.Getenv("NomieTheBotHTTPAPIToken"), text, os.Getenv("TGChannelHN"), "4")
}

// func TestAutoRetrieveHN(t *testing.T) {
// 	mBot.AutoRetrieveHN("", "-1002072292994")
// }

func TestQueryRows(t *testing.T) {
	var items = []SavedNews{}
	DB.QueryRows([]string{
		"38526610", "38526590", "38526584", "38526579", "38526559", "38526478", "38526439", "38526433", "38526401", "38526400", "38526376", "38526374", "38526359", "38526349", "38526343", "38526277", "38526275", "38526263", "38526247", "38526184", "38526175", "38526166", "38526112", "38526092", "38526085", "38526078", "38526075", "38526070", "38526062", "38526055", "38526051", "38526037", "38526036", "38526026", "38526019", "38526004", "38525978", "38525968", "38525952", "38525941", "38525909", "38525904", "38525876", "38525867", "38525858", "38525851", "38525832", "38525827", "38525825", "38525811", "38525810", "38525797", "38525696", "38525679", "38525655", "38525645", "38525629", "38525625", "38525619", "38525594", "38525591", "38525583", "38525527", "38525523", "38525515", "38525497", "38525484", "38525483", "38525454", "38525446", "38525429", "38525416", "38525385", "38525369", "38525329", "38525314", "38525291", "38525265", "38525250", "38525215", "38525157", "38525148", "38525138", "38525135", "38525119", "38525111", "38525098", "38525078", "38525069", "38524992", "38524985", "38524983", "38524921", "38524907", "38524903", "38524828", "38524803", "38524797", "38524760", "38524738", "38524698", "38524682", "38524666", "38524653", "38524627", "38524599", "38524595", "38524550", "38524498", "38524459", "38524458", "38524444", "38524427", "38524424", "38524408", "38524388", "38524372", "38524336", "38524335", "38524262", "38524260", "38524253", "38524243", "38524223", "38524222", "38524212", "38524211", "38524205", "38524190", "38524186", "38524171", "38524153", "38524127", "38524113", "38524107", "38524094", "38524079", "38524067", "38524063", "38524058", "38524056", "38524047", "38524046", "38524043", "38524042", "38524033", "38524028", "38524002", "38524001", "38523999", "38523995", "38523961", "38523927", "38523898", "38523893", "38523892", "38523852", "38523848", "38523805", "38523774", "38523736", "38523706", "38523704", "38523699", "38523692", "38523693", "38523691", "38523679", "38523672", "38523647", "38523633", "38523625", "38523601", "38523523", "38523504", "38523473", "38523463", "38523455", "38523441", "38523440", "38523438", "38523433", "38523422", "38523408", "38523394", "38523371", "38523354", "38523352", "38523349", "38523340", "38523330", "38523324", "38523318", "38523302", "38523296", "38523239",
	}, &items)
	t.Log("len(items):", len(items))
	t.Logf("items[0]:%+v", items[0])
}

func init() {
	var dbDialector gorm.Dialector = postgres.New(postgres.Config{
		DSN:                  os.Getenv("PostgreSQLString"),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	})
	var dbConfig *gorm.Config = &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel:      logger.Silent, // Log level
				// IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				// ParameterizedQueries:      true,          // Don't include params in the SQL log
				Colorful: true, // Disable color
			},
		)}

	DB.Init("", dbDialector, dbConfig) // creates db file or table if doesn't exist, doesn't do anything if exists.
}
