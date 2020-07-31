package kill

import (
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"os"
	"os/exec"
	"strconv"
)

func Run(args []string)  {
	sig, processVal, err := argsValid(args)
	if err != nil {
		app.Log.Fatal(err.Error())
	}

	var process model.Process
	err = app.Db().First(&process, "id = ? or name = ?", processVal, processVal).Error
	if err != nil {
		app.Log.Fatal("process not found: %s", processVal)
	}

	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", process.Pid))
	if err == nil { // process exist
		// kill process
		cmd := exec.Command("kill", sig, strconv.Itoa(process.Pid))
		err := cmd.Run()
		if err != nil {
			app.Log.Fatal(err)
		}
	}

	// update db
	process.Status = model.StatusStopped
	if app.Db().Save(&process).Error != nil {
		app.Log.Fatal("try to update process status err")
	}

	app.Log.Infof("kill process %s success \n", processVal)

}
