package restart

import (
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
	"github.com/ntt360/pmon2/app/svc/process"
	"github.com/spf13/cobra"
	"os"
)

var Cmd = &cobra.Command{
	Use:   "restart",
	Short: "restart some process by id or name",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func cmdRun(args []string) {
	if len(args) == 0 {
		app.Log.Fatal("please input restart process id or name")
	}

	val := args[0]
	var m model.Process
	if err := app.Db().First(&m, "id = ? or name = ?", val, val).Error; err != nil {
		app.Log.Fatal(fmt.Sprintf("the process %s not exist", val))
	}

	// checkout process state
	if process.IsRunning(m.Pid) {
		if err := process.TryStop(false, &m); err != nil {
			app.Log.Fatalf("restart error: %s", err.Error())
		}
	}

	rel, err := process.TryStart(m)
	if err != nil {
		if len(os.Getenv("PMON2_DEBUG")) > 0 {
			app.Log.Fatalf("%+v", err)
		} else {
			app.Log.Fatal(err.Error())
		}
	}

	output.TableOne(rel)
}
