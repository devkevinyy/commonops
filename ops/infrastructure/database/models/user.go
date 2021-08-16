package models

import (
	"errors"
	"fmt"
	"github.com/chujieyang/commonops/ops/opslog"
	"strings"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/jinzhu/gorm"
)

var authSkipMap = map[string]bool{
	"GET:/user/tokenRefresh": true,
	"GET:/user/permissions": true,
	"GET:/user/login": true,
}

/*
	系统用户表
*/
type User struct {
	gorm.Model
	UserId       string `json:"user_id" gorm:"user_id"`
	UnionId      string `json:"union_id" gorm:"union_id"`
	Mobile       string `json:"mobile" gorm:"mobile"`
	IsLeader     bool   `json:"is_leader" gorm:"is_leader"`
	UserName     string `json:"username" gorm:"column:user_name"`
	Password     string `json:"password" gorm:"column:password"`
	EmpId        string `json:"empId" gorm:"column:emp_id"`
	DepartmentId int    `json:"department_id" gorm:"column:department_id"`
	Active       bool   `json:"active" gorm:"active"`
	Position     string `json:"position" gorm:"position"`
	Email        string `json:"email" gorm:"email"`
	FirstIn      int8   `json:"first_in" gorm:"first_in;not null;default:1"`
}

/*
	系统角色表
*/
type Role struct {
	Id          int    `gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*
	用户-角色关系表
*/
type UserRoles struct {
	Id     int `gorm:"primary_key"`
	UserId int `json:"userId"`
	RoleId int `json:"roleId"`
}

/*
	系统权限表
*/
type Permissions struct {
	Id          int    `gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UrlPath     string `json:"urlPath"`
	CanDelete   int  `json:"canDelete"`
	AuthType       string `json:"authType"`
}

/*
	角色-权限关系表
*/
type RolePermissions struct {
	Id           int `gorm:"primary_key"`
	RoleId       int `json:"roleId"`
	PermissionId int `json:"permissionId"`
}

/*
	角色-资源关系表
*/
type RoleResources struct {
	Id           int    `gorm:"primary_key"`
	RoleId       int    `json:"roleId"`
	ResourceType string `json:"resourceType"`
	ResourceId   int    `json:"resourceId"`
}

func GetUsersCount() (total int) {
	total = 0
	database.Mysql().Model(&User{}).Where("active=1").Count(&total)
	return
}

func GetUserInfo(email string) (user User, err error) {
	err = database.Mysql().First(&user, " active = 1 and email=?", email).Error
	return
}

func GetUserInfoByEmpId(empId string) (user User, err error) {
	err = database.Mysql().First(&user, " active = 1 and emp_id=?", empId).Error
	return
}

func UpdateUserFirstInStatus(id uint) {
	database.Mysql().Exec("update users set first_in=0 where id=?", id)
	return
}

func UpdateUserPassword(empId string, password string) (user User, err error) {
	err = database.Mysql().First(&user, "emp_id=?", empId).Error
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	user.Password = password
	err = database.Mysql().Save(&user).Error
	return
}

func GetUserList(offset int, pageSize int) (users []User) {
	users = []User{}
	database.Mysql().Offset(offset).Limit(pageSize).Find(&users).Where("active=1")
	return
}

func GetRolesCount() (total int) {
	total = 0
	database.Mysql().Model(&Role{}).Count(&total)
	return
}

func GetRolesList(offset int, pageSize int) (roles []Role) {
	roles = []Role{}
	database.Mysql().Offset(offset).Limit(pageSize).Find(&roles)
	return
}

func GetUsersByRoleId(roleId int) (users []User, err error) {
	err = database.Mysql().Raw("SELECT ur.user_id as id FROM user_roles as ur "+
		"WHERE ur.role_id = ?", roleId).Scan(&users).Error
	return
}

func GetAllUsers() (users []User, err error) {
	err = database.Mysql().Raw("SELECT id, user_name, position FROM users where active = 1").Scan(&users).Error
	return
}

func AddNewRole(roleName string, description string) (err error) {
	role := Role{
		Name:        roleName,
		Description: description,
	}
	count := 0
	database.Mysql().Where("name=?", roleName).Find(&role).Count(&count)
	if count == 0 {
		database.Mysql().Create(&role)
		return
	} else {
		return errors.New("角色名称已经存在")
	}
}

func UpdateRole(id int, roleName string, description string) (err error) {
	err = database.Mysql().Exec("update roles set name=?, description=? where id=?",
		roleName, description, id).Error
	return
}

func DeleteRole(id int) (err error) {
	tx := database.Mysql().Begin()
	err1 := tx.Delete(&Role{}, "id=?", id).Error
	err2 := tx.Exec("delete from user_roles where role_id = ?", id).Error
	err3 := tx.Exec("delete from role_resources where role_id = ?", id).Error
	err4 := tx.Exec("delete from role_permissions where role_id = ?", id).Error
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		tx.Rollback()
		if err1 != nil {
			err = err1
		} else if err2 != nil {
			err = err2
		} else if err3 != nil {
			err = err3
		} else {
			err = err4
		}
		return
	}
	tx.Commit()
	return
}

func AddRoleUsers(roleId int64, userIdList []string) (err error) {
	users := strings.Join(userIdList, ",")
	err = database.Mysql().Exec("delete from user_roles where role_id = ? "+
		"and user_id not in (?)", roleId, users).Error
	for _, userId := range userIdList {
		exist := 0
		database.Mysql().Raw("select count(*) from user_roles "+
			"where role_id = ? and user_id = ?", roleId, userId).Count(&exist)
		if exist == 0 {
			err = database.Mysql().Exec("insert into user_roles(role_id, user_id) "+
				"values(?, ?)", roleId, userId).Error
		}
	}
	return
}

func GetAllAuthLinks() (permissions []Permissions, err error) {
	err = database.Mysql().Raw("SELECT id, name, url_path, can_delete, auth_type FROM permissions").Scan(&permissions).Error
	return
}

func GetAuthLinksByRoleId(roleId int) (permissions []Permissions, err error) {
	err = database.Mysql().Raw("SELECT permission_id as id FROM role_permissions "+
		"WHERE role_id = ?", roleId).Scan(&permissions).Error
	return
}

func AddAuthLink(name string, description string, path string) (err error) {
	permission := Permissions{
		Name:        name,
		Description: description,
		UrlPath:     path,
		CanDelete: 0,
	}
	count := 0
	database.Mysql().Where("name = ? or url_path = ?", name, path).Find(&permission).Count(&count)
	if count == 0 {
		database.Mysql().Create(&permission)
		return
	} else {
		return errors.New("权限名称或路径已经存在")
	}
}

func DeleteAuthLink(id int) (err error) {
	err = database.Mysql().Delete(&Permissions{}, "id=?", id).Error
	return
}

func AddRoleAuthLinks(roleId int64, authLinkIdList []string) (err error) {
	links := strings.Join(authLinkIdList, ",")
	err = database.Mysql().Exec("delete from role_permissions where role_id = ? "+
		"and permission_id not in (?)", roleId, links).Error
	for _, linkId := range authLinkIdList {
		exist := 0
		database.Mysql().Raw("select count(*) from role_permissions "+
			"where role_id = ? and permission_id = ?", roleId, linkId).Count(&exist)
		if exist == 0 {
			err = database.Mysql().Exec("insert into role_permissions(role_id, permission_id) "+
				"values(?, ?)", roleId, linkId).Error
		}
	}
	return
}

func GetPermissionsCount() (total int) {
	total = 0
	database.Mysql().Model(&Permissions{}).Count(&total)
	return
}

func GetPermissionsList(offset int, pageSize int) (permissions []Permissions) {
	permissions = []Permissions{}
	database.Mysql().Offset(offset).Limit(pageSize).Find(&permissions)
	return
}

func AddUser(username string, password string, email string, position string) (id uint, err error) {
	count := 0
	database.Mysql().Raw(fmt.Sprintf("select count(*) from users where email = '%s' ", email)).Count(&count)
	if count > 0 {
		err = errors.New("邮箱已经被注册")
		return
	}
	user := User{
		UserName: username,
		Password: utils.GenUserPassword(password),
		Email:    email,
		Active:   true,
		Position: position,
		EmpId:    utils.GetUUID(),
	}
	result := database.Mysql().Create(&user)
	return user.ID, result.Error
}

func UpdateUserActiveStatus(email string, active bool) (err error) {
	activeStatus := 0
	if active == true {
		activeStatus = 1
	}
	err = database.Mysql().Exec(fmt.Sprintf("update users set active = %d where email = '%s' ",
		activeStatus, email)).Error
	return
}

func GetUserPermissions(userId int, authType string) (permissionList []Permissions) {
	database.Mysql().Raw("select distinct(p.id) as id, p.name, p.url_path, p.can_delete, p.auth_type from "+
		"permissions as p inner join role_permissions as rp on p.id = rp.permission_id "+
		"inner join user_roles as ur on rp.role_id = ur.role_id "+
		"where ur.user_id = ? and p.auth_type = ? ", userId, authType).Scan(&permissionList)
	return
}

// 用户是否是系统超级管理员
func IsUserSuperAdmin(userId uint) (isSuperAdmin bool) {
	count := 0
	database.Mysql().Raw("select count(*) from "+
		"roles as r inner join user_roles as ur on r.id = ur.role_id "+
		"where r.name= '超级管理员' and ur.user_id = ?", userId).Count(&count)
	if count > 0 {
		isSuperAdmin = true
	} else {
		isSuperAdmin = false
	}
	return
}

// 用户是否是系统"超级管理员"或者"运维中心"员工
func IsUserSuperAdminOrOps(userId uint) (isSuperAdminOrOps bool) {
	count := 0
	database.Mysql().Raw("select count(*) from "+
		"roles as r inner join user_roles as ur on r.id = ur.role_id "+
		"where (r.name = '超级管理员' or r.name = '运维中心') and ur.user_id = ?", userId).Count(&count)
	if count > 0 {
		isSuperAdminOrOps = true
	} else {
		isSuperAdminOrOps = false
	}
	return
}

// 用户是否状态正常
func IsUserValid(userId int32) (isValid bool) {
	count := 0
	database.Mysql().Raw("select count(*) from users where active=1 and id = ?", userId).Count(&count)
	if count > 0 {
		isValid = true
	} else {
		isValid = false
	}
	return
}

func IsUserHasActionPermision(userId int32, authLink string) bool {
	if _, ok := authSkipMap[authLink]; ok {
		return true
	}
	count := 0
	database.Mysql().Raw("select count(distinct ur.user_id) from "+
		"permissions as p inner join role_permissions as rp on p.id = rp.permission_id "+
		"inner join user_roles as ur on rp.role_id = ur.role_id "+
		"where ur.user_id = ? and p.url_path = ?", userId, authLink).Count(&count)
	return count > 0
}
