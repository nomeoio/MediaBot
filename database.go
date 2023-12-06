package mediabot

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"

	"gorm.io/gorm"
)

// const (
// 	sqliteFile string = "file:./data/ids.db"
// )

// func init() {
// 	var CurrentWorkingDir string
// 	var err error
// 	var envAbs string
// 	CurrentWorkingDir, _ = AbsCwd()
// 	if envAbs, err = filepath.Abs(filepath.Join(CurrentWorkingDir, ".env")); err != nil {
// 		log.Fatalln(err)
// 	}

// 	if err = godotenv.Load(envAbs); err != nil {
// 		log.Fatalln("Error loading .env file: ", err)
// 	}

// }

type Database struct {
	gormDB *gorm.DB
}

type SavedNews struct { // for saving into gorm
	Id       string `json:"id"`
	Platform string `json:"platform"`
	Scores   int    `json:"scores"`
}

func (db *Database) Init(slackWebHookUrlHN string, dbDialector gorm.Dialector, dbConfig *gorm.Config) {
	var err error
	db.gormDB, err = gorm.Open(dbDialector, dbConfig)
	if err != nil {
		log.Panicln(err)
	}
}

func (db Database) CreateTable() {
	// db.gormDB.AutoMigrate(&SavedNews{})
}

func (db Database) InsertRow(item SavedNews) {
	// item := SavedNews{Id: newId, Platform: "HackerNews"}
	var result *gorm.DB = db.gormDB.Create(&item)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
}

func (db Database) InsertRows(items []SavedNews) {
	// item := SavedNews{Id: newId, Platform: "HackerNews"}
	var result *gorm.DB = db.gormDB.Create(&items)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
}

func (db Database) QueryRow(id string) (item SavedNews) {
	var result *gorm.DB = db.gormDB.First(&item, "id = ?", id)
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
	return
}

func (db Database) QueryRows(ids []string, items *[]SavedNews) (result *gorm.DB) {
	return db.gormDB.Raw("SELECT * FROM saved_news WHERE id IN (?)", ids).Scan(items)
}

func (db Database) ReturnAllRecords(platform string) (savedItems []SavedNews) {
	savedItems = []SavedNews{}
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

func (db Database) DeleteItem(item *SavedNews) (result *gorm.DB) {
	result = db.gormDB.Delete(item)
	return
}

func (db Database) UpdateXkcd() (item SavedNews) {
	item = SavedNews{}
	var result *gorm.DB = db.gormDB.First(&item, "Platform = ?", "xkcd")
	if result.Error != nil {
		if result.Error.Error() != "record not found" {
			log.Panicln(result.Error)
		}
	}
	result = db.DeleteItem(&item)
	_ = result

	if item == (SavedNews{}) { // if there's no record in the db
		item = SavedNews{Id: "10", Platform: "xkcd"} // create a new record starting from 10
	} else {
		var id int
		id, _ = strconv.Atoi(item.Id)
		item.Id = fmt.Sprint(id + 1)
	}
	db.InsertRow(item)
	return
}

func (db Database) UpdateRow(targetId, newPlatform string) (item SavedNews) {
	item = SavedNews{}
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

func AbsCwd() (cwd string, err error) {
	// os.Getwd() reuturns where you're in terminal window.
	// this func returns the directory of the executable
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		err = errors.New("unable to get the current filename")
		return
	}
	cwd = filepath.Dir(filename)
	return
}
