package process

import (
	"encoding/json"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/client/proxy"
	"os"
)

func FindByProcessFile(pFile string) *model.Process{
	var rel model.Process
	err := app.Db().First(&rel, "process_file = ?", pFile).Error
	if err != nil {
		return nil
	}

	return &rel
}

func IsRunning(pid int) bool {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return !os.IsNotExist(err)
	}

	return true
}

func TryStart(m model.Process) ([]string, error) {
	// restart process -- --user --log
	args := []string{"--user", m.Username, "--log", m.Log}
	if len(m.Args) > 0 {
		args = append(args, "--", m.Args)
	}
	args = append([]string{"restart", m.ProcessFile}, args...)
	data, err := proxy.RunProcess(args)
	if err != nil {
		return nil, err
	}

	var tb []string
	_ = json.Unmarshal(data, &tb)

	return tb, nil
}
