package model

import "encoding/json"

type ExecFlags struct {
	User          string `json:"user"`
	Log           string `json:"log"`
	NoAutoRestart bool   `json:"no_auto_restart"`
	Args          string `json:"args"`
	Name          string `json:"name"`
}

func (ExecFlags) Parse(jsonStr string) (*ExecFlags, error) {
	var m ExecFlags
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (e *ExecFlags) Json() string {
	content, _ := json.Marshal(e)

	return string(content)
}
