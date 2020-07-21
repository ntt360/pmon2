package main

import (
	"github.com/ntt360/pmon2/app/network"
	"log"
	"net"
	"os"
	"syscall"
)

func main() {
	prjDir := os.Getenv("HOME") + "/.pmon/run/"
	_, err := os.Stat(prjDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(prjDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	unixSocketFile := prjDir + "pmon.sock"
	_, err = os.Stat(unixSocketFile)
	if os.IsExist(err) {
		syscall.Unlink(unixSocketFile)
		err = os.Remove(unixSocketFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	l, err := net.Listen("unix", unixSocketFile)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			// TODO
			log.Fatal(err)
		}

		// handler the data
		go network.HandlerConn(conn)
	}
}
