package god

import (
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"os"
)

var MonQueue chan *model.Process

type Monitor struct {
}

func NewMonitor() {
	MonQueue = make(chan *model.Process)

	go runMonitor()
}

func runMonitor() {
	var all []model.Process
	err := app.Db().Find(&all).Error
	if err != nil {
		return
	}

	for _, process := range all {
		if process.Status == model.StatusStopped || process.Status == model.StatusInit{
			continue
		}

		// just check failed process
		go restartProcess(process)
	}
}

func restartProcess(p model.Process) {
	_, err := os.Stat(fmt.Sprintf("/proc/%d/status", p.Pid))
	if err == nil {
		return
	}

	// proc status file not exit
	if os.IsNotExist(err) && (p.Status == model.StatusRunning || p.Status == model.StatusFailed) {
		execStartProc(p)
	}
}

func execStartProc(p model.Process) {
	//args = append([]string{"restart", p.ProcessFile}, args[1:]...)
	//data, err := proxy.RunProcess(args)
}
