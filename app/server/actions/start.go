package actions

import (
	"encoding/json"
	"fmt"
	"github.com/ntt360/pmon2/app"
	"github.com/ntt360/pmon2/app/model"
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

	return renderOutput(&m).ToJson(suffix), nil
}

func renderOutput(m *model.Process) model.Rsp {
	rel := m.RenderTable()
	jsonStr, _ := json.Marshal(rel)
	rsp := model.Rsp{
		Code: 0,
		Msg:  "success",
		Data: string(jsonStr),
	}

	return rsp
}
