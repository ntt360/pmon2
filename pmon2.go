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
}
