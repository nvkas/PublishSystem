package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"publish_server_user/datamodels"
	"publish_server_user/datasource"
	"testing"
)

//初始化配置文件
func Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

type JsonConf struct {
	MysqlName     string
	MysqlHost     string
	MysqlPort     string
	MysqlPassword string
	MysqlDatabase string
	PublishPort   string
	ConsulAddr    string
	BasePath      string
}

var ConfJson JsonConf
var userService = NewUserService()

func GetInit() {
	Load("conf.json", &ConfJson)
	fmt.Println(ConfJson)
	var sql = ConfJson.MysqlName + ":" + ConfJson.MysqlPassword + "@(" + ConfJson.MysqlHost + ":" + ConfJson.MysqlPort + ")/" + ConfJson.MysqlDatabase + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", sql)
	if err != nil {
		fmt.Println(err)
	}
	db.DB().SetMaxIdleConns(0)    //最大打开的连接数
	db.DB().SetMaxOpenConns(1000) //设置最大闲置个数
	db.SingularTable(true)
	// 启用Logger，显示详细日志
	db.LogMode(true)
	datasource.MysqlBD = db
}

//用户登录方法测试
func TestLogin(t *testing.T) {
	GetInit()
	log.Println("初始化", ConfJson)
	u := datamodels.User{}
	u.LoginName = "admin"
	u.Password = "123456"
	var result datamodels.Result
	err := userService.Login(context.Background(), &u, &result)
	if err != nil {
		log.Println("出现异常：", err)
	} else {
		if result.Status == false {
			log.Println("账号或密码错误")
		} else {
			log.Println("登录成功")
		}
	}
}

//检测用户名是否存在方法测试
func TestFindLoginName(t *testing.T) {
	GetInit()
	log.Println("初始化", ConfJson)
	u := datamodels.User{}
	u.LoginName = "admin"
	var result datamodels.Result
	err := userService.FindLoginName(context.Background(), &u, &result)
	if err != nil {
		log.Println("出现异常：", err)
	} else {
		if result.Status == false {
			log.Println("用户名已存在,不能注册了")
		} else {
			log.Println("用户名不存在,能注册")
		}
	}
}

//根据用户ID查询所对应的菜单权限
func TestFindMenuRole(t *testing.T) {
	GetInit()
	log.Println("初始化", ConfJson)
	u := datamodels.User{}
	u.ID = 1
	var result datamodels.Result
	err := userService.FindMenuRole(context.Background(), &u, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//用户的添加修改操作
func TestAddOrUpdateUser(t *testing.T) {
	GetInit()
	jsonStr := `{
       "user":{
        "DepartmentName": "654321",
        "Email": "654321",
        "LoginName": "123456789",
        "Name": "adinadminss",
        "ProjectWorkPath": "/usr/local/app/projects/zhangzeyi",
        "Password": "123456",
        "Phone": "18273886975",
        "Sex": 1,
		 "userTypeId":1
       },
       "roleId":[1],
       "id": 5
}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.AddOrUpdateUser(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//用户查询操作
func TestFindUserAll(t *testing.T) {
	GetInit()
	jsonStr := `{
              "user":{
                     "Name": "张泽意",
                     "LoginName": "admin"
                  },
              "page":{
                  "CurrentPage":1,
                  "PageSize":10
                  },
              "Id":1
          }`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.ResultPage
	err := userService.FindUserAll(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//用户修改密码
func TestUpdateUserPassword(t *testing.T) {
	GetInit()
	jsonStr := `{
	"user":{
       "Password":"123456"
	},
	"newPassword":"654321",
   "id":1
}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.UpdateUserPassword(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

// 角色的添加和修改以及分配对应的菜单
func TestAddOrUpdateRole(t *testing.T) {
	GetInit()
	jsonStr := `{
	"role":{
       "name":"系统管理员",
		"remarks":"eee"
	},
	"menuId":[1,2,3,4],
	"id":0
}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.AddOrUpdateRole(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//权限查询
func TestFindRoleAll(t *testing.T) {
	GetInit()
	jsonStr := `{
					"role":{
   					    "Name": ""
						},
					"page":{
						"CurrentPage":1,
						"PageSize":10
						}
				}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindRoleAll(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

/**
根据角色ID获取角色对应的菜单权限ID
*/
func TestFindUserRoleIdAll(t *testing.T) {
	GetInit()
	jsonStr := `{
					"id":1
				}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindMenuIdAll(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

/**
项目的添加和修改
*/
func TestAddOrUpdateProject(t *testing.T) {
	GetInit()
	jsonStr := `{
	"id":9,
	"Name":"项目名",
	"ProjectTypeId":1,
	"ServerAddress":"/usr/local/app/projects/sheke",
	"GitAddress":"http://server.spacej.tech:11300/ZyAng/spacej.git",
	"Remarks":"社科项目",
	"OnlineAccessAddress":"www.spacej.tech/sheke",
	"State":1,
	"Port":"8080",
	"WarehouseName":"ss",
	"UserId":5,
	"ConfAddr":"conf.json"
}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.AddOrUpdateProject(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

/**
项目环境的配置的添加和修改
*/
func TestAddOrUpdateSetting(t *testing.T) {
	GetInit()
	jsonStr := `[{
	"id":1,
	"ProjectId":1,
	"Keys":"openUsername",
	"Values":"zzysss",
	"Enable":"true"
},
{
	"id":2,
	"ProjectId":1,
	"Keys":"openPassword",
	"Values":"123456",
	"Enable":"true"
},
{
	"id":3,
	"ProjectId":1,
	"Keys":"deployKey",
	"Values":"789456123456789123",
	"Enable":"true"
},
{
	"id":4,
	"ProjectId":1,
	"Keys":"pullPath",
	"Values":"我也不知道是什么路径",
	"Enable":"true"
},
{
	"id":5,
	"ProjectId":1,
	"Keys":"runPath",
	"Values":"运行路径",
	"Enable":"true"
},
{
	"id":6,
	"ProjectId":1,
	"Keys":"beforeCommand",
	"Values":"发布前的命令",
	"Enable":"true"
},
{
	"id":7,
	"ProjectId":1,
	"Keys":"publishCommand",
	"Values":"发布的命令",
	"Enable":"true"
},
{
	"id":8,
	"ProjectId":1,
	"Keys":"afterCommand",
	"Values":"发布后的命令",
	"Enable":"true"
}
]`
	var maps1 []interface{}
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.AddOrUpdateSetting(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

/**
项目环境的配置的查询 根据项目ID
*/
func TestFindProjectSetting(t *testing.T) {
	GetInit()
	jsonStr := `{
	"projectId":1
}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindProjectSetting(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//查询仓库名是否重复
func TestFindWarehouseName(t *testing.T) {
	GetInit()
	jsonStr := `{
		"warehouseName":"spacej"
	}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindWarehouseName(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//查询角色名是否重复
func TestFindRoleName(t *testing.T) {
	GetInit()
	jsonStr := `{
		"roleName":"spacej"
	}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindRoleName(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//项目查询All
func TestFindProjectAll(t *testing.T) {
	GetInit()
	jsonStr := `{
    "project":{
            "Name": "",
			"WarehouseName":""
     },
    "page":{
            "CurrentPage":1,
            "PageSize":10
    },
	"projectTypeId":"1",
	"userId":1
}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.ResultPage
	err := userService.FindProjectAll(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//项目类型查询接口
func TestFindProjectTypeAll(t *testing.T) {
	GetInit()
	jsonStr := `{
    "project":{
            "Name": "",
			"projectTypeId":0,
			"warehouseName":"spacejss"
     },
    "page":{
            "CurrentPage":1,
            "PageSize":10
    }
}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindProjectTypeAll(context.Background(), nil, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//根据用户ID查询
func TestFindUserId(t *testing.T) {
	GetInit()
	jsonStr := `{
    "id":1
}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindUserId(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//用户删除接口
func TestDeleteUserAll(t *testing.T) {
	GetInit()
	jsonStr := `{
    "userId":[16]
	}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.DeleteUserAll(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//项目删除
func TestDeleteProjectAll(t *testing.T) {
	GetInit()
	jsonStr := `{
    "projectId":[15,16]
	}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.DeleteProjectAll(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//角色删除
func TestDeleteRoleAll(t *testing.T) {
	GetInit()
	jsonStr := `{
    "roleId":[18]
	}`
	maps1 := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.DeleteRoleAll(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}
func TestSaveSystemSetting(t *testing.T) {
	GetInit()
	jsonStr := `[{
	"id":0,
	"Keys":"NginxConfPath",
	"Values":"zzysss",
	"Enable":"true"
    },
    {
	"id":0,
	"Keys":"NginxInstallPath",
	"Values":"123456",
	"Enable":"true"
    },
   {
	"id":0,
	"Keys":"GoPath",
	"Values":"789456123456789123",
	"Enable":"true"
   },
  {
	"id":0,
	"Keys":"VuePath",
	"Values":"我也不知道是什么路径",
	"Enable":"true"
}
]`
	var maps1 []interface{}
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.SaveSystemSetting(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

//系统配置查询接口
func TestFindSystemSetting(t *testing.T) {
	GetInit()
	jsonStr := ``
	var maps1 map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindSystemSetting(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}
func TestFindsSystemSetting(t *testing.T) {
	GetInit()
	jsonStr := `{true success 200 map[User:map[MaxLifeTime:0 Role:<nil> ID:5 LoginName:admin123456789 Password:123456 Phone:18273886975 AccessTime:0001-01-01 00:00:00 +0000 UTC RemoteAddr: UserTypeId:1 UserType:map[Name: ID:0 CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC DeletedAt:<nil>] UpdatedAt:2019-02-21 14:49:35 +0800 CST Sex:1 DepartmentName:654321 Session: LastAccessTime:0001-01-01 00:00:00 +0000 UTC ProjectWorkPath:/usr/local/app/projects/test/ Name:张泽意 RoleId:<nil> CreatedAt:2019-02-21 14:49:35 +0800 CST DeletedAt:<nil> Email:1324204490@qq.com Online:false] ID:9 CreatedAt:2019-02-28 16:13:54 +0800 CST DeletedAt:<nil> ProjectType:map[ID:0 CreatedAt:0001-01-01 00:00:00 +0000 UTC UpdatedAt:0001-01-01 00:00:00 +0000 UTC DeletedAt:<nil> Name:] GitAddress:http://server.spacej.tech:11300/shehao/publish_test.git ServerAddress:/usr/local/app/projects/test Remarks: ProjectTypeId:2 OnlineAccessAddress:http://www.hnsdrl.com:8888/login Setting:[map[ProjectId:9 Keys:openUsername Values:shehao Enable:true ID:25 CreatedAt:2019-02-28 17:09:23 +0800 CST UpdatedAt:2019-02-28 17:09:23 +0800 CST DeletedAt:<nil>] map[DeletedAt:<nil> ProjectId:9 Keys:openPassword Values:123456 Enable:true ID:26 CreatedAt:2019-02-28 17:09:23 +0800 CST UpdatedAt:2019-02-28 17:09:23 +0800 CST] map[Enable:true ID:27 CreatedAt:2019-02-28 17:09:23 +0800 CST UpdatedAt:2019-02-28 17:09:23 +0800 CST DeletedAt:<nil> ProjectId:9 Keys:deployKey Values:] map[CreatedAt:2019-02-28 17:09:23 +0800 CST UpdatedAt:2019-02-28 17:09:23 +0800 CST DeletedAt:<nil> ProjectId:9 Keys:pullPath Values:pull Enable:true ID:28] map[Values:run Enable:true ID:29 CreatedAt:2019-02-28 17:09:23 +0800 CST UpdatedAt:2019-02-28 17:09:23 +0800 CST DeletedAt:<nil> ProjectId:9 Keys:runPath] map[Values:go build -o publish_test Enable:true ID:30 CreatedAt:2019-02-28 17:09:23 +0800 CST UpdatedAt:2019-02-28 17:09:23 +0800 CST DeletedAt:<nil> ProjectId:9 Keys:beforeCommand] map[Keys:publishCommand Values:nohup ./publish_test & Enable:true ID:31 CreatedAt:2019-02-28 17:09:23 +0800 CST UpdatedAt:2019-02-28 17:09:23 +0800 CST DeletedAt:<nil> ProjectId:9] map[UpdatedAt:2019-02-28 17:09:23 +0800 CST DeletedAt:<nil> ProjectId:9 Keys:afterCommand Values:cd / Enable:true ID:32 CreatedAt:2019-02-28 17:09:23 +0800 CST]] ConfAddr:config.json,/tool/HttpServer.go UpdatedAt:2019-02-28 16:13:54 +0800 CST Name:publish_test State:0 Port:8888 WarehouseName:shehao/publish_test] }`
	var maps1 map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &maps1)
	fmt.Println(maps1)
	log.Println("初始化", ConfJson)
	var result datamodels.Result
	err := userService.FindSystemSetting(context.Background(), maps1, &result)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(result)
	}
}
