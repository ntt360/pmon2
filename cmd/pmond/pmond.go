package main

import (
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/conf"
	"github.com/ntt360/pmon2/app/god"
	"log"
)

func main() {
	config := conf.GetDefaultConf()
	err := app.Instance(config)
	if err != nil {
		log.Fatal(err)
	}

	// start monitor process file
	god.NewMonitor()
}
