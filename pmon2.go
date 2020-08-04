package main

import (
	"flag"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/client/actions/del"
	"github.com/ntt360/pmon2/client/actions/desc"
	"github.com/ntt360/pmon2/client/actions/help"
	"github.com/ntt360/pmon2/client/actions/list"
	"github.com/ntt360/pmon2/client/actions/reload"
	"github.com/ntt360/pmon2/client/actions/start"
	"github.com/ntt360/pmon2/client/actions/stop"
	"os"
)

func main() {
	flag.Parse()
	app.Instance("/etc/pmon2/config/config.yml")

	args := flag.Args()
	argsLen := len(args)
	if argsLen == 0 {
		// TODO print help
		os.Exit(0)
	}

	firstParam := args[0]
	leftParams := args[1:]

	switch firstParam {
	case "start":
		start.Run(leftParams)
		break
	case "list", "ls":
		list.Run(leftParams)
		break
	case "desc": // show process detail info
		desc.Run(leftParams)
		break
	case "stop":
		stop.Run(leftParams)
		break
	case "del", "delete":
		del.Run(leftParams)
		break
	case "reload":
		reload.Run(leftParams)
		break
	default:
		help.Run(leftParams)
	}

}
