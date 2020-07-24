package app

import (
	"github.com/ntt360/pmon2/app/model"
	"os"
	"sync"
)
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var dbOnce sync.Once
var db *gorm.DB

func Db() *gorm.DB {
	dbOnce.Do(func() {
		appDbDir := os.Getenv("HOME") + "/.pmon/db"
		_, err := os.Stat(appDbDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(appDbDir, 0755)
			if err != nil {
				panic(err)
			}
		}

		initDb, err := gorm.Open("sqlite3", appDbDir+"/data.db")
		if err != nil {
			panic(err)
		}
		db = initDb

		// init table
		if !db.HasTable(&model.Process{}) {
			db.CreateTable(&model.Process{})
		}
	})

	return db
}
