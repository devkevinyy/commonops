package models

import (
	"encoding/json"
	"fmt"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type DmsLog struct {
	gorm.Model
	DataStatus int8    `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	EmpId string `json:"EmpId" gorm:"column:emp_id;type:varchar(32)"`
	Username string `json:"Username" gorm:"column:username;type:varchar(32)"`
	DatabaseId  string    `json:"DatabaseId" gorm:"column:database_id;type:varchar(128)"`
	DatabaseName string    `json:"DatabaseName"  gorm:"column:database_name;type:varchar(32)"`
	StartTime  string    `json:"StartTime" gorm:"column:start_time;type:varchar(32)"`
	SqlContent  string    `json:"SqlContent" gorm:"column:sql_content;type:text"`
	ExecStatus int    `json:"ExecStatus" gorm:"column:exec_status;type:tinyint"`
	Duration  int    `json:"Duration" gorm:"column:duration;type:int"`
	EffectRows  int    `json:"EffectRows" gorm:"column:effect_rows;type:int"`
	Result string    `json:"Result" gorm:"column:result;type:text"`
	ExceptionOutput string `json:"ExceptionOutput" gorm:"column:exception_output;type:text"`
	HasExecuted int `json:"HasExecuted" gorm:"column:has_executed;type:tinyint;not null;default:0"` // 是否已经执行 0-未执行 1-已执行
	RollbackTableName string  `json:"RollbackTableName" gorm:"column:rollback_table_name;type:varchar(1024)"`
	HasRollback int `json:"HasRollback" gorm:"column:has_rollback;type:tinyint;not null;default:0"` // 是否已经执行 0-未回滚 1-已回滚
	RollbackTime string `json:"RollbackTime" gorm:"column:rollback_time;type:varchar(32);not null;"`
	SqlType string `json:"SqlType" gorm:"column:sql_type;type:varchar(32);not null;"`
}

func (DmsLog) TableName() string {
	return "dms_log"
}

func SaveDmsQueryLog(empId string, username string, databaseId string, databaseName string, startTime string,
	sql string, status int, duration int, effectRows int, result []map[string]string, exceptionOutput string, hasExecuted int) (err error) {
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return
	}
	saveSql := fmt.Sprintf("INSERT INTO dms_log(emp_id, username, database_id, database_name, " +
		"start_time, sql_content, exec_status, duration, effect_rows, result, exception_output, has_executed) " +
		"VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
	err = database.Mysql().Exec(saveSql, empId, username, databaseId, databaseName, startTime,
		sql, status, duration, effectRows, string(resultBytes), exceptionOutput, hasExecuted).Error
	return
}

func UpdateDmsQueryLog(id int, status int, duration int, effectRows int, result []map[string]string, exceptionOutput string, rollbackTableName string) (err error) {
	resultContent := "******"
	saveSql := fmt.Sprintf("update dms_log set exec_status = ?, duration = ?, effect_rows = ?, " +
		"result = ?, exception_output = ?, rollback_table_name = ?, has_executed = 1 where id = ?")
	err = database.Mysql().Exec(saveSql, status, duration, effectRows, resultContent, exceptionOutput, rollbackTableName, id).Error
	return
}

func GetDmsQueryLogCount(empId string) (count uint) {
	querySql := "select count(*) from dms_log where data_status = 1 and has_executed = 1 and emp_id = ?"
	database.Mysql().Raw(querySql, empId).Count(&count)
	return
}

func GetDmsQueryLogByPage(empId string, offset uint, limit uint) (log []DmsLog) {
	querySql := "select * from dms_log where data_status = 1 and has_executed = 1 and emp_id = ? order by id desc limit ?, ?"
	database.Mysql().Raw(querySql, empId, offset, limit).Scan(&log)
	return
}

func CheckDmsSqlLogRollbackAuth(empId string, logId string) (hasAuth bool, err error) {
	var count int
	querySql := "select count(*) from dms_log where data_status = 1 and has_executed = 1 " +
		"and exec_status = 1 and emp_id = ? and id = ?"
	err = database.Mysql().Raw(querySql, empId, logId).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		hasAuth = true
	}
	return
}

func GetDmsLogById(id string) (log DmsLog, err error) {
	querySql := "select * from dms_log where id = ? limit 1"
	err = database.Mysql().Raw(querySql, id).Scan(&log).Error
	return
}

func UpdateDmsRollbackLogResult(id string, hasRollback int, rollbackTime string) (err error) {
	updateSql := "update dms_log set has_rollback = ?, rollback_time = ? where id = ?"
	err = database.Mysql().Exec(updateSql, hasRollback, rollbackTime, id).Error
	return
}

func GetLatestSqlList(empId string, sqlType string) (logs []DmsLog, err error) {
	querySql := "select sql_content from dms_log where data_status = 1 and has_executed = 1 " +
		"and sql_type = ? and emp_id = ? and start_time >= DATE_SUB(NOW(), INTERVAL 8 HOUR)"
	err = database.Mysql().Raw(querySql, sqlType, empId).Scan(&logs).Error
	return
}