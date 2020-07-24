package model

import (
	"encoding/json"
)

type Rsp struct {
	Code int
	Msg  string
	Data string
}

func (r Rsp) ToJson(suffix string) []byte {
	data, _ := json.Marshal(r)

	return append(data, suffix...)
}
