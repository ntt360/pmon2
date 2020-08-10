package del

import (
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/spf13/cobra"
	"os"
)

var Cmd = &cobra.Command{
	Use:   "del",
	Short: "del process by id or name",
	Run: func(cmd *cobra.Command, args []string) {
		runCmd(args)
	},
}

func runCmd(args []string) {
	if len(args) == 0 {
		app.Log.Fatalf("missing del process id or name \n")
	}

	val := args[0]
	var m model.Process
	err := app.Db().First(&m, "id = ? or name = ?", val, val).Error
	if err != nil {
		app.Log.Fatalf("del process err:%s \n", err.Error())
	}

	if m.Status == model.StatusRunning {
		app.Log.Fatalf("the process %s is running, must stop it firstly \n", val)
	}

	clearData(m)
	app.Log.Info("del process")
}

func clearData(m model.Process) {
	app.Db().Delete(&m)
	_ = os.Remove(m.Log)
}
