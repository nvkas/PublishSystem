package controller

import (
	"context"
	"fmt"
	"github.com/kataras/iris"
	"log"
	"publish_client_user/datamodels"
)

//需要获取token认真的user控制层
type UsersController struct {
	Ctx iris.Context
}

/**
验证登录名是否存在
*/
func (c *UsersController) PostSelectLoginName() (result *datamodels.Result) {
	var user datamodels.User
	c.Ctx.ReadJSON(&user)
	fmt.Println(user)
	err := GetClient().Call(context.Background(), "FindLoginName", &user, &result)
	if err != nil {
		log.Printf("ERROR failed to call: %v", err)
		return
	}
	return
}

//根据用户ID查询菜单权限
func (c *UsersController) PostSelectMenu() (result *datamodels.Result) {
	var user datamodels.User
	c.Ctx.ReadJSON(&user)
	err := GetClient().Call(context.Background(), "FindMenuRole", &user, &result)
	if err != nil {
		log.Printf("ERROR failed to call: %v", err)
		return
	}
	return
}

//用户的添加修改操作
func (c *UsersController) PostSaveUser() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//err := GetClient().Call(context.Background(), "AddOrUpdateUser", &maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("AddOrUpdateUser", maps, result)
	return
}

//根据用户ID查询用户
func (c *UsersController) PostFindUserId() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//err := GetClient().Call(context.Background(), "FindUserId", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("FindUserId", maps, result)
	return
}

//用户的查询操作
func (c *UsersController) PostSelectUser() (result *datamodels.ResultPage) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//err := GetClient().Call(context.Background(), "FindUserAll", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = CallPage("FindUserAll", maps, result)
	return
}

//修改用户密码
func (c *UsersController) PostUpdateUserPassword() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//err := GetClient().Call(context.Background(), "UpdateUserPassword", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("UpdateUserPassword", maps, result)
	return
}

//获取所有菜单
func (c *UsersController) GetMenu() (result *datamodels.Result) {
	//err := GetClient().Call(context.Background(), "GetMenuAll", nil, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("GetMenuAll", maps, result)
	return
}

//获取所有角色
func (c *UsersController) GetRole() (result *datamodels.Result) {
	//err := GetClient().Call(context.Background(), "GetRoleAll", nil, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("GetRoleAll", maps, result)
	return
}

//角色查询
func (c *UsersController) PostFindRoleAll() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//err := GetClient().Call(context.Background(), "FindRoleAll", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("FindRoleAll", maps, result)
	return
}

//添加或者修改角色以及菜单权限
func (c *UsersController) PostSaveRole() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//err := GetClient().Call(context.Background(), "AddOrUpdateRole", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("AddOrUpdateRole", maps, result)
	return
}

//根据角色ID查询角色对应的菜单权限Id
func (c *UsersController) PostFindMenuId() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//err := GetClient().Call(context.Background(), "FindMenuIdAll", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("FindMenuIdAll", maps, result)
	return
}

//添加或修改一个项目
func (c *UsersController) PostSaveProject() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//err := GetClient().Call(context.Background(), "AddOrUpdateProject", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("AddOrUpdateProject", maps, result)
	return
}

//修改和添加一个项目的配置
func (c *UsersController) PostSaveSetting() (result *datamodels.Result) {
	var maps []interface{}
	c.Ctx.ReadJSON(&maps)
	fmt.Println(maps)
	err := GetClient().Call(context.Background(), "AddOrUpdateSetting", maps, &result)
	if err != nil {
		log.Printf("ERROR failed to call: %v", err)
		return result
	}
	return
}

//项目的配置的查询 根据项目Id
func (c *UsersController) PostFindProjectSetting() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//fmt.Println(maps)
	//err := GetClient().Call(context.Background(), "FindProjectSetting", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("FindProjectSetting", maps, result)
	return
}

//查询仓库名是否存在
func (c *UsersController) PostFindWarehouseName() (result *datamodels.Result) {
	//var maps map[string]interface{}
	//c.Ctx.ReadJSON(&maps)
	//fmt.Println(maps)
	//err := GetClient().Call(context.Background(), "FindWarehouseName", maps, &result)
	//if err != nil {
	//	log.Printf("ERROR failed to call: %v", err)
	//	return result
	//}
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("FindWarehouseName", maps, result)
	return
}

//查询角色名是否存在
func (c *UsersController) PostFindRoleName() (result *datamodels.Result) {
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("FindRoleName", maps, result)
	return
}

//项目查询接口
func (c *UsersController) PostFindProjectAll() (result *datamodels.ResultPage) {
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = CallPage("FindProjectAll", maps, result)
	return
}

//项目类型查询接口
func (c *UsersController) GetFindProjectTypeAll() (result *datamodels.Result) {
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("FindProjectTypeAll", maps, result)
	return
}

//用户类型查询接口
func (c *UsersController) GetFindUserTypeAll() (result *datamodels.Result) {
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("FindUserTypeAll", maps, result)
	return
}

//用户删除
func (c *UsersController) PostDeleteUserAll() (result *datamodels.Result) {
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("DeleteUserAll", maps, result)
	return
}

//项目删除
func (c *UsersController) PostDeleteProjectAll() (result *datamodels.Result) {
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("DeleteProjectAll", maps, result)
	return
}

//角色删除
func (c *UsersController) PostDeleteRoleAll() (result *datamodels.Result) {
	var maps map[string]interface{}
	c.Ctx.ReadJSON(&maps)
	result = Call("DeleteRoleAll", maps, result)
	return
}
