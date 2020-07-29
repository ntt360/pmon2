package stop

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"os"
	"os/exec"
	"strconv"
)

func Run(args []string) {
	val := args[0]
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
			if err:= app.Db().Save(&process).Error; err != nil {
				app.Log.Fatalf("stop process %s err \n", val)
			}

			app.Log.Infof("stop process %s success \n", val)
			return
		}
	}

	// try to kill the process
	cmd := exec.Command("kill", strconv.Itoa(process.Pid))
	err = cmd.Run()
	if err != nil {
		app.Log.Fatal(err)
	}

	process.Status = model.StatusStopped
	if app.Db().Save(&process).Error != nil {
		app.Log.Fatalf("stop the process %s failed", val)
	}

	app.Log.Infof("stop process %s success \n", val)
}
