package main

import (
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/network"
	"net"
	"os"
)

func main() {
	prjDir := os.Getenv("HOME") + "/.pmon/run/"
	_, err := os.Stat(prjDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(prjDir, 0755)
		if err != nil {
			app.Log.Fatal(err)
		}
	}

	god.NewMonitor()

	unixSocketFile := prjDir + "pmon.sock"
	_, err = os.Stat(unixSocketFile)
	if !os.IsNotExist(err) {
		app.Log.WithField("unix_sock", unixSocketFile).Debug("unix socket file already exist, del it")
		err = os.Remove(unixSocketFile)
		if err != nil {
			app.Log.Fatal(err)
		}
	}

	l, err := net.Listen("unix", unixSocketFile)
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
		go network.HandlerConn(conn)
	}
}
