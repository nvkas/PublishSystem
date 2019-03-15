package datamodels

import "github.com/jinzhu/gorm"

/*
openUsername	公共Git账号
openPassword	公共Git密码
deployKey		部署密钥
pullPath		项目Pull路径
runPath			项目运行路径
beforeCommand	发布前命令
publishCommand	发布命令
afterCommand	发布后命令
 */

type Setting struct {
	gorm.Model
	ProjectId	uint	//所属项目
	Keys	string `gorm:"type:varchar(100);not null";json:"Key"`	//Key
	Values	string `gorm:"type:varchar(500)";json:"Value"`				//Value
	Enable	string 	//是否启用
}

type SystemSetting struct {
	gorm.Model
	Keys   string `gorm:"type:varchar(100)`  //Keys
	Values string `gorm:"type:varchar(500)"` //Values
	Enable string `gorm:"type:varchar(10)"`  //是否启用true为启用false为不启用
}