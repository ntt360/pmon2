package main

import (
	"flag"
	"github.com/ntt360/pmon2/cmd"
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
		cmd.Start(args[1:])
		break
	case "list":
		cmd.List(args[1:])
		break
	}

}
