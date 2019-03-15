package datamodels

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	LoginName       string
	Password        string
	Name            string
	Phone           string
	Sex             uint      `gorm:"default:1"` //性别	1男,2女
	DepartmentName  string    //所属部门
	Email           string    //邮件
	Session         string    `gorm:"-"` //Session
	AccessTime      time.Time `gorm:"-"` //登陆时间
	LastAccessTime  time.Time `gorm:"-"` //最后访问时间
	RemoteAddr      string    `gorm:"-"` //IP
	MaxLifeTime     int64     `gorm:"-"` //session过期时间
	Online          bool      `gorm:"-"` //是否在线
	Role            []Role    `gorm:"-"`
	RoleId          []uint    `gorm:"-"`
	ProjectWorkPath string
	UserTypeId      uint
	UserType        UserType `gorm:"ForeignKey:UserTypeId"`
}
type UserRole struct {
	gorm.Model
	UserId uint
	RoleId uint
}
type UserType struct {
	gorm.Model
	Name string
}
