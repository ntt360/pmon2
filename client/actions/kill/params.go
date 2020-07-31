package kill

import (
	"errors"
	"fmt"
	"github.com/ntt360/pmon2/app/utils/array"
	"strings"
)

const (
	SigHup  = "-HUP"
	SigKill = "-9"
	SigQuit = "-QUIT"
	SigTerm = "-TERM"
	SigUsr1 = "-USR1"
	SigUsr2 = "-USR2"
)

var supportSignal = []string{SigHup, SigKill, SigQuit, SigTerm, SigUsr1, SigUsr2}

func argsValid(args []string) (string, string,  error) {
	if len(args) == 0 {
		return "", "", errors.New("please input process id or name")
	}

	// only one params mean signal is term
	if len(args) == 1 {
		return SigTerm, args[0], nil
	}

	curSig := strings.ToUpper(args[0])
	if !array.In(supportSignal, curSig) {
		return "", "", fmt.Errorf("input signal not support: %s", curSig)
	}

	return curSig, args[1], nil
}
