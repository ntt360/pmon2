package model

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

const (
	StatusInit    = "init"    // init process
	StatusRunning = "running" // success running
	StatusStopped = "stopped" // success finished or run stop success
	StatusReload  = "reload"  //
	StatusFailed  = "failed"  // run error
)

type Process struct {
	ID          uint        `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `sql:"index" json:"deleted_at"`
	Pid         int         `gorm:"column:pid" json:"pid"`
	Log         string      `gorm:"column:log" json:"log"`
	Name        string      `json:"name"`
	ProcessFile string      `json:"process_file"`
	Args        string      `json:"args"`
	Status      string      `json:"status"`
	Pointer     *os.Process `gorm:"-" json:"-"`
	Uid         string
	Username    string
	Gid         string
}

func (Process) TableName() string {
	return "process"
}

func (p Process) MustJson() string {
	data, _ := json.Marshal(&p)

	return string(data)
}

func (p Process) RenderTable() []string {
	return []string{
		strconv.Itoa(int(p.ID)),
		p.Name,
		strconv.Itoa(p.Pid),
		p.Status,
		p.Username,
		p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
