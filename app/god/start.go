package god

import (
	"errors"
	"fmt"
	"github.com/ntt360/pmon2/app/crypto"
	"github.com/ntt360/pmon2/app/utils"
	"os"
	"strings"
	"syscall"
)

type StartCmd struct {
	args []string
}

func (s *StartCmd) Run() ([]byte, error) {
	processFile := s.args[0]
	file, err := os.Stat(processFile)
	if os.IsNotExist(err) || file.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s not exist", processFile))
	}

	logKey := crypto.Crc32Hash(processFile)
	a := utils.ParseArgs(s.args)

	// --log 获取自定义日志路径
	customLogFile := a.Get("log")
	prjDir := os.Getenv("HOME") + "/.pmon"
	if len(customLogFile) <= 0 {
		defLogDir := prjDir + "/logs/"
		_, err := os.Stat(defLogDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(defLogDir, 0755)
			if err != nil {
				return nil, err
			}
		}
		customLogFile = defLogDir + logKey
	}

	// 创建进程日志文件
	logFile, err := os.OpenFile(customLogFile, syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}

	attr := &os.ProcAttr{
		Env:   os.Environ(),
		Files: []*os.File{nil, logFile, logFile},
	}

	var processParams = []string{processFile}
	extArgs := a.Get("def_params")
	if len(extArgs) > 0 {
		processParams = append(processParams, strings.Split(extArgs, " ")...)
	}

	process, err := os.StartProcess(processFile, processParams, attr)
	if err != nil {
		return nil, err
	}

	_, err = process.Wait()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
