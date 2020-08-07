package god

import (
	"fmt"
	"github.com/goinbox/shell"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/utils/iconv"
	"github.com/ntt360/pmon2/client/proxy"
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
func runningTask()  {
	var all []model.Process
	err := app.Db().Find(&all).Error
	if err != nil {
		return
	}

	for _, process := range all {
		if process.Status == model.StatusStopped || process.Status == model.StatusReload || process.Status == model.StatusInit{
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
	defer func() {
		fmt.Printf("check process fork state %v", rel.Ok)
	}()

	if rel.Ok {
		newPidStr := strings.TrimSpace(string(rel.Output))
		newPid := iconv.MustInt(newPidStr)
		if newPid != 0 && newPid != process.Pid {
			process.Pid = newPid
			process.Status = model.StatusRunning
			if app.Db().Save(&process).Error != nil {
				return false
			}
			fmt.Printf("process fork new pid %d", &process.Pid)
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

		return execStartProc(p)
	}

	return nil
}

func execStartProc(p model.Process) error {
	// restart process -- --user --log
	args := []string{"--user", p.Username, "--log", p.Log}
	if len(p.Args) > 0 {
		args = append(args, "--", p.Args)
	}
	args = append([]string{"restart", p.ProcessFile}, args...)
	_, err := proxy.RunProcess(args)
	if err != nil {
		return err
	}
	return nil
}
