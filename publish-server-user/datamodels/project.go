package datamodels

import (
	"github.com/jinzhu/gorm"
)

type Project struct {
	gorm.Model
	Name                string      //名称
	ProjectTypeId       uint        //所属类型
	ProjectType         ProjectType `gorm:"ForeignKey:ProjectTypeId"`
	Remarks             string      `gorm:"type:varchar(100)"` //备注
	GitAddress          string      `gorm:"type:varchar(500)"` //git地址
	ServerAddress       string      `gorm:"type:varchar(500)"` //服务器项目地址
	OnlineAccessAddress string      `gorm:"type:varchar(500)"` //项目线上访问地址
	State               uint        //状态
	Port                string      `gorm:"type:varchar(20)"`  //端口号
	WarehouseName       string      `gorm:"type:varchar(100)"` //仓库名
	Setting             []Setting   `gorm:"-"`
	User                User        `gorm:"-"`
	ConfAddr            string      `gorm:"type:varchar(500)"` //conf配置文件地址 逗号隔开
	Enable              uint        `default:2`                //1.正在启动,2.未启动
}

type UserProject struct {
	gorm.Model
	UserId    uint
	ProjectId uint
}
type ProjectType struct {
	gorm.Model
	Name string
}
