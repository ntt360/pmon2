package process

import (
	"encoding/json"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/client/proxy"
	"os"
	"os/exec"
	"strconv"
)

func IsRunning(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

func TryStop(forced bool, p *model.Process) error {
	var cmd *exec.Cmd
	if forced {
		cmd = exec.Command("kill", "-9", strconv.Itoa(p.Pid))
	} else {
		cmd = exec.Command("kill", strconv.Itoa(p.Pid))
	}

	err := cmd.Run()
	if err != nil {
		app.Log.Fatal(err)
	}

	p.Status = model.StatusStopped

	return app.Db().Save(p).Error
}

func TryStart(m model.Process, flags string) ([]string, error) {
	var flagsModel = model.ExecFlags{
		User:          m.Username,
		Log:           m.Log,
		NoAutoRestart: !m.AutoRestart,
		Args:          m.Args,
		Name:          m.Name,
	}

	if len(flags) > 0 {
		app.Log.Debugf("start with flags: %s \n", flags)
		execFlags := model.ExecFlags{}
		curFlag, err := execFlags.Parse(flags)
		if err != nil {
			return nil, err
		}

		if len(curFlag.Log) > 0 {
			flagsModel.Log = curFlag.Log
		}

		// log dir
		if len(curFlag.LogDir) > 0 && len(curFlag.Log) == 0 {
			flagsModel.LogDir = curFlag.LogDir
			flagsModel.Log = ""
		}
	}

	data, err := proxy.RunProcess([]string{"restart", m.ProcessFile, flagsModel.Json()})
	if err != nil {
		return nil, err
	}

	var tb []string
	_ = json.Unmarshal(data, &tb)

	return tb, nil
}
