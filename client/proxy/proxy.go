package proxy

import (
	"errors"
	"github.com/ntt360/pmon2/app"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

func RunProcess(args []string) ([]byte, error) {
	workerExec := strings.TrimRight(app.Config.Bin, "/") + "/worker"
	_, err := os.Stat(workerExec)
	if err != nil {
		return nil, err
	}

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	worker, err := os.StartProcess(workerExec, args, &os.ProcAttr{
		Dir: "/",
		Env: append([]string{"PMON2_CONF=" + app.Config.Conf}, os.Environ()...),
		Files: []*os.File{nil, w, w},
		Sys: &syscall.SysProcAttr{
			Setsid:     true,
		},
	})

	if err != nil {
		return nil, err
	}

	workerState, err := worker.Wait()
	if err != nil {
		return nil, err
	}

	var data []byte
	if workerState.Exited() {
		_ = w.Close()
		data, err = ioutil.ReadAll(r)
		if err != nil {
			_ = r.Close()
			return nil, err
		}

		// run error msg
		if workerState.ExitCode() >= 2 {
			return nil, errors.New(string(data))
		}
	}

	return data, nil
}
