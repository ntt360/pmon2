package god

import "encoding/json"

// 数据传输结束标记
const EOF = "\r\n\r\n"

type Package struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func (p *Package) MustToJson() []byte {
	suffix := []byte(EOF)

	rel, err := json.Marshal(p)
	if err != nil {
		return suffix
	}

	// 插入数据结束符
	rel = append(rel, suffix...)
	return rel
}
