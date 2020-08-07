package start

import (
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
	"github.com/ntt360/pmon2/app/svc/process"
)

func Run(args []string) {
	if len(args) == 0 {
		app.Log.Fatal("please input start process id or name")
	}

	val := args[0]
	var m model.Process
	if err := app.Db().Debug().First(&m, "id = ? or name = ?", val, val).Error; err != nil {
		app.Log.Fatal(fmt.Sprintf("the process %s not exist", val))
	}

	// checkout process state
	if process.IsRunning(m.Pid) {
		if m.Status != model.StatusRunning {
			m.Status = model.StatusRunning
			app.Db().Save(&m)
		}
		output.TableOne(m.RenderTable())
		return
	}

	rel, err := process.TryStart(m)
	if err != nil {
		app.Log.Fatal(err.Error())
	}

	output.TableOne(rel)
}
