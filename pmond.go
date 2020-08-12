package main

import (
	"flag"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/server"
	"log"
	"path"
	"path/filepath"
)

func main() {
	var config string
	flag.StringVar(&config, "f", "/etc/pmon2/config/config.yml", "")
	flag.Parse()

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
