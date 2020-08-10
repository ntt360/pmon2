package main

import (
	"flag"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/client/cmd"
	"log"
)

func main() {
	flag.Parse()
	app.Instance("/etc/pmon2/config/config.yml")

	err := cmd.Exec()
	if err != nil {
		log.Fatal(err)
	}

	//switch firstParam {
	//case "exec":
	//	exec.Run(leftParams)
	//	break
	//case "start":
	//	start.Run(leftParams)
	//	break
	//case "list", "ls":
	//	list.Run(leftParams)
	//	break
	//case "desc", "show": // show process detail info
	//	desc.Run(leftParams)
	//	break
	//case "stop":
	//	stop.Run(leftParams)
	//	break
	//case "del", "delete":
	//	del.Run(leftParams)
	//	break
	//case "reload":
	//	reload.Run(leftParams)
	//	break
	//default:
	//	help.Run(leftParams)
	//}

}
