package datasource

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"publish_server_user/datamodels"
	"publish_server_user/tool"
)

var MysqlBD *gorm.DB

func init() {
	var sql = tool.ConfJson.MysqlName + ":" + tool.ConfJson.MysqlPassword + "@(" + tool.ConfJson.MysqlHost + ":" + tool.ConfJson.MysqlPort + ")/" + tool.ConfJson.MysqlDatabase + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", sql)
	if err != nil {
		fmt.Println(err)
	}
	db.DB().SetMaxIdleConns(0)    //最大打开的连接数
	db.DB().SetMaxOpenConns(1000) //设置最大闲置个数
	db.SingularTable(true)
	// 启用Logger，显示详细日志
	db.LogMode(true)
	MysqlBD = db
}
func GetDB() *gorm.DB {
	return MysqlBD
}

func CreateTable() {
	GetDB().AutoMigrate(
		&datamodels.User{},
		&datamodels.Role{},
		&datamodels.Menu{},
		&datamodels.UserRole{},
		&datamodels.RoleMenu{},
		&datamodels.UserProject{},
		&datamodels.Project{},
		&datamodels.ProjectType{},
		&datamodels.Setting{},
		&datamodels.UserType{},
		&datamodels.SystemSetting{},
	)
}
