package reload

import (
	"context"
	"fmt"
	"github.com/goinbox/shell"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
	"github.com/ntt360/pmon2/app/utils/iconv"
	"os"
	"strings"
	"syscall"
	"time"
)

func Run(args []string) {
	processVal, err := argsValid(args)
	if err != nil {
		app.Log.Fatal(err.Error())
	}

	var process model.Process
	err = app.Db().First(&process, "id = ? or name = ?", processVal, processVal).Error
	if err != nil {
		app.Log.Fatal("process not found: %s", processVal)
	}

	// update db state
	var oldState = process.Status
	process.Status = model.StatusReload
	if app.Db().Save(&process).Error != nil {
		app.Log.Fatal(err)
	}

	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", process.Pid))
	if err == nil { // process exist
		// kill process
		p, _ := os.FindProcess(process.Pid)
		err := p.Signal(syscall.SIGUSR2)
		if err != nil {
			app.Log.Fatal(err)
		}
	}

	// try to get process new pid
	pidChannel := make(chan int, 1)
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer ctxCancel()
	go func(ctx context.Context) {
		timer := time.NewTicker(time.Millisecond * 300)
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				rel := shell.RunCmd(fmt.Sprintf("ps -ef | grep %s | grep -v grep | awk '{print $2}'", process.ProcessFile))
				if rel.Ok {
					newPidStr := strings.TrimSpace(string(rel.Output))
					newPid := iconv.MustInt(newPidStr)
					if newPid != 0 && newPid != process.Pid {
						pidChannel <- newPid
						return
					}
				}
			case <-ctx.Done():
				pidChannel <- -1
			}
		}
	}(ctx)

	newPid := <-pidChannel
	if newPid > 0 { // 进程重启成功
		process.Status = model.StatusRunning
		if err = app.Db().Save(&process).Error; err != nil {
			app.Log.Fatal(err)
		}
		output.Table([][]string{process.RenderTable()})
	} else {
		process.Status = oldState
		if err = app.Db().Save(&process).Error; err != nil {
			app.Log.Fatal(err)
		}

		app.Log.Fatal(fmt.Printf("process %s reload failed", processVal))
	}
}
