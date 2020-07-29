package del

import (
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"os"
)

func Run(args []string) {
	if len(args) == 0 {
		app.Log.Fatalf("missing del process id or name")
	}

	val := args[0]
	var m model.Process
	err := app.Db().First(&m, "id = ? or name = ?", val, val).Error
	if err != nil {
		app.Log.Fatalf("del process err:" + err.Error())
	}

	if m.Status == model.StatusRunning {
		app.Log.Fatalf("the process %s is running, must stop it firstly")
	}

	clearData(m)
	app.Log.Info("del process")
}

func clearData(m model.Process) {
	app.Db().Delete(&m)
	_ = os.Remove(m.Log)
}
