package models

import (
	"fmt"
	"strings"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type DmsAuth struct {
	gorm.Model
	DataStatus   int    `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	EmpId        int    `json:"EmpId" gorm:"column:emp_id;type:int;not null"`
	Username     string `json:"Username" gorm:"column:username;type:varchar(256);not null"`
	AuthType     int    `json:"AuthType" gorm:"column:auth_type;type:tinyint;not null;default:1"` // 1-实例，2-库(现在都是2)，3-表，4-表操作
	InstanceType string `json:"InstanceType" gorm:"column:instance_type;type:varchar(256);not null"`
	InstanceId   string `json:"InstanceId" gorm:"column:instance_id;type:varchar(512);not null"`
	InstanceName string `json:"InstanceName" gorm:"column:instance_name;type:varchar(512);not null"`
	DatabaseId   string `json:"DatabaseId" gorm:"column:database_id;type:varchar(256);not null"`
	DatabaseName string `json:"DatabaseName" gorm:"column:database_name;type:varchar(256);not null"`
	OperType     int    `json:"OperType" gorm:"column:oper_type;type:tinyint;not null;default:1"` // 1-可见  2-查询 3-修改
	ValidTime    string `json:"ValidTime" gorm:"column:valid_time;type:varchar(256);not null"`
	OperCount    int    `json:"OperCount" gorm:"column:oper_count;type:int;not null;default:5"`
	AllowTables  string `json:"AllowTables" gorm:"column:allow_tables;type:varchar(2048);not null"`
	ApproveEmpId int    `json:"ApproveEmpId" gorm:"column:approve_emp_id;type:int"`
}

func (DmsAuth) TableName() string {
	return "dms_auth"
}

/*
	添加用户权限
*/
func AddUserDmsAuth(empId string, username string, authType int, instanceType string, instanceId string,
	instanceName string, databaseId string, databaseName string, operType int, validTime string,
	operCount int, allowTables string, approveEmpId string) (err error) {
	count := 0
	querySql := fmt.Sprintf("select count(*) from dms_auth where emp_id = ? and auth_type = ? " +
		"and instance_id = ? and database_id = ? and oper_type = ? and data_status = 1 and valid_time >= now() and oper_count > 0")
	err = database.Mysql().Raw(querySql, empId, authType, instanceId, databaseId, operType).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		err = errors.New("该用户存在此表对应的权限且在生效中！")
		return
	}
	err = database.Mysql().Exec("insert into dms_auth(emp_id, username, auth_type, instance_type, "+
		"instance_id, instance_name, database_id, database_name, oper_type, valid_time, "+
		"oper_count, allow_tables, approve_emp_id) values (?,?,?,?,?,?,?,?,?,?,?,?,?)",
		empId, username, authType, instanceType, instanceId, instanceName, databaseId, databaseName,
		operType, validTime, operCount, allowTables, approveEmpId).Error
	return
}

func GetUserAuthCount(empId string, instanceId string, operType string) (total uint) {
	total = 0
	querySql := "select count(*) from dms_auth where data_status = 1"
	var args []interface{}
	if empId != "" {
		querySql = fmt.Sprintf("%s and emp_id = ?", querySql)
		args = append(args, empId)
	}
	if instanceId != "" {
		querySql = fmt.Sprintf("%s and instance_id = ?",
			querySql)
		args = append(args, instanceId)
	}
	if operType != "" {
		querySql = fmt.Sprintf("%s and oper_type = ? ", querySql)
		args = append(args, operType)
	}
	database.Mysql().Raw(querySql, args...).Count(&total)
	return
}

func GetUserAuthList(empId string, instanceId string, operType string, offset uint, limit uint) (authList []DmsAuth) {
	querySql := "select * from dms_auth where data_status = 1"
	var args []interface{}
	if empId != "" {
		querySql = fmt.Sprintf("%s and emp_id = ?", querySql)
		args = append(args, empId)
	}
	if instanceId != "" {
		querySql = fmt.Sprintf("%s and instance_id = ?",
			querySql)
		args = append(args, instanceId)
	}
	if operType != "" {
		querySql = fmt.Sprintf("%s and oper_type = ? ", querySql)
		args = append(args, operType)
	}
	querySql += " order by id desc limit ?, ? "
	args = append(args, offset, limit)
	database.Mysql().Raw(querySql, args...).Scan(&authList)
	return
}

/*
	禁用用户权限
*/
func DeleteUserDmsAuth(authId int) (err error) {
	err = database.Mysql().Exec("update dms_auth set data_status = 0 where id = ?", authId).Error
	return
}

/*
	校验用户是否有库查询权限
*/
func CheckUserDMSAuthDatabaseQuery(empId string, databaseId string, tableNames string) (hasAuth bool, err error) {
	var auth DmsAuth
	err = database.Mysql().Raw("select * from dms_auth where data_status = 1 and auth_type = 2 "+
		"and oper_type = 2 and emp_id = ? and database_id = ? and oper_count > 0 and valid_time >= now()", empId, databaseId).Scan(&auth).Error
	if err != nil {
		return false, err
	}
	if auth.DatabaseId == databaseId { // 查询到了权限数据
		if strings.TrimSpace(auth.AllowTables) == "" {
			return true, err
		}
		userUseTableList := strings.Split(tableNames, ";")
		for _, tableName := range userUseTableList {
			if strings.Contains(auth.AllowTables, strings.TrimSpace(tableName)) {
				continue
			} else {
				return false, err
			}
		}
		return true, err
	}
	return false, err
}

func CheckUserSqlExistJoin(sqlContent string) (err error) {
	exist := strings.Contains(strings.ToLower(strings.TrimSpace(sqlContent)), "join")
	if exist {
		err = errors.New("当前不支持join相关的更新操作!")
	}
	return
}

/*
	校验用户是否有库执行权限
*/
func CheckUserDMSAuthDatabaseModify(empId string, databaseId string, tableNames string) (hasAuth bool, err error) {
	var auth DmsAuth
	var count int
	err = database.Mysql().Raw("select count(*) from dms_auth where data_status = 1 and auth_type = 2 "+
		"and oper_type = 3 and emp_id = ? and database_id = ? and oper_count > 0 and valid_time >= now()", empId, databaseId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	err = database.Mysql().Raw("select * from dms_auth where data_status = 1 and auth_type = 2 "+
		"and oper_type = 3 and emp_id = ? and database_id = ? and oper_count > 0 and valid_time >= now()", empId, databaseId).
		Scan(&auth).Error
	if err != nil {
		return false, err
	}
	if auth.DatabaseId == databaseId { // 查询到了权限数据
		if strings.TrimSpace(auth.AllowTables) == "" {
			return true, err
		}
		userUseTableList := strings.Split(tableNames, ";")
		for _, tableName := range userUseTableList {
			if strings.Contains(auth.AllowTables, strings.TrimSpace(tableName)) {
				continue
			} else {
				return false, err
			}
		}
		return true, err
	}
	return false, err
}

func GetUserInstanceData(empId string) (data []DmsAuth, err error) {
	err = database.Mysql().Raw("select emp_id, instance_type, instance_id, instance_name "+
		" from dms_auth where data_status = 1 and emp_id = ? and oper_count > 0 "+
		"and valid_time >= now() group by emp_id, instance_type, instance_id, "+
		"instance_name", empId).Scan(&data).Error
	return
}

func GetUserDatabaseData(empId string, instanceId string) (data []DmsAuth, err error) {
	err = database.Mysql().Raw("select database_id, database_name from dms_auth "+
		"where data_status = 1 and emp_id = ? and instance_id = ? and oper_count > 0 and valid_time >= now() "+
		"group by database_id, database_name", empId, instanceId).Scan(&data).Error
	return
}

func DecreaseUserOperatorCount(empId string, operType int, databaseId string) (err error) {
	err = database.Mysql().Exec("update dms_auth set oper_count = oper_count-1 where data_status = 1 and "+
		"oper_type = ? and emp_id = ? and database_id = ?", operType, empId, databaseId).Error
	return
}

type ProcesureResult struct {
	SqlResult int `json:"sql_result"`
}

func DmsBigTableProcedureCheck(host string, databaseName string, sqlContent string) (result ProcesureResult, err error) {
	err = database.Mysql().Raw("call p_dms_bigtable(?, ?, ?)", host, databaseName, sqlContent).Scan(&result).Error
	return
}

func DmsDeleteUserAuth(id string) (err error) {
	err = database.Mysql().Exec("update dms_auth set data_status = 0 where id = ?", id).Error
	return
}

func DmsUserAuthInfo(operType string, empId string, databaseId string) (auth DmsAuth, err error) {
	oper := 3
	if operType == "select" {
		oper = 2
	}
	err = database.Mysql().Raw("select * from dms_auth where data_status = 1 and auth_type = 2 "+
		"and oper_type = ? and emp_id = ? and database_id = ? and oper_count > 0 and valid_time >= now()",
		oper, empId, databaseId).Scan(&auth).Error
	return
}
