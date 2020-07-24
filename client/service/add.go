package service

import (
	"encoding/json"
	"errors"
	"github.com/ntt360/pmon2/app/model"
	"github.com/ntt360/pmon2/client/sock"
)

func AddData(process *model.Process) (string, error) {
	conn, err := sock.GetConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	p := model.Package{
		Cmd:  model.CmdStart,
		Data: process.MustJson(),
	}

	msg := p.MustToJson()
	_, err = conn.Write(msg)
	if err != nil {
		return "", err
	}

	rspData, err := validRspData(sock.ReadData(conn))
	if err != nil {
		return "", err
	}

	return rspData.Data, nil
}

func validRspData(relData string) (*model.Rsp, error) {
	var rel model.Rsp
	err := json.Unmarshal([]byte(relData), &rel)
	if err != nil {
		return nil, err
	}

	if rel.Code == 0 {
		return &rel, nil
	}

	return nil, errors.New(rel.Msg)
}
