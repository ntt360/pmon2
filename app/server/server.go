package server

import (
	"bytes"
	"encoding/json"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/app/server/actions"
	"net"
	"strings"
)

// socket serve
type SServer struct {
	sock string
}

func New(sockFile string) *SServer {
	s := &SServer{
		sock: sockFile,
	}

	return s
}

func (s *SServer) Run()  {
	l, err := net.Listen("unix", s.sock)
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
		go HandlerConn(conn)
	}
}

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
		if bytes.Contains(buffer, []byte(model.EOF)) {
			break
		}
	}

	var p model.Package
	dataStr := strings.TrimRight(data.String(), model.EOF)
	err := json.Unmarshal([]byte(dataStr), &p)
	if err != nil {
		app.Log.Debug(err)
		errRsp := model.Rsp{
			Code: 1,
			Msg:  err.Error(),
		}
		_, _ = conn.Write(errRsp.ToJson(model.EOF))
		return
	}

	rsp, err := runModule(p)
	if err != nil {
		errRsp := model.Rsp{
			Code: 1,
			Msg:  err.Error(),
		}
		_, _ = conn.Write(errRsp.ToJson(model.EOF))
		return
	}

	_, _ = conn.Write(rsp)
}

func runModule(p model.Package) ([]byte, error) {
	switch p.Cmd {
	case model.CmdStart:
		return actions.NewStart(p).Rsp(model.EOF)
	case model.CmdStop:
		break
	}

	return nil, nil
}
