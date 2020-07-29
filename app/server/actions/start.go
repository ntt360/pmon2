package actions

import (
	"encoding/json"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/god"
	"github.com/ntt360/pmon2/app/model"
	"strconv"
)

type Start struct {
	p model.Package
}

func NewStart(p model.Package) *Start {
	return &Start{
		p: p,
	}
}

func (s *Start) Rsp(suffix string) ([]byte, error) {
	var m model.Process
	_ = json.Unmarshal([]byte(s.p.Data), &m)

	// save to db
	var originOne model.Process
	err := app.Db().First(&originOne, "process_file = ?", m.ProcessFile).Error
	if err == nil && originOne.ID > 0 { // process already exist
		m.ID = originOne.ID
	}

	err = app.Db().Save(&m).Error
	if err != nil {
		return nil, fmt.Errorf("pmon2 run err: %v", err)
	}

	// push to process watch
	god.MonQueue <- &m

	return renderOutput(&m).ToJson(suffix), nil
}

func renderOutput(m *model.Process) model.Rsp {
	rel := []string{
		strconv.Itoa(int(m.ID)),
		m.Name,
		strconv.Itoa(m.Pid),
		m.Status,
		m.Username,
		m.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	jsonStr, _ := json.Marshal(rel)
	rsp := model.Rsp{
		Code: 0,
		Msg:  "success",
		Data: string(jsonStr),
	}

	return rsp
}
