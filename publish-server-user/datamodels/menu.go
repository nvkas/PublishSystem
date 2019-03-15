package datamodels

import "github.com/jinzhu/gorm"

//菜单
type Menu struct {
	gorm.Model
	MenuId               string //菜单编号
	MenuName             string //菜单名称
	MenuUrl              string //菜单url地址
	Level                int    //等级
	Superior             string //上级菜单
	Icon                 string //图片
	PermissionIdentifier string //许可标识符
	Menu                 []Menu
}
