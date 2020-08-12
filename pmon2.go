package main

import (
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/client/cmd"
	"log"
)

func main() {
	//conf := "/etc/pmon2/config/config.yml"
	conf := "/home/liangqi1/devspace/pmon2/config/config.dev.yml"
	err := app.Instance(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Exec()
	if err != nil {
		log.Fatal(err)
	}
}
