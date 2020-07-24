package service

import (
	"fmt"
	"github.com/ntt360/pmon2/app/model"
	"os"
	"time"
)

type ProcStat struct {
	Done    chan int
	Process *model.Process
}

func NewProcStat(p *model.Process) *ProcStat {
	stat := &ProcStat{
		Done:    make(chan int),
		Process: p,
	}

	stat.Run()

	return stat
}

func (r *ProcStat) Run() {
	go r.ProcessWait(r.Process)
	go r.ProcessExistCheck(r.Process.Pid)
}

func (r *ProcStat) Wait() *model.Process {
	// waiting result
	status := <-r.Done

	if status > 0 {
		r.Process.Status = model.StatusFailed
	} else {
		r.Process.Status = model.StatusRunning
	}

	return r.Process
}

func (r *ProcStat) ProcessWait(process *model.Process) {
	processState, err := process.Pointer.Wait()
	if err != nil {
		r.Done <- 1
		return
	}

	if processState.Exited() {
		r.Done <- processState.ExitCode()
	}
}

func (r *ProcStat) ProcessExistCheck(pid int) {
	timer := time.NewTicker(time.Millisecond * 200)
	defer timer.Stop()
	for {
		select {
		case existCode := <-r.Done:
			if existCode != 0 { // process exist exception
				r.Done <- 1
				return
			}
		case <-timer.C: // check process status by proc file
			_, err := os.Stat(fmt.Sprintf("/proc/%d/status", pid))
			if !os.IsNotExist(err) {
				r.Done <- 0
				return
			}
		}
	}
}
