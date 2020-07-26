package main

import (
	"flag"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/server"
)

func main() {
	confDir := flag.String("conf", "/etc/pmon2/config/config.yaml", "pmon2 boot config file path")
	flag.Parse()
	app.Instance(*confDir)

	// start monitor process file
	god.NewMonitor()

	// init socket server
	s := server.New(app.Config.GetSockFile())
	s.Run()
}
