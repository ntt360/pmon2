package main

import (
	"flag"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/client/actions/desc"
	"github.com/ntt360/pmon2/client/actions/list"
	"github.com/ntt360/pmon2/client/actions/start"
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
	switch args[0] {
	case "start":
		start.Run(args[1:])
		break
	case "list", "ls":
		list.Run(args[1:])
		break
	case "desc": // show process detail info
		desc.Run(args[1:])
	default:
		// TODO show print help
	}

}
