package app

import (
	"github.com/ntt360/pmon2/app/boot"
	"github.com/ntt360/pmon2/app/conf"
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger
var Config *conf.Tpl

func init() {
	Log = logrus.New()
	Log.SetLevel(logrus.DebugLevel)
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
	//Log.SetReportCaller(true)
}

func Instance(confDir string) {
	tpl, err := boot.Conf(confDir)
	if err != nil {
		Log.Fatal(err)
	}

	Config = tpl
}
