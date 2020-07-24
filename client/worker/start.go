package worker

import (
	"errors"
	"fmt"
	"github.com/ntt360/pmon2/app/executor"
	"github.com/ntt360/pmon2/app/utils"
	"github.com/ntt360/pmon2/client/service"
	"os"
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

	// start process
	process, err := executor.Exec(processFile, a.Get("log"), a.Get("name"), a.Get("def_params"))
	if err != nil {
		return "", err
	}
	process.CreatedAt = time.Now()

	// waiting process state
	var stat = service.NewProcStat(process).Wait()

	// return process data
	return service.AddData(stat)
}
