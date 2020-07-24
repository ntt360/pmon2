package proxy

import (
	"errors"
	"io/ioutil"
	"os"
)

func RunProcess(prjDir string, args []string) ([]byte, error) {
	workerExec := prjDir + "/bin/worker"

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	worker, err := os.StartProcess(workerExec, args, &os.ProcAttr{
		Files: []*os.File{nil, w, w},
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
