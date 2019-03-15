package repository

import (
	"fmt"
	"github.com/prometheus/common/log"
	"publish_server_user/datamodels"
	"publish_server_user/datasource"
	"publish_server_user/tool"
)

func NewUserRepo() UserRepo {
	return &userRepo{}
}

type userRepo struct{}
type UserRepo interface {
	FindUserNamePassword(user *datamodels.User)
	FindUserName(user *datamodels.User)
	FindUserId(user *datamodels.User)
	FindRoleName(roleName string) (role datamodels.Role)
	FindUserAll(user datamodels.User, page *datamodels.Page) (userSlice []datamodels.User)
	SaveUser(user *datamodels.User, roleId []uint) (bool bool)
	UpdateUserPassword(userId uint, newPassword string) (bool bool)
	BasisUserIdFindRole(user datamodels.User) (roleSlice []datamodels.Role)
	FindOneMenuRole(roleSlice []datamodels.Role) (menuSlice []datamodels.Menu)
	GetMenuAll() (menuSlice []datamodels.Menu)
	GetRoleAll() (menuSlice []datamodels.Role)
	SaveRole(role datamodels.Role, menuId []uint) (bool bool)
	FindRoleAll(role datamodels.Role, page *datamodels.Page) (roleSlice []datamodels.Role)
	SaveProject(project datamodels.Project, up datamodels.UserProject) (bool bool)
	SaveSetting(setting []datamodels.Setting) (bool bool, s []datamodels.Setting)
	FindProjectIdSetting(projectId uint) (setting []datamodels.Setting)
	FindProjectAll(project datamodels.Project, page *datamodels.Page, id []uint) (projectSlice []datamodels.Project)
	FindWarehouseNameProject(WarehouseName string) (project datamodels.Project)
	FindProjectTypeAll() (p []datamodels.ProjectType)
	FindUserTypeAll() (ut []datamodels.UserType)
	FindUserIdToProjectIds(userId uint) (up []datamodels.UserProject)
	FindProjectIdToUserId(projectId uint) (user datamodels.User)
	DeleteUserAll(userId []uint) (bool bool)
	DeleteProjectAll(projectId []uint) (bool bool)
	FindProjectId(projectId uint) (project datamodels.Project)
	DeleteRoleAll(roleId []uint) (bool bool)
	SaveSystemSetting(systemSetting []datamodels.SystemSetting) (bool bool)
	FindSystemSetting() (systemSetting []datamodels.SystemSetting)
}

/**
根据用户ID查询用户信息
*/
func (u *userRepo) FindUserId(user *datamodels.User) {
	datasource.GetDB().Where("id = ? ", user.ID).First(&user)
}

/**
根据用户登录名和密码查询用户 用于登录
*/
func (u *userRepo) FindUserNamePassword(user *datamodels.User) {
	datasource.GetDB().Where("login_name = ? and password = ?", user.LoginName, user.Password).Preload("UserType").First(&user)
}

/**
根据用户登录名查询用户 用于注册
*/
func (u *userRepo) FindUserName(user *datamodels.User) {
	datasource.GetDB().Where("login_name = ?", user.LoginName).First(&user)
}

/**
用户修改密码
*/
func (u *userRepo) UpdateUserPassword(userId uint, newPassword string) (bool bool) {
	err := datasource.GetDB().Exec("update `User` set password = ? where id = ?", newPassword, userId).Error
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

//角色名查重
func (u *userRepo) FindRoleName(roleName string) (role datamodels.Role) {
	datasource.GetDB().Where("name = ?", roleName).First(&role)
	return
}

/**
查询所有用户信息
*/
func (u *userRepo) FindUserAll(user datamodels.User, page *datamodels.Page) (userSlice []datamodels.User) {
	db := datasource.GetDB()
	if user.Name != "" {
		db = db.Where("name = ?", user.Name)
	}
	if user.LoginName != "" {
		db = db.Where("login_name = ?", user.LoginName)
	}
	if user.Sex > 0 {
		db = db.Where("sex = ?", user.Sex)
	}
	if user.ID > 0 {
		db = db.Where("id = ?", user.ID)
	}
	var count int
	db.Model(&datamodels.User{}).Count(&count)
	page.PageCount = count
	db.Preload("UserType").Limit(page.PageSize).Offset(datamodels.PageIndex(page)).Find(&userSlice)
	return userSlice
}

/**
修改或添加用户信息,以及对应的角色
*/
func (u *userRepo) SaveUser(user *datamodels.User, roleId []uint) (bool bool) {
	tx := datasource.GetDB().Begin()
	err := tx.Save(&user).Error
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return false
	}
	//添加权限,先删除用户之前的角色
	err = tx.Exec("delete from user_role where user_id = ?", user.ID).Error
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return false
	}
	//再添加角色
	for _, v := range roleId {
		var userRole datamodels.UserRole
		userRole.RoleId = v
		userRole.UserId = user.ID
		err = tx.Save(&userRole).Error
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return false
		}
	}
	tx.Commit()
	return true
}

/**
根据用户ID查询所对应的角色
*/
func (u *userRepo) BasisUserIdFindRole(user datamodels.User) (roleSlice []datamodels.Role) {
	rows, err := datasource.GetDB().Raw("select * from role where id in (select role_id from user_role where user_id = ?)", user.ID).Rows()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	for rows.Next() {
		var role datamodels.Role
		rows.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt, &role.Name, &role.Remarks)
		roleSlice = append(roleSlice, role)
	}
	return
}

/**
根据角色ID查询菜单权限
*/
func (u *userRepo) FindOneMenuRole(roleSlice []datamodels.Role) (menuSlice []datamodels.Menu) {
	var id []uint
	for _, v := range roleSlice {
		id = append(id, v.ID)
	}
	rows, err := datasource.GetDB().Raw("select * from menu where id in (select menu_id from role_menu where role_id in (?)) ORDER BY menu_id ", id).Rows()
	if err != nil {
		log.Error(err)
		return nil
	}
	for rows.Next() {
		var menu datamodels.Menu
		rows.Scan(&menu.ID, &menu.CreatedAt, &menu.UpdatedAt, &menu.DeletedAt, &menu.MenuId, &menu.MenuName, &menu.MenuUrl, &menu.Level, &menu.Superior, &menu.Icon, &menu.PermissionIdentifier)
		menuSlice = append(menuSlice, menu)
	}
	return
}

//获取所有的菜单
func (u *userRepo) GetMenuAll() (menuSlice []datamodels.Menu) {
	datasource.GetDB().Find(&menuSlice)
	return
}

//获取所有的角色
func (u *userRepo) GetRoleAll() (roleSlice []datamodels.Role) {
	datasource.GetDB().Find(&roleSlice)
	return
}

//添加或者修改角色以及菜单权限
func (u *userRepo) SaveRole(role datamodels.Role, menuId []uint) (bool bool) {
	tx := datasource.GetDB().Begin()
	err := tx.Save(&role).Error
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return false
	}
	//先删除之前的菜单
	err = tx.Exec("delete from role_menu where role_id = ?", role.ID).Error
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return false
	}
	for _, v := range menuId {
		var roleMenu datamodels.RoleMenu
		roleMenu.RoleId = role.ID
		roleMenu.MenuId = v
		err = tx.Save(&roleMenu).Error
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return false
		}
	}
	tx.Commit()
	return true
}

//查询角色
func (u *userRepo) FindRoleAll(role datamodels.Role, page *datamodels.Page) (roleSlice []datamodels.Role) {
	db := datasource.GetDB()
	if role.Name != "" {
		db = db.Where("name like ?", "%"+role.Name+"%")
	}
	var count int
	db.Model(&datamodels.Role{}).Count(&count)
	db.Limit(page.PageSize).Offset(datamodels.PageIndex(page)).Find(&roleSlice)
	return
}

//项目的添加和修改
func (u *userRepo) SaveProject(project datamodels.Project, up datamodels.UserProject) (bool bool) {
	tx := datasource.GetDB().Begin()
	if project.ID == 0 {
		//ID==0时,执行添加操作
		err := tx.Save(&project).Error
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return false
		}
		//当添加项目成功后添加项目ID和UserId到中间表,添加一条用户项目表记录
		up.ProjectId = project.ID
		err = tx.Save(&up).Error
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return false
		}
	} else {
		//单纯的执行修改操作
		err := tx.Save(&project).Error
		if err != nil {
			log.Error(err)
			tx.Rollback()
			return false
		}
	}
	tx.Commit()
	return true
}

//项目的配置的添加和修改
func (u *userRepo) SaveSetting(setting []datamodels.Setting) (bool bool, s []datamodels.Setting) {
	tx := datasource.GetDB().Begin()
	for _, v := range setting {
		if v.ID > 0 {
			v.CreatedAt = *tool.GetTimeNow()
		}
		err := tx.Save(&v).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return false, s
		}
		s = append(s, v)
	}
	tx.Commit()
	return true, s
}

//根据项目ID查询对应的配置信息
func (u *userRepo) FindProjectIdSetting(projectId uint) (setting []datamodels.Setting) {
	datasource.GetDB().Where("project_id = ?", projectId).Find(&setting)
	return
}

//根据仓库名查询项目
func (u *userRepo) FindWarehouseNameProject(WarehouseName string) (project datamodels.Project) {
	datasource.GetDB().Where("warehouse_name = ?", WarehouseName).First(&project)
	return
}

//项目的查询接口
func (u *userRepo) FindProjectAll(project datamodels.Project, page *datamodels.Page, id []uint) (projectSlice []datamodels.Project) {
	db := datasource.GetDB()
	if project.Name != "" {
		db = db.Where("name like ?", "%"+project.Name+"%")
	}
	if project.ProjectTypeId > 0 {
		db = db.Where("project_type_id = ?", project.ProjectTypeId)
	}
	if project.WarehouseName != "" {
		db = db.Where("warehouse_name like ?", "%"+project.WarehouseName+"%")
	}
	db = db.Where("id in (?)", id)
	var count int
	db.Model(&datamodels.Project{}).Count(&count)
	page.PageCount = count
	db.Preload("ProjectType").Limit(page.PageSize).Offset(datamodels.PageIndex(page)).Find(&projectSlice)
	return
}

//根据项目ID查询项目
func (u *userRepo) FindProjectId(projectId uint) (project datamodels.Project) {
	datasource.GetDB().Where("id = ?", projectId).Preload("ProjectType").First(&project)
	return
}

//项目类型的查询
func (u *userRepo) FindProjectTypeAll() (p []datamodels.ProjectType) {
	datasource.GetDB().Find(&p)
	return
}

//根据UserId查询项目ID
func (u *userRepo) FindUserIdToProjectIds(userId uint) (up []datamodels.UserProject) {
	db := datasource.GetDB()
	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}
	db.Find(&up)
	return
}

//根据项目ID查询用户ID
func (u *userRepo) FindProjectIdToUserId(projectId uint) (user datamodels.User) {
	var userProject datamodels.UserProject
	datasource.GetDB().Where("project_id = ?", projectId).First(&userProject)
	datasource.GetDB().Where("id = ?", userProject.UserId).First(&user)
	return
}

//查询用户类型
func (u *userRepo) FindUserTypeAll() (ut []datamodels.UserType) {
	datasource.GetDB().Find(&ut)
	return
}

//删除用户&&会删除用户关联下的所有数据！！！请三思而后行
func (u *userRepo) DeleteUserAll(userId []uint) (bool bool) {
	tx := datasource.GetDB().Begin()
	//1.删除用户
	err := tx.Exec("delete from user where id in (?)", userId).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return false
	}
	//2.删除角色权限
	err = tx.Exec("delete from user_role where user_id in (?)", userId).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return false
	}
	//3.删除用户下对应的项目
	var up []datamodels.UserProject
	//3.1先查询出用户对应的项目ID
	tx.Where("user_id in (?)", userId).Find(&up)
	var projectId []uint
	for _, v := range up {
		projectId = append(projectId, v.ProjectId)
	}
	//3.2 删除项目
	err = tx.Exec("delete from project where id in (?)", projectId).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return false
	}
	//4.删除用户项目的中间表数据
	err = tx.Exec("delete from user_project where user_id in (?)", userId).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return false
	}
	//5.删除项目对应的环境配置表数据
	err = tx.Exec("delete from setting where project_id in (?)", projectId).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return false
	}
	tx.Commit()
	return true
}

//删除项目&&会删除项目关联下的所有数据！！！请三思而后行
func (u *userRepo) DeleteProjectAll(projectId []uint) (bool bool) {
	tx := datasource.GetDB().Begin()
	//1.删除项目
	err := tx.Exec("delete from project where id in (?)", projectId).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return false
	}
	//2.删除项目对应的用户和项目的中间表
	err = tx.Exec("delete from user_project where project_id in (?)", projectId).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return false
	}
	//3.删除项目下的环境配置表
	err = tx.Exec("delete from setting where project_id in (?)", projectId).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return false
	}
	tx.Commit()
	return true
}

//删除角色和对应的菜单权限
func (u *userRepo) DeleteRoleAll(roleId []uint) (bool bool) {
	tx := datasource.GetDB().Begin()
	err := tx.Exec("delete from role where id in (?)", roleId).Error
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return false
	}
	err = tx.Exec("delete from role_menu where role_id in (?)", roleId).Error
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

//添加或修改系统配置
func (u *userRepo) SaveSystemSetting(systemSetting []datamodels.SystemSetting) (bool bool) {
	tx := datasource.GetDB().Begin()
	for _, v := range systemSetting {
		if v.ID > 0 {
			v.CreatedAt = *tool.GetTimeNow()
		}
		err := tx.Save(&v).Error
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return false
		}
	}
	tx.Commit()
	return true
}

//查询系统配置
func (u *userRepo) FindSystemSetting() (systemSetting []datamodels.SystemSetting) {
	datasource.GetDB().Find(&systemSetting)
	return
}
