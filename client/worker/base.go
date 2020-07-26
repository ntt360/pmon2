package worker

import (
	"github.com/ntt360/pmon2/app/utils"
	"os"
	"os/user"
	"strconv"
)

func GetProcUser(a *utils.Args) (*user.User, error) {
	runUser := a.Get("user")
	var curUser *user.User
	var err error

	if len(runUser) <= 0 {
		curUser, err = user.LookupId(strconv.Itoa(os.Getuid()))
	} else {
		curUser, err = user.Lookup(runUser)
	}

	if err != nil {
		return nil, err
	}

	return curUser, nil
}
