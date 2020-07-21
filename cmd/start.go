package cmd

import (
	"bytes"
	"fmt"
	"github.com/ntt360/pmon2/app/god"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func Start(args []string) {
	prjDir := os.Getenv("HOME") + "/.pmon"
	sockFile := prjDir + "/run/pmon.sock"
	_, err := os.Stat(sockFile)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	conn, err := net.Dial("unix", sockFile)
	if err != nil {
		log.Fatal(err)
	}
	// timeout 5 seconds
	conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer conn.Close()

	p := god.Package{
		Cmd:  god.CmdStart,
		Args: args,
	}

	msg := p.MustToJson()
	pos, err := conn.Write(msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("本次发送数据：%s, length: %d\n", msg, pos)

	var rspData strings.Builder
	for {
		buffer := make([]byte, 4096)
		pos, err := conn.Read(buffer)
		if err != nil {
			break
		}
		rspData.Write(buffer[0:pos])
		if bytes.Contains(buffer, []byte(god.EOF)) { // data end
			break
		}
	}

	// TODO
	fmt.Println(rspData.String())
}
