package app

import (
	"github.com/jinzhu/gorm"
	"github.com/ntt360/pmon2/app/boot"
	"github.com/ntt360/pmon2/app/conf"
	"github.com/ntt360/pmon2/app/model"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var Log *logrus.Logger
var Config *conf.Tpl
var dbOnce sync.Once
var db *gorm.DB

func init() {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	Log.SetReportCaller(true)
}

func Instance(confDir string) error {
	tpl, err := boot.Conf(confDir)
	if err != nil {
		return err
	}

	Config = tpl

	return nil
}

func Db() *gorm.DB {
	dbOnce.Do(func() {
		appDbDir := Config.Data + "/db"
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

		if !db.HasTable(&model.App{}) {
			db.CreateTable(&model.App{})
		}

		// sync data
		var appModel model.App
		err = db.First(&appModel).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound { // first version
				db.Create(&model.App{Version: conf.Version})
			}
		}
	})

	return db
}
