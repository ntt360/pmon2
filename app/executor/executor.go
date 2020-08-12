package executor

import (
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/utils/crypto"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func Exec(processFile, customLogFile, name, extArgs string, user *user.User, autoRestart bool) (*model.Process, error) {
	logPath, err := getLogPath(customLogFile, crypto.Crc32Hash(processFile))
	if err != nil {
		return nil, err
	}
	logOutput, err := getLogFile(logPath)
	if err != nil {
		return nil, err
	}

	uid, _ := strconv.Atoi(user.Uid)
	gid, _ := strconv.Atoi(user.Gid)

	attr := &os.ProcAttr{
		Env:   os.Environ(),
		Files: []*os.File{nil, logOutput, logOutput},
		Sys: &syscall.SysProcAttr{
			Credential: &syscall.Credential{
				Uid: uint32(uid),
				Gid: uint32(gid),
			},
			Setsid: true,
		},
	}

	var processParams = []string{processFile}
	if len(extArgs) > 0 {
		processParams = append(processParams, strings.Split(extArgs, " ")...)
	}

	process, err := os.StartProcess(processFile, processParams, attr)
	if err != nil {
		return nil, err
	}

	pModel := model.Process{
		Pid:         process.Pid,
		Log:         logPath,
		Name:        name,
		ProcessFile: processFile,
		Args:        strings.Join(processParams[1:], " "),
		Pointer:     process,
		Status:      model.StatusInit,
		Uid:         user.Uid,
		Gid:         user.Gid,
		Username:    user.Username,
		AutoRestart: autoRestart,
	}

	return &pModel, nil
}
