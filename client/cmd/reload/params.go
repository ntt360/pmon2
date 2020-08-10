package reload

import (
	"errors"
)

func argsValid(args []string) (string,  error) {
	if len(args) == 0 {
		return  "", errors.New("please input process id or name")
	}

	return args[0], nil
}
