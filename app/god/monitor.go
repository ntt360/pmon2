package god

import (
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/client/proxy"
	"os"
	"strconv"
	"sync"
	"time"
)

var MonQueue chan *model.Process

type Monitor struct {
}

func NewMonitor() {
	MonQueue = make(chan *model.Process)

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
		if process.Status == model.StatusStopped || process.Status == model.StatusInit{
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

func restartProcess(p model.Process) error {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err == nil { // process already running
		return nil
	}

	// proc status file not exit
	if os.IsNotExist(err) && (p.Status == model.StatusRunning || p.Status == model.StatusFailed) {
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
