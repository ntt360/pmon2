package start

import (
	"github.com/ntt360/pmon2/app"
	"log"
	"os"
)

func Run(args []string) {
	_, err := os.Stat(app.Config.Sock)
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
		rel, err = restart(m, args)
	} else {
		rel, err = loadFirst(execPath, args)
	}

	if err != nil {
		app.Log.Fatal(err)
	}

	log.Println(rel)

	//output.Table([][]string{rel})
}
