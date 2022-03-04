package mediabot

import (
	"log"
	"os"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	if Hostname, err = os.Hostname(); err != nil {
		log.Fatalln(err)
	}
	IsLocal = strings.Contains(Hostname, "MacBook")   // checking if the app in local
	IsTestMode = strings.Contains(os.Args[0], "test") // checking if it's in test mode
	// db.CreateTable()                                  // creates db file or table if doesn't exist, doesn't do anything if exists.
}
