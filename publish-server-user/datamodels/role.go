package datamodels

import "github.com/jinzhu/gorm"

//角色表
type Role struct {
	gorm.Model
	Name    string //角色名
	Remarks string //备注
}
type RoleMenu struct {
	gorm.Model
	RoleId uint
	MenuId uint
}
