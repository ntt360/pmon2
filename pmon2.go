package main

import (
	"flag"
	"github.com/ntt360/pmon2/client"
	"github.com/ntt360/pmon2/client/tasks/start"
)

func main() {
	flag.Parse()
	args := flag.Args()
	argsLen := len(args)
	if argsLen == 0 {
		// TODO print
	}
	switch args[0] {
	case "start":
		start.Start(args[1:])
		break
	case "list":
		client.List(args[1:])
		break
	}

}
