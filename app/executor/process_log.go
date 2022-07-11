package executor

import (
	"github.com/ntt360/errors"
	"github.com/ntt360/pmon2/app"
	"os"
	"strings"
	"syscall"
)

const logSuffix = ".log"

func getLogPath(customLogFile string, hash string) (string, error) {
	prjDir := strings.TrimRight(app.Config.GetLogsDir(), "/")
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
