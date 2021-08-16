package models

import (
	"fmt"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type DmsInstance struct {
	gorm.Model
	DataStatus    int8   `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	Port          int    `json:"Port" gorm:"column:port;type:smallint;not null"`
	InstanceType  string `json:"InstanceType" gorm:"column:instance_type;type:varchar(32);not null"`
	Host          string `json:"Host" gorm:"column:host;type:varchar(128);not null"`
	State         string `json:"State" gorm:"column:state;type:varchar(32);not null"`
	InstanceId    string `json:"InstanceId" gorm:"column:instance_id;type:varchar(128);not null;unique_index:instance_id_idx"`
	InstanceAlias string `json:"InstanceAlias" gorm:"column:instance_alias;type:varchar(128);not null"`
	OperUser      string `json:"OperUser" gorm:"column:oper_user;type:varchar(128);not null"`
	OperPwd       string `json:"OperPwd" gorm:"column:oper_pwd;type:varchar(256);not null"`
}

func (DmsInstance) TableName() string {
	return "dms_instance"
}

func SaveDmsInstance(port int, instanceType string, host string, state string, instanceId string,
	instanceAlias string, operUser string, operPwd string) (err error) {
	saveSql := fmt.Sprintf("INSERT INTO dms_instance(port, instance_type, host, state, instance_id, instance_alias, oper_user, oper_pwd) VALUES(?,?,?,?,?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE port = ?, instance_type = ?, host = ?, state = ?, instance_alias = ?, oper_user = ?, oper_pwd = ?")
	err = database.Mysql().Exec(saveSql, port, instanceType, host, state, instanceId, instanceAlias, operUser, operPwd, port,
		instanceType, host, state, instanceAlias, operUser, operPwd).Error
	return
}

func DeleteDmsInstance(instanceId string) (err error) {
	err = database.Mysql().Exec("update dms_instance set data_status = 0 where instance_id = ?", instanceId).Error
	return
}

func GetDmsInstanceCount(instanceName string) (total uint) {
	total = 0
	querySql := "select count(*) from dms_instance where data_status = 1"
	var args []interface{}
	if instanceName != "" {
		querySql = fmt.Sprintf("%s and instance_alias like ? ", querySql)
		args = append(args, "%"+instanceName+"%")
	}
	database.Mysql().Raw(querySql, args...).Count(&total)
	return
}

func GetDmsInstanceList(instanceName string, offset uint, limit uint) (instanceList []DmsInstance) {
	querySql := "select * from dms_instance where data_status = 1"
	var args []interface{}
	if instanceName != "" {
		querySql = fmt.Sprintf("%s and instance_alias like ? ", querySql)
		args = append(args, "%"+instanceName+"%")
	}
	querySql += " order by id desc limit ?, ? "
	args = append(args, offset, limit)
	database.Mysql().Raw(querySql, args...).Scan(&instanceList)
	return
}

func GetDmsAllInstanceData() (instances []DmsInstance, err error) {
	err = database.Mysql().Raw("select instance_id, instance_type, instance_alias, host, port from dms_instance").Scan(&instances).Error
	return
}

func GetDmsInstanceByInstanceId(instanceId string) (instance DmsInstance, err error) {
	err = database.Mysql().Raw("select * from dms_instance where instance_id = ?", instanceId).Scan(&instance).Error
	return
}

func GetDmsInstanceOperUserInfoByInstanceId(instanceId string) (instance DmsInstance, err error) {
	err = database.Mysql().Raw("select oper_user, oper_pwd from dms_instance where instance_id = ?", instanceId).Scan(&instance).Error
	return
}
