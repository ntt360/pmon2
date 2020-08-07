package worker

import (
	"errors"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/executor"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/utils"
	"github.com/ntt360/pmon2/client/service"
	"os"
	"path/filepath"
	"time"
)

func Start(args []string) (string, error) {
	// prepare params
	processFile := args[0]
	file, err := os.Stat(processFile)
	if os.IsNotExist(err) || file.IsDir() {
		return "", errors.New(fmt.Sprintf("%s not exist", processFile))
	}
	a := utils.ParseArgs(args)

	// get run process user
	runUser, err := GetProcUser(a)
	if err != nil {
		return "", nil
	}

	name :=  a.Get("name")
	// get process file name
	if len(name) <= 0 {
		name = filepath.Base(processFile)
	}

	// checkout process name whether exist
	if app.Db().First(&model.Process{}, "name = ?", name).Error == nil {
		return "", fmt.Errorf("process name: %s already exist, please set other name by --name", name)
	}

	// start process
	process, err := executor.Exec(processFile, a.Get("log"), name, a.Get("def_params"), runUser)
	if err != nil {
		return "", err
	}
	process.CreatedAt = time.Now()
	process.UpdatedAt = time.Now()

	// waiting process state
	var stat = service.NewProcStat(process).Wait()

	// return process data
	return service.AddData(stat)
}
