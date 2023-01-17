package app

import (
	"os"
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/ntt360/pmon2/app/boot"
	"github.com/ntt360/pmon2/app/conf"
	"github.com/ntt360/pmon2/app/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var Log *logrus.Logger
var Config *conf.Tpl
var dbOnce sync.Once
var db *gorm.DB

func init() {
	Log = logrus.New()
	if os.Getenv("PMON2_DEBUG") == "true" {
		Log.SetLevel(logrus.DebugLevel)
		Log.SetReportCaller(true)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
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

		initDb, err := gorm.Open(sqlite.Open(appDbDir+"/data.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db = initDb

		// init table
		if !db.Migrator().HasTable(&model.Process{}) {
			db.Migrator().CreateTable(&model.Process{})
		}

		if !db.Migrator().HasTable(&model.App{}) {
			db.Migrator().CreateTable(&model.App{})
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
