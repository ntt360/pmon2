package stop

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
	process2 "github.com/ntt360/pmon2/app/svc/process"
	"github.com/spf13/cobra"
	"os"
)

var Cmd = &cobra.Command{
	Use:     "stop",
	Short:   "stop running process",
	Example: "sudo pmon2 stop [id or name]",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func cmdRun(args []string) {
	var val string
	if len(args) <= 0 {
		app.Log.Fatalf("must input some process id or name")
	}

	if len(args) == 1 {
		val = args[0]
	}

	// stop process force
	forced := false
	if len(args) == 2 {
		val = args[1]
		if args[0] == "-f" {
			forced = true
		}
	}

	var process model.Process
	err := app.Db().Where("id = ? or name = ?", val, val).First(&process).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			app.Log.Fatalf("%s not exist", val)
		}
	}

	// check process is running
	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", process.Pid))
	if os.IsNotExist(err) {
		if process.Status == model.StatusRunning {
			process.Status = model.StatusStopped
			if err := app.Db().Save(&process).Error; err != nil {
				app.Log.Fatalf("stop process %s err \n", val)
			}

			app.Log.Infof("stop process %s success \n", val)
			return
		}
	}

	// try to kill the process
	err = process2.TryStop(forced, &process)
	if err != nil {
		app.Log.Fatalf("stop the process %s failed", val)
	}

	output.TableOne(process.RenderTable())
}
