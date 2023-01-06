package executor

import (
	"github.com/ntt360/errors"
	"github.com/ntt360/pmon2/app"
	"os"
	"strings"
	"syscall"
)

const logSuffix = ".log"

func getLogPath(customLogFile string, hash string, logDir string) (string, error) {
	if len(logDir) <= 0 {
		app.Log.Debugf("custom log dir: %s \n", logDir)
		logDir = app.Config.GetLogsDir()
	}

	prjDir := strings.TrimRight(logDir, "/")
	if len(customLogFile) <= 0 {
		_, err := os.Stat(prjDir)
		if os.IsNotExist(err) {
			err := os.MkdirAll(prjDir, 0755)
			if err != nil {
				return "", errors.Wrapf(err, "err: %s, logs dir: '%s'", err.Error(), prjDir)
			}
		}
		customLogFile = prjDir + "/" + hash + logSuffix
	}

	app.Log.Debugf("log file is: %s \n", customLogFile)

	return customLogFile, nil
}

func getLogFile(customLogFile string) (*os.File, error) {
	// 创建进程日志文件
	logFile, err := os.OpenFile(customLogFile, syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, 0755)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return logFile, nil
}
