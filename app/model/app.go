package model

import "github.com/jinzhu/gorm"

type App struct {
	gorm.Model
	Version string
}

func (App) TableName() string {
	return "app"
}
