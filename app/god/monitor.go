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
	go runMonitor()
}

func runMonitor() {
	timer := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-timer.C:
			runningTask()
		}
	}
}

var pendingTask sync.Map

func runningTask() {
	var all []model.Process
	err := app.Db().Find(&all).Error
	if err != nil {
		return
	}

	for _, process := range all {
		if process.Status == model.StatusStopped || process.Status == model.StatusReload || process.Status == model.StatusInit {
			// state no need restart
			continue
		}

		key := "process_id:" + strconv.Itoa(int(process.ID))
		_, ok := pendingTask.LoadOrStore(key, process.ID)
		if ok {
			// start process already running
			continue
		}

		// just check failed process
		go func(p model.Process) {
			_ = restartProcess(p)
			pendingTask.Delete(key)
		}(process)
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
			if app.Db().Save(&process).Error != nil {
				return false
			}

			return true
		}
	}

	return false
}

func restartProcess(p model.Process) error {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err == nil { // process already running
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

		_, err := process2.TryStart(p)
		if err != nil {
			return err
		}
	}

	return nil
}
