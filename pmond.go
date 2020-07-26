package main

import (
	"flag"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/server"
	"net"
)

func main() {
	confDir := flag.String("conf", "/etc/pmon2/config/config.yaml", "pmon2 boot config file path")
	flag.Parse()
	app.Instance(*confDir)

	// start monitor process file
	god.NewMonitor()

	l, err := net.Listen("unix", app.Config.GetSockFile())
	if err != nil {
		app.Log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			app.Log.Fatal(err)
		}

		// handler the data
		go server.HandlerConn(conn)
	}
}
