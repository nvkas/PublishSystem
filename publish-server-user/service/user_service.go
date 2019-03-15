package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/nsqio/go-nsq"
	"github.com/prometheus/common/log"
	"github.com/smallnest/rpcx/client"
	"publish_server_user/datamodels"
	"publish_server_user/repository"
	"publish_server_user/tool"
	"strconv"
	"strings"
)

var (
	repoUser = repository.NewUserRepo()
	msgError = "系统异常,操作失败,请联系管理员"
	xClient  client.XClient
	producer *nsq.Producer
)

func init() {
	p, err := nsq.NewProducer(tool.ConfJson.NsqAddr, nsq.NewConfig())
	if err != nil {
		log.Fatal("连接NSQ出现异常")
	}
	producer = p
}

type UserService interface {
	Login(ctx context.Context, user *datamodels.User, result *datamodels.Result) error
	FindUserId(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	FindLoginName(ctx context.Context, user *datamodels.User, result *datamodels.Result) error
	FindUserAll(ctx context.Context, maps map[string]interface{}, resultPage *datamodels.ResultPage) error
	FindMenuRole(ctx context.Context, user *datamodels.User, result *datamodels.Result) error
	AddOrUpdateUser(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	UpdateUserPassword(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	GetMenuAll(ctx context.Context, null, result *datamodels.Result) error
	GetRoleAll(ctx context.Context, null, result *datamodels.Result) error
	AddOrUpdateRole(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	FindRoleAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	FindMenuIdAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	AddOrUpdateProject(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	AddOrUpdateSetting(ctx context.Context, maps []interface{}, result *datamodels.Result) error
	FindProjectSetting(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	FindWarehouseProject(ctx context.Context, maps *map[string]interface{}, result *datamodels.Result) error
	FindWarehouseName(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	FindRoleName(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	FindProjectAll(ctx context.Context, maps map[string]interface{}, result *datamodels.ResultPage) error
	FindProjectTypeAll(ctx context.Context, null, result *datamodels.Result) error
	FindUserTypeAll(ctx context.Context, null, result *datamodels.Result) error
	DeleteUserAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	DeleteProjectAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	DeleteRoleAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	SaveSystemSetting(ctx context.Context, maps []interface{}, result *datamodels.Result) error
	FindSystemSetting(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	ProjectStart(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
	ProjectStop(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error
}

type UserServices struct {
	repo repository.UserRepo
}

func NewUserService() UserService {
	return &UserServices{}
}

/**
用户修改密码
*/
func (u *UserServices) UpdateUserPassword(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	fmt.Println(maps)
	var user datamodels.User
	newPassword := ""
	var Id uint
	err := mapstructure.Decode(maps["user"], &user)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	usedPassword := user.Password
	err = mapstructure.Decode(maps["newPassword"], &newPassword)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["id"], &Id)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	user.ID = Id
	//判断原密码是否正确
	repoUser.FindUserId(&user)
	if user.Password != usedPassword {
		datamodels.GetResult(result, nil, "原密码不正确")
		return nil
	} else {
		bool := repoUser.UpdateUserPassword(user.ID, newPassword)
		if bool == false {
			datamodels.GetResult(result, nil, msgError)
			return nil
		}
		datamodels.GetResult(result, true)
		return nil
	}
	return nil
}

/**
根据用户ID查询用户
*/
func (u *UserServices) FindUserId(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	fmt.Println(maps)
	var user datamodels.User
	var id uint
	err := mapstructure.Decode(maps["id"], &id)
	if err != nil {
		log.Error(err)
		return nil
	}
	user.ID = id
	fmt.Println(id)
	repoUser.FindUserId(&user)
	datamodels.GetResult(result, user)
	return nil
}

//查询所有用户信息
func (u *UserServices) FindUserAll(ctx context.Context, maps map[string]interface{}, resultPage *datamodels.ResultPage) error {
	var user datamodels.User
	var page datamodels.Page
	var Id uint
	err := mapstructure.Decode(maps["user"], &user)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["page"], &page)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["Id"], &Id)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	user.ID = Id
	userSlice := repoUser.FindUserAll(user, &page)
	//查询用户对应的角色
	for i, v := range userSlice {
		roleSlice := repoUser.BasisUserIdFindRole(v)
		userSlice[i].Role = roleSlice
	}
	datamodels.GetResultPage(resultPage, userSlice, page.PageCount)
	return nil
}

//用户登录
func (u *UserServices) Login(ctx context.Context, user *datamodels.User, result *datamodels.Result) error {
	repoUser.FindUserNamePassword(user)
	if user.ID > 0 {
		datamodels.GetResult(result, user)
	} else {
		datamodels.GetResult(result, user, "用户名或密码错误")
	}
	return nil
}

/**
修改或添加用户
*/
func (u *UserServices) AddOrUpdateUser(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var user datamodels.User
	var roleId []uint
	var id uint
	err := mapstructure.Decode(maps["user"], &user)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["roleId"], &roleId)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["id"], &id)
	if err != nil {
		fmt.Println("map转struct出现异常：", err)
		return err
	}
	user.ID = id
	if user.ID > 0 {
		user.CreatedAt = *tool.GetTimeNow()
	}
	bool := repoUser.SaveUser(&user, roleId)
	if bool {
		datamodels.GetResult(result, user)
	} else {
		datamodels.GetResult(result, user, msgError)
	}
	return nil
}

//登录账号查重
func (u *UserServices) FindLoginName(ctx context.Context, user *datamodels.User, result *datamodels.Result) error {
	repoUser.FindUserName(user)
	if user.ID > 0 {
		datamodels.GetResult(result, user, "该用户已存在")
	} else {
		datamodels.GetResult(result, user)
	}
	return nil
}

//根据用户ID查询所拥有的菜单权限
func (u *UserServices) FindMenuRole(ctx context.Context, user *datamodels.User, result *datamodels.Result) error {
	//根据用户ID查询角色
	roleSlice := repoUser.BasisUserIdFindRole(*user)
	//根据角色查询菜单权限
	menuSlice := repoUser.FindOneMenuRole(roleSlice)
	datamodels.GetResult(result, menuSlice)
	return nil
}

//获取所有菜单
func (u *UserServices) GetMenuAll(ctx context.Context, null, result *datamodels.Result) error {
	menuSlice := repoUser.GetMenuAll()
	datamodels.GetResult(result, menuSlice)
	return nil
}

//获取所有角色
func (u *UserServices) GetRoleAll(ctx context.Context, null, result *datamodels.Result) error {
	roleSlice := repoUser.GetRoleAll()
	datamodels.GetResult(result, roleSlice)
	return nil
}

// 角色的添加和修改以及分配对应的菜单
func (u *UserServices) AddOrUpdateRole(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var role datamodels.Role
	var id uint
	err := mapstructure.Decode(maps["role"], &role)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	var menuId []uint
	err = mapstructure.Decode(maps["menuId"], &menuId)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["id"], &id)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	role.ID = id
	bool := repoUser.SaveRole(role, menuId)
	if bool == false {
		datamodels.GetResult(result, bool, msgError)
		return nil
	}
	datamodels.GetResult(result, bool)
	return nil
}

//分页查询角色
func (u *UserServices) FindRoleAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var role datamodels.Role
	var page datamodels.Page
	fmt.Println("当前map的值为", maps)
	err := mapstructure.Decode(maps["role"], &role)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["page"], &page)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	roleSlice := repoUser.FindRoleAll(role, &page)
	datamodels.GetResult(result, roleSlice)
	return nil
}

//根据角色ID获取角色对应的菜单权限ID
func (u *UserServices) FindMenuIdAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var role datamodels.Role
	var id uint
	fmt.Println("当前map的值为", maps)
	err := mapstructure.Decode(maps["id"], &id)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	role.ID = id
	fmt.Println(role)
	var roleSlice []datamodels.Role
	roleSlice = append(roleSlice, role)
	menuSlice := repoUser.FindOneMenuRole(roleSlice)
	var menuId []uint
	for _, v := range menuSlice {
		menuId = append(menuId, v.ID)
	}
	datamodels.GetResult(result, menuId)
	return nil
}

//项目的添加和修改
func (u *UserServices) AddOrUpdateProject(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var project datamodels.Project
	var id uint
	var userId uint
	err := mapstructure.Decode(maps, &project)
	if err != nil {
		fmt.Println("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["id"], &id)
	if err != nil {
		fmt.Println("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["UserId"], &userId)
	if err != nil {
		fmt.Println("map转struct出现异常：", err)
		return err
	}
	project.ID = id
	fmt.Println(project)
	if id > 0 {
		project.CreatedAt = *tool.GetTimeNow()
	}
	up := datamodels.UserProject{}
	up.UserId = userId
	bool := repoUser.SaveProject(project, up)
	if bool {
		datamodels.GetResult(result, bool)
		return nil
	}
	datamodels.GetResult(result, bool, msgError)
	return nil
}

//项目的环境配置的添加和修改
func (u *UserServices) AddOrUpdateSetting(ctx context.Context, maps []interface{}, result *datamodels.Result) error {
	var setting []datamodels.Setting
	//将切片中的interface转成需要的结构体
	for _, v := range maps {
		var s datamodels.Setting
		jsonStr, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		json.Unmarshal(jsonStr, &s)
		setting = append(setting, s)
	}
	bool, s := repoUser.SaveSetting(setting)
	if bool {
		datamodels.GetResult(result, s)
		return nil
	}
	datamodels.GetResult(result, bool, msgError)
	return nil
}

//项目的环境配置的查询
func (u *UserServices) FindProjectSetting(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	fmt.Println(maps)
	var projectId uint
	err := mapstructure.Decode(maps["projectId"], &projectId)
	if err != nil {
		fmt.Println("map转struct出现异常：", err)
		return err
	}
	fmt.Println(projectId)
	settingSlice := repoUser.FindProjectIdSetting(projectId)
	datamodels.GetResult(result, settingSlice)
	return nil
}

func GetNSQWarehouseProjectName() {
	maps := make(map[string]interface{})
	var m map[string]interface{}
	var warehouseName string
	//消费者
	consumer, err := nsq.NewConsumer("publish_system", "publish_system_one", nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	//生产者
	//producer, err := nsq.NewProducer(tool.ConfJson.NsqAddr, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	// 设置消息处理函数
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		fmt.Println(string(message.Body))
		if string(message.Body) == "" {
			return nil
		} else {
			err := json.Unmarshal([]byte(string(message.Body)), &m)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			warehouseName = m["repository_name"].(string)
			log.Info("仓库名为:", m["repository_name"])
		}
		project := repoUser.FindWarehouseNameProject(warehouseName)
		setting := repoUser.FindProjectIdSetting(project.ID)
		user := repoUser.FindProjectIdToUserId(project.ID)
		s := strings.Split(user.ProjectWorkPath, ",")
		fmt.Println(s)
		if project.ProjectTypeId == 1 {
			user.ProjectWorkPath = s[0]
		}
		if project.ProjectTypeId == 2 {
			user.ProjectWorkPath = s[1]
		}
		project.Setting = setting
		project.User = user
		if project.ProjectTypeId == 0 {
			log.Info("project是空对象")
			return nil
		}
		maps["project"] = project
		//err := GetClient().Call(context.Background(), "ExecAllCommand", &maps, datamodels.Result{})
		//if err != nil {
		//	log.Error("ERROR failed to call: %v", err)
		//	return nil
		//}
		//jsonBytes, err := json.Marshal(project)
		//if err != nil {
		//	log.Error("对象转JSON出错:", err)
		//	return nil
		//}
		//if len(jsonBytes) == 0 {
		//	log.Error("Json数据为空")
		//	return nil
		//}
		//if err := producer.Publish("publish_core", jsonBytes); err != nil {
		//	log.Fatal("publish error: " + err.Error())
		//}
		//log.Info("发送成功")
		return nil
	}))
	if err := consumer.ConnectToNSQD(tool.ConfJson.NsqAddr); err != nil {
		log.Fatal(err)
		return
	}
}

//提供给另外一个服务的接口 返回map
func (u *UserServices) FindWarehouseProject(ctx context.Context, maps *map[string]interface{}, result *datamodels.Result) error {
	fmt.Println("maps的值为", maps)
	var warehouseName string
	err := mapstructure.Decode((*maps)["repository_name"], &warehouseName)
	if err != nil {
		fmt.Println("map转struct出现异常：", err)
		return err
	}
	project := repoUser.FindWarehouseNameProject(warehouseName)
	setting := repoUser.FindProjectIdSetting(project.ID)
	user := repoUser.FindProjectIdToUserId(project.ID)
	project.Setting = setting
	//s := strings.Split(user.ProjectWorkPath, ",")
	//fmt.Println(s)
	//if len(s) == 1 {
	//	if project.ProjectTypeId == 1 {
	//		user.ProjectWorkPath = s[0]
	//	}
	//}
	//if len(s) == 2 {
	//	if project.ProjectTypeId == 1 {
	//		user.ProjectWorkPath = s[0]
	//	} else if project.ProjectTypeId == 2 {
	//		user.ProjectWorkPath = s[1]
	//	}
	//}
	project.User = user
	mapsTwo := make(map[string]interface{})
	systemSetting := repoUser.FindSystemSetting()
	mapsTwo["project"] = project
	mapsTwo["systemSetting"] = systemSetting
	datamodels.GetResult(result, mapsTwo)
	return nil
}

/**
查询仓库名是否重复
*/
func (u *UserServices) FindWarehouseName(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	fmt.Println("maps的值为", maps)
	var warehouseName string
	err := mapstructure.Decode(maps["warehouseName"], &warehouseName)
	if err != nil {
		fmt.Println("map转struct出现异常：", err)
		return err
	}
	project := repoUser.FindWarehouseNameProject(warehouseName)
	if project.ID > 0 {
		datamodels.GetResult(result, false, "该仓库名已经存在,不能使用这个仓库名")
		return nil
	}
	datamodels.GetResult(result, true)
	return nil
}

/**
查询角色名是否重复
*/
func (u *UserServices) FindRoleName(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	fmt.Println("maps的值为", maps)
	var roleName string
	err := mapstructure.Decode(maps["roleName"], &roleName)
	if err != nil {
		fmt.Println("map转struct出现异常：", err)
		return err
	}
	role := repoUser.FindRoleName(roleName)
	if role.ID > 0 {
		datamodels.GetResult(result, false, "该角色名已经存在,不能使用这个角色名")
		return nil
	}
	datamodels.GetResult(result, true)
	return nil
}

/**
项目的查询接口
*/
func (u *UserServices) FindProjectAll(ctx context.Context, maps map[string]interface{}, result *datamodels.ResultPage) error {
	var project datamodels.Project
	var page datamodels.Page
	var projectTypeId string
	var userId uint
	err := mapstructure.Decode(maps["project"], &project)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["page"], &page)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["projectTypeId"], &projectTypeId)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	err = mapstructure.Decode(maps["userId"], &userId)
	if err != nil {
		log.Error("map转struct出现异常：", err)
		return err
	}
	if projectTypeId != "" {
		i, _ := strconv.Atoi(projectTypeId)
		project.ProjectTypeId = uint(i)
	}
	//根据用户ID查询出对应的项目列表ID
	up := repoUser.FindUserIdToProjectIds(userId)
	var projectId []uint
	for _, v := range up {
		projectId = append(projectId, v.ProjectId)
	}
	projectSlice := repoUser.FindProjectAll(project, &page, projectId)
	fmt.Println("当前project切片数组", projectSlice)
	for i, v := range projectSlice {
		fmt.Println("当前project对象", v)
		if v.ProjectTypeId == 2 {
			projectName := v.WarehouseName[strings.Index(v.WarehouseName, "/")+1:]
			fmt.Println("当前项目名", projectName)
			pid := tool.GetPidByProcessName(projectName)
			fmt.Println("当前PID为", pid)
			if pid != "" {
				projectSlice[i].Enable = 1
			} else {
				projectSlice[i].Enable = 2
			}
		}
	}
	datamodels.GetResultPage(result, projectSlice, page.PageCount)
	return nil
}

//项目类型查询接口
func (u *UserServices) FindProjectTypeAll(ctx context.Context, null, result *datamodels.Result) error {
	p := repoUser.FindProjectTypeAll()
	datamodels.GetResult(result, p)
	return nil
}

//用户类型查询接口
func (u *UserServices) FindUserTypeAll(ctx context.Context, null, result *datamodels.Result) error {
	p := repoUser.FindUserTypeAll()
	datamodels.GetResult(result, p)
	return nil
}

//用户删除接口
func (u *UserServices) DeleteUserAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var userId []uint
	mapstructure.Decode(maps["userId"], &userId)
	bool := repoUser.DeleteUserAll(userId)
	if bool {
		datamodels.GetResult(result, bool)
		return nil
	}
	datamodels.GetResult(result, bool, msgError)
	return nil

}

//项目删除接口
func (u *UserServices) DeleteProjectAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var projectId []uint
	mapstructure.Decode(maps["projectId"], &projectId)
	bool := repoUser.DeleteProjectAll(projectId)
	if bool {
		datamodels.GetResult(result, bool)
		return nil
	}
	datamodels.GetResult(result, bool, msgError)
	return nil
}

//角色删除接口
func (u *UserServices) DeleteRoleAll(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var roleId []uint
	mapstructure.Decode(maps["roleId"], &roleId)
	bool := repoUser.DeleteRoleAll(roleId)
	if bool {
		datamodels.GetResult(result, bool)
		return nil
	}
	datamodels.GetResult(result, bool, msgError)
	return nil
}

//系统配置的添加和修改
func (u *UserServices) SaveSystemSetting(ctx context.Context, maps []interface{}, result *datamodels.Result) error {
	var systemSetting []datamodels.SystemSetting
	//将切片中的interface转成需要的结构体
	for _, v := range maps {
		var s datamodels.SystemSetting
		jsonStr, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err)
			return err
		}
		json.Unmarshal(jsonStr, &s)
		systemSetting = append(systemSetting, s)
	}
	bool := repoUser.SaveSystemSetting(systemSetting)
	if bool {
		datamodels.GetResult(result, bool)
		return nil
	}
	datamodels.GetResult(result, bool, msgError)
	return nil
}

//系统配置的查询
func (u *UserServices) FindSystemSetting(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	slice := repoUser.FindSystemSetting()
	datamodels.GetResult(result, slice)
	return nil
}

//启动项目
func (u *UserServices) ProjectStart(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	log.Info(maps)
	var projectId uint
	err := mapstructure.Decode(maps["projectId"], &projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	//根据项目ID查询对应的用户信息
	user := repoUser.FindProjectIdToUserId(projectId)
	systemSlice := repoUser.FindSystemSetting()
	//取出Project地址
	var projectAddr string
	for _, v := range systemSlice {
		if v.Keys == "ProjectPath" {
			projectAddr = v.Values
			break
		}
	}
	project := repoUser.FindProjectId(projectId)
	mapsTwo := make(map[string]interface{})
	mapsTwo["userName"] = user.LoginName
	mapsTwo["projectAddr"] = projectAddr
	mapsTwo["warehouseName"] = project.WarehouseName
	jsonBytes, err := json.Marshal(mapsTwo)
	fmt.Println("Byte[]===", jsonBytes)
	if err != nil {
		log.Error("map转byte[]出现异常")
		datamodels.GetResult(result, true, msgError)
		return nil
	}
	if err := producer.Publish("publish_core_start", jsonBytes); err != nil {
		log.Error("publish error: " + err.Error())
		datamodels.GetResult(result, true, msgError)
		return nil
	}
	datamodels.GetResult(result, true)
	return nil
}

//停止项目
func (u *UserServices) ProjectStop(ctx context.Context, maps map[string]interface{}, result *datamodels.Result) error {
	var warehouseName string
	err := mapstructure.Decode(maps["warehouseName"], &warehouseName)
	if err != nil {
		log.Error(err)
		datamodels.GetResult(result, false, msgError)
		return err
	}
	mapsTwo := make(map[string]interface{})
	mapsTwo["warehouseName"] = warehouseName
	jsonBytes, err := json.Marshal(mapsTwo)
	if err != nil {
		log.Error("map转byte[]出现异常")
		datamodels.GetResult(result, true, msgError)
		return nil
	}
	if err := producer.Publish("publish_core_stop", jsonBytes); err != nil {
		log.Error("publish error: " + err.Error())
		datamodels.GetResult(result, true, msgError)
		return nil
	}
	datamodels.GetResult(result, true)
	return nil
}
