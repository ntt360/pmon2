package model

import "time"

type Process struct {
	Id        int
	Log       string
	Name      string
	StartCmd  string
	CreatedAt time.Time
}
