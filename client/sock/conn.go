package sock

import (
	"bytes"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"net"
	"os"
	"strings"
	"time"
)

func GetConn() (net.Conn, error){
	sockFile := app.Config.Sock
	_, err := os.Stat(sockFile)
	if os.IsNotExist(err) {
		return nil, err
	}

	conn, err := net.Dial("unix", sockFile)
	if err != nil {
		return nil, err
	}

	// timeout 5 seconds
	_ = conn.SetDeadline(time.Now().Add(time.Second * 5))

	return conn, nil
}

func ReadData(conn net.Conn) string {
	var rspData strings.Builder
	for {
		buffer := make([]byte, 4096)
		pos, err := conn.Read(buffer)
		if err != nil {
			break
		}

		rspData.Write(buffer[0:pos])
		if bytes.Contains(buffer, []byte(model.EOF)) { // data end
			break
		}
	}

	return strings.TrimRight(rspData.String(), model.EOF)
}

