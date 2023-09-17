package mediabot

import (
	"fmt"
	"log"
	"strconv"

	"gorm.io/gorm"
)

// const (
// 	sqliteFile string = "file:./data/ids.db"
// )

type Database struct {
	gormDB *gorm.DB
}

type SavedItem struct { // for saving into gorm
	Id       string
	Platform string
}

func (db *Database) Init(dbDialector gorm.Dialector, dbConfig *gorm.Config) {
	// Init("file:./data/ids.db")
	// var err error
	// if db.gormDB, err = gorm.Open(sqlite.Open(sqliteFile), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Silent),
	// }); err != nil {
	// 	log.Panicln(err)
	// }
	// db.CreateTable()

	// var dbDialector gorm.Dialector = postgres.New(postgres.Config{
	// 	DSN:                  "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
	// 	PreferSimpleProtocol: true, // disables implicit prepared statement usage
	// })

	// var dbConfig *gorm.Config = &gorm.Config{}

	var err error
	db.gormDB, err = gorm.Open(dbDialector, dbConfig)
	if err != nil {
		log.Panicln(err)
	}
}

func (db Database) CreateTable() {
	// db.gormDB.AutoMigrate(&SavedItem{})
}

func (db Database) InsertRow(item SavedItem) {
	// item := SavedItem{Id: newId, Platform: "HackerNews"}
	var result *gorm.DB = db.gormDB.Create(&item)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
}

func (db Database) InsertRows(items []SavedItem) {
	// item := SavedItem{Id: newId, Platform: "HackerNews"}
	var result *gorm.DB = db.gormDB.Create(&items)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
}

func (db Database) QueryRow(id string) (item SavedItem) {
	item = SavedItem{}
	var result *gorm.DB = db.gormDB.First(&item, "Id = ?", id)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
	return
}

func (db Database) ReturnAllRecords(platform string) (savedItems []SavedItem) {
	savedItems = []SavedItem{}
	var result *gorm.DB
	if len(platform) == 0 { // to return all records across platform, pass in platform=""
		result = db.gormDB.Find(&savedItems)
	} else {
		result = db.gormDB.Where("platform = ?", platform).Find(&savedItems)
	}
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
	_ = result
	return
}

func (db Database) DeleteItem(item *SavedItem) (result *gorm.DB) {
	result = db.gormDB.Delete(item)
	return
}

func (db Database) UpdateXkcd() (item SavedItem) {
	item = SavedItem{}
	var result *gorm.DB = db.gormDB.First(&item, "Platform = ?", "xkcd")
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
	result = db.DeleteItem(&item)
	_ = result

	if item == (SavedItem{}) { // if there's no record in the db
		item = SavedItem{Id: "10", Platform: "xkcd"} // create a new record starting from 10
	} else {
		var id int
		id, _ = strconv.Atoi(item.Id)
		item.Id = fmt.Sprint(id + 1)
	}
	db.InsertRow(item)
	return
}

func (db Database) UpdateRow(targetId, newPlatform string) (item SavedItem) {
	item = SavedItem{}
	var result *gorm.DB = db.gormDB.First(&item, "Id = ?", targetId)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
	item.Platform = newPlatform
	db.gormDB.Save(&item)
	return
}
