package start

import (
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/output"
	"log"
	"os"
)

func Run(args []string) {
	prjDir := os.Getenv("HOME") + "/.pmon"
	sockFile := prjDir + "/run/pmon.sock"
	_, err := os.Stat(sockFile)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	// get exec abs file path
	execPath, err := getExecFile(args)
	if err != nil {
		app.Log.Error(err.Error())
		return
	}

	m, exist := processExist(execPath)
	var rel []string
	if exist {
		rel, err = restart(m, prjDir, args)
	} else {
		rel, err = loadFirst(execPath, prjDir, args)
	}

	if err != nil {
		app.Log.Fatal(err)
	}

	output.Table([][]string{rel})
}
