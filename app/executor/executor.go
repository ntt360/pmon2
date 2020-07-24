package executor

import (
	"github.com/ntt360/pmon2/app/crypto"
	"github.com/ntt360/pmon2/app/model"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func Exec(processFile, customLogFile, name, extArgs string) (*model.Process, error) {
	logPath, err := getLogPath(customLogFile, crypto.Crc32Hash(processFile))
	if err != nil {
		return nil, err
	}
	logOutput, err := getLogFile(logPath)
	if err != nil {
		return nil, err
	}

	attr := &os.ProcAttr{
		Env:   os.Environ(),
		Files: []*os.File{nil, logOutput, logOutput},
	}

	var processParams = []string{processFile}
	if len(extArgs) > 0 {
		processParams = append(processParams, strings.Split(extArgs, " ")...)
	}

	process, err := os.StartProcess(processFile, processParams, attr)
	if err != nil {
		return nil, err
	}

	// get process file name
	if len(name) <= 0 {
		name = filepath.Base(processFile)
	}

	pModel := model.Process{
		Pid:         process.Pid,
		Log:         logPath,
		Name:        name,
		ProcessFile: processFile,
		Args:        strings.Join(processParams[1:], " "),
		Pointer:     process,
		Status:      model.StatusInit,
	}

	return &pModel, nil
}

func getLogPath(customLogFile string, hash string) (string, error) {
	prjDir := os.Getenv("HOME") + "/.pmon"
	if len(customLogFile) <= 0 {
		defLogDir := prjDir + "/logs/"
		_, err := os.Stat(defLogDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(defLogDir, 0755)
			if err != nil {
				return "", err
			}
		}
		customLogFile = defLogDir + hash
	}

	return customLogFile, nil
}

func getLogFile(customLogFile string) (*os.File, error) {
	// 创建进程日志文件
	logFile, err := os.OpenFile(customLogFile, syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}
