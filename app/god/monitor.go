package god

import (
	"fmt"
	"github.com/goinbox/shell"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	process2 "github.com/ntt360/pmon2/app/svc/process"
	"github.com/ntt360/pmon2/app/utils/iconv"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Monitor struct {
}

func NewMonitor() {
	runMonitor()
}

func runMonitor() {
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		<-timer.C
		runningTask()
	}
}

var pendingTask sync.Map

func runningTask() {
	var all []model.Process
	err := app.Db().Find(&all, "status = ? or status = ?", model.StatusRunning, model.StatusFailed).Error
	if err != nil {
		return
	}

	for _, process := range all {
		// just check failed process
		key := "process_id:" + strconv.Itoa(int(process.ID))
		_, ok := pendingTask.LoadOrStore(key, process.ID)
		if ok {
			return
		}

		go func(p model.Process, key string) {
			var cur model.Process
			defer func() {
				pendingTask.Delete(key)
			}()
			err = app.Db().First(&cur, p.ID).Error
			if err != nil {
				return
			}

			if cur.Status != model.StatusRunning && cur.Status != model.StatusFailed {
				return
			}

			// 启动大于5秒后的进程才进行检查
			if time.Since(cur.UpdatedAt).Seconds() <= 5 {
				return
			}

			err = restartProcess(p)
			if err != nil {
				app.Log.Error(err)
			}
		}(process, key)
	}
}

// Detects whether a new process is created
func checkFork(process model.Process) bool {
	// try to get process new pid
	rel := shell.RunCmd(fmt.Sprintf("ps -ef | grep '%s ' | grep -v grep | awk '{print $2}'", process.ProcessFile))
	if rel.Ok {
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := iconv.MustInt(newPidStr)
		if newPid != 0 && newPid != process.Pid {
			process.Pid = newPid
			process.Status = model.StatusRunning
			return app.Db().Save(&process).Error == nil
		}
	}

	return false
}

func restartProcess(p model.Process) error {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err == nil { // process already running
		//fmt.Printf("Monitor: process (%d) already running \n", p.Pid)
		return nil
	}

	// proc status file not exit
	if os.IsNotExist(err) && (p.Status == model.StatusRunning || p.Status == model.StatusFailed) {
		if checkFork(p) {
			return nil
		}

		// check whether set auto restart
		if !p.AutoRestart {
			if p.Status == model.StatusRunning { // but process is dead, update db state
				p.Status = model.StatusFailed
				app.Db().Save(&p)
			}
			return nil
		}

		_, err := process2.TryStart(p, "")
		if err != nil {
			return err
		}
	}

	return nil
}
