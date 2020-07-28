package main

import (
	"flag"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/server"
)

func main() {
	flag.Parse()
	app.Instance("/etc/pmon2/config/config.yml")

	// start monitor process file
	god.NewMonitor()

	// init socket server
	s := server.New(app.Config.GetSockFile())
	s.Run()
}
