package main

import (
	"flag"
	"log"

	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/client/cmd"
)

var mode string
var config string

func main() {
	flag.StringVar(&config, "f", "/etc/pmon2/config/config.yml", "")
	flag.Parse()

	switch mode {
	case "pmond":
		pmond()
	case "worker":
		worker()
	default:
		err := app.Instance(config)
		if err != nil {
			log.Fatal(err)
		}

		err = cmd.Exec()
		if err != nil {
			log.Fatal(err)
		}
	}
}
