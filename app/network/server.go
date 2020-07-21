package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ntt360/pmon2/app/god"
	"log"
	"net"
	"strings"
)

func HandlerConn(conn net.Conn) {
	var data strings.Builder
	defer conn.Close()

	for {
		buffer := make([]byte, 4096)
		pos, err := conn.Read(buffer)
		if err != nil {
			return
		}

		readData := buffer[0:pos]
		data.Write(readData)
		if bytes.Contains(buffer, []byte(god.EOF)) {
			break
		}
	}

	var p god.Package
	dataStr := strings.TrimRight(data.String(), god.EOF)
	err := json.Unmarshal([]byte(dataStr), &p)
	if err != nil {
		errMsg := fmt.Sprintf("write error %v", err)
		_, err := conn.Write([]byte(errMsg + god.EOF))
		if err != nil {
			log.Printf(errMsg)
		}
	}

	stdout, err := runModule(p)
	if err != nil {
		_, err := conn.Write([]byte(err.Error() + god.EOF))
		if err != nil {
			log.Printf(err.Error())
		}
	}

	stdout = append(stdout, []byte(god.EOF)...)
	_, err = conn.Write(stdout)
	if err != nil {
		log.Printf(err.Error())
	}
}

func runModule(p god.Package) ([]byte, error) {
	switch p.Cmd {
	case god.CmdStart:
		cmd := &god.StartCmd{}
		return cmd.Run()
	case god.CmdStop:
		break
	}

	return nil, nil
}
