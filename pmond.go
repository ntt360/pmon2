package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"

	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/server"
)

func pmond() {
	if !path.IsAbs(config) {
		newConf, err := filepath.Abs(config)
		config = newConf
		if err != nil {
			log.Fatal(fmt.Errorf("try to path the config path err: %s", config))
		}
	}

	err := app.Instance(config)
	if err != nil {
		log.Fatal(err)
	}

	// start monitor process file
	god.NewMonitor()

	// init socket server
	s := server.New(app.Config.GetSockFile())
	s.Run()
}
