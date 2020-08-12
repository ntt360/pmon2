package worker

import (
	"errors"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/executor"
	"github.com/ntt360/pmon2/app/model"
	process2 "github.com/ntt360/pmon2/app/svc/process"
	"github.com/ntt360/pmon2/client/service"
	"time"
)

func Restart(pFile string, flags *model.ExecFlags) (string, error) {
	m := process2.FindByProcessFile(pFile)
	if m == nil {
		return "", errors.New("try to get process data error")
	}

	cstLog := flags.Log
	if len(cstLog) > 0 && cstLog != m.Log {
		m.Log = cstLog
	}

	cstName := flags.Name
	if len(cstName) > 0 && cstName != m.Name {
		m.Name = cstName
	}

	// checkout process name whether exist except itself
	if app.Db().First(&model.Process{}, "name = ? AND id != ?", cstName, m.ID).Error == nil {
		return "", fmt.Errorf("process name: %s already used, please set other name by --name", cstName)
	}

	extArgs := flags.Args
	if len(extArgs) > 0 {
		m.Args = extArgs
	}

	// get run process user
	runUser, err := GetProcUser(flags)
	if err != nil {
		return "", nil
	}

	process, err := executor.Exec(m.ProcessFile, m.Log, m.Name, m.Args, runUser, !flags.NoAutoRestart)
	if err != nil {
		return "", err
	}

	// update process extra data
	process.ID = m.ID
	process.CreatedAt = m.CreatedAt
	process.UpdatedAt = time.Now()

	waitData := service.NewProcStat(process).Wait()

	return service.AddData(waitData)
}
