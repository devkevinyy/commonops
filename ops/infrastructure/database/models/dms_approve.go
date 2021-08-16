package models

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/jinzhu/gorm"
)

type DmsApprove struct {
	gorm.Model
	DataStatus     int    `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	EmpId          string `json:"EmpId" gorm:"column:emp_id;type:varchar(32);not null"`
	Username       string `json:"Username" gorm:"column:username;type:varchar(256);not null"`
	InstanceType   string `json:"InstanceType" gorm:"column:instance_type;type:varchar(256);not null"`
	InstanceId     string `json:"InstanceId" gorm:"column:instance_id;type:varchar(512);not null"`
	InstanceName   string `json:"InstanceName" gorm:"column:instance_name;type:varchar(512);not null"`
	DatabaseId     string `json:"DatabaseId" gorm:"column:database_id;type:varchar(256);not null"`
	DatabaseName   string `json:"DatabaseName" gorm:"column:database_name;type:varchar(256);not null"`
	SqlContent     string `json:"SqlContent" gorm:"column:sql_content;type:text;not null;"`
	SqlDescription string `json:"SqlDescription" gorm:"column:sql_description;type:varchar(2048)"`
	CreateTime     string `json:"CreateTime" gorm:"column:create_time;type:varchar(256);not null"`
	ApproveId      string `json:"ApproveId" gorm:"column:approve_id;type:varchar(256);not null;default:''"`
	ApproveContent string `json:"ApproveContent" gorm:"column:approve_content;type:varchar(128);not null;default:'审批中'"`
	ApproveStatus  int    `json:"ApproveStatus" gorm:"column:approve_status;type:tinyint;not null;default:0"` // -1-审批拒绝 0-审批中  1-审批通过
	LogId          int    `json:"LogId" gorm:"column:log_id;type:int;not null"`
	HasExecuted    int    `json:"HasExecuted" gorm:"column:has_executed;type:tinyint;not null;default:0"` // 是否已经执行 0-未执行 1-已执行
}

func (DmsApprove) TableName() string {
	return "dms_approve"
}

func SaveDmsSqlApprove(emdId string, username string, instanceType string, instanceId string, instanceName string,
	databaseId string, databaseName string, sqlContent string, approveId string, logId int) (err error) {
	createTime := utils.GetCurrentTime()
	insertSql := "insert into dms_approve(emp_id, username, instance_type, instance_id, instance_name, database_id, " +
		"database_name, sql_content, create_time, approve_id, log_id) values (?,?,?,?,?,?,?,?,?,?,?,?)"
	err = database.Mysql().Exec(insertSql, emdId, username, instanceType, instanceId, instanceName, databaseId,
		databaseName, sqlContent, createTime, approveId, logId).Error
	return
}

func GetDmsUserApproveCount(empId string) (count int, err error) {
	querySql := "select count(*) from dms_approve where data_status = 1 and emp_id = ?"
	err = database.Mysql().Raw(querySql, empId).Count(&count).Error
	return
}

func GetDmsUserApproveByPage(empId string, offset int, limit int) (approves []DmsApprove, err error) {
	querySql := "select * from dms_approve where data_status = 1 and emp_id = ? order by id desc limit ?, ?"
	err = database.Mysql().Raw(querySql, empId, offset, limit).Scan(&approves).Error
	return
}

func GetApproveInfoById(id string) (approve DmsApprove, err error) {
	querySql := "select * from dms_approve where id = ?"
	err = database.Mysql().Raw(querySql, id).Scan(&approve).Error
	return
}

func SetApproveHasExecuted(id string) (err error) {
	err = database.Mysql().Exec("update dms_approve set has_executed = 1 where id = ?", id).Error
	return
}

func IsApproveCanExecute(id string) (canExecute bool, err error) {
	count := 0
	err = database.Mysql().Raw("select count(*) from dms_approve where id = ? and data_status = 1 "+
		"and approve_status = 1 and has_executed = 0", id).Count(&count).Error
	if count > 0 {
		canExecute = true
	}
	return
}
