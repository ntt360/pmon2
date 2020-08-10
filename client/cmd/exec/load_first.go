package exec

import (
	"encoding/json"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/client/proxy"
	"os"
	"path"
	"path/filepath"
)

func loadFirst(execPath string, args []string) ([]string, error) {
	args = append([]string{"start", execPath }, args[1:]...)
	data, err := proxy.RunProcess(args)
	if err != nil {
		return nil, err
	}

	var tbData []string
	_ = json.Unmarshal(data, &tbData)

	return tbData, nil
}

// check the process already have
func processExist(execPath string) (*model.Process, bool) {
	var process model.Process
	err := app.Db().First(&process, "process_file = ?", execPath).Error
	if err != nil {
		return nil, false
	}

	return &process, true
}

func getExecFile(args []string) (string, error) {
	execFile := args[0]
	_, err := os.Stat(execFile)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%s not exist", execFile)
	}

	if path.IsAbs(execFile) {
		return execFile, nil
	}

	absPath, err := filepath.Abs(execFile)
	if err != nil {
		return "", fmt.Errorf("%s get file path error: %v", err.Error())
	}

	return absPath, nil
}
