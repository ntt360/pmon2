package reload

import (
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/god/proc"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/output"
	"github.com/ntt360/pmon2/app/utils/array"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"strings"
	"syscall"
	"time"
)

var signals = []string{"HUP", "USR1", "USR2"}
var sigFlag string

var Cmd = &cobra.Command{
	Use:   "reload",
	Short: "reload some process",
	Long:  "pmon2 will send -SIGUSR2 signal to the process.",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRun(args)
	},
}

func init() {
	Cmd.Flags().StringVarP(&sigFlag, "sig", "s", "", "--sig [unix signal name, such as: HUP]")
}

func cmdRun(args []string) {
	// check user
	uid := os.Geteuid()
	if uid != 0 {
		app.Log.Fatal("command need root or with sudo permission!")
	}

	// check flag
	if len(sigFlag) > 0 {
		if !array.In(signals, strings.ToUpper(sigFlag)) {
			app.Log.Error("the signal only support: HUP，USR1，USR2")
			return
		}
	}

	nl, err := proc.DialPCNWithEvents([]proc.EventType{proc.ProcEventExit, proc.ProcEventFork})
	if err != nil {
		app.Log.Error(err)
		return
	}

	processVal, err := argsValid(args)
	if err != nil {
		app.Log.Fatal(err.Error())
	}

	var process model.Process
	err = app.Db().First(&process, "id = ? or name = ?", processVal, processVal).Error
	if err != nil {
		app.Log.Fatalf("process not found: %s", processVal)
	}

	// 验证进程id的状态
	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", process.Pid))
	if os.IsNotExist(err) {
		app.Log.Errorf("the process %s pid %d not running \n", process.Name, process.Pid)
		return
	}

	// update db state
	var oldState = process.Status
	process.Status = model.StatusReload
	if err = app.Db().Save(&process).Error; err != nil {
		app.Log.Fatalf("try to set process %s state err: %s", processVal, err.Error())
	}

	_, err = os.Stat(fmt.Sprintf("/proc/%d/status", process.Pid))
	if err == nil { // process exist
		p, _ := os.FindProcess(process.Pid)
		err = p.Signal(getSignal())
		app.Log.Debugf("send signal to process: %s", getSignal().String())
		if err != nil {
			app.Log.Fatalf("try send signal to process err: %s", err)
		}
	}

	// try to get process new pid
	pidChannel := make(chan int, 1)
	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer ctxCancel()
	go func() {
		go func() {
			var rel = make(map[string]uint32)
			for {
				event, err := nl.ReadPCN()
				if err != nil {
					continue
				}
				for _, procEvent := range event {
					switch procEvent.What {
					case proc.ProcEventExit:
						if procEvent.EventData.Tgid() != uint32(process.Pid) {
							continue
						}
						rel["exist"] = procEvent.EventData.Tgid()
					case proc.ProcEventFork:
						f := procEvent.EventData.(proc.Fork)
						if f.ParentTgid == uint32(process.Pid) {
							rel["fork"] = f.ParentPid
							rel["new_pid"] = f.ChildPid
						}
					}
				}

				_, ok := rel["exist"]
				if !ok {
					continue
				}

				_, ok = rel["fork"]
				if !ok {
					continue
				}

				// 平滑重启需要保证子进程重启，且父进程退出
				pidChannel <- int(rel["new_pid"])
				_ = nl.ClosePCN()
				break
			}
		}()

		<-ctx.Done()
		_ = nl.ClosePCN()
		pidChannel <- -1
	}()

	newPid := <-pidChannel
	if newPid > 0 { // 进程重启成功
		app.Log.Debugf("reload success process id: %d", newPid)
		process.Pid = newPid
		process.Status = model.StatusRunning
		if err = app.Db().Save(&process).Error; err != nil {
			app.Log.Fatalf("try to update process %s state err %s", processVal, err.Error())
		}
		output.TableOne(process.RenderTable())
	} else { // 重启失败，db数据还原
		app.Log.Debugf("reload failed, roll process status: %s", oldState)
		process.Status = model.StatusFailed
		if err = app.Db().Save(&process).Error; err != nil {
			app.Log.Fatalf("reloadl failed, try rollback process err: %s", err.Error())
		}
		app.Log.Fatalf("process %s reload failed", processVal)
	}
}

func getSignal() syscall.Signal {
	switch strings.ToUpper(sigFlag) {
	case "HUP":
		return syscall.SIGHUP
	case "USR1":
		return syscall.SIGUSR1
	default:
		return syscall.SIGUSR2
	}
}
