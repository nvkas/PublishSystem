package datamodels

import "github.com/jinzhu/gorm"

type System struct {
	gorm.Model
	GoPath        string
	NginxJsonAddr string
}
