package models

import (
	"fmt"
	other2 "github.com/chujieyang/commonops/ops/value_objects/other"
	"strings"
	"time"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type OtherRes struct {
	gorm.Model
	DataStatus     int8   `json:"DataStatus" gorm:"not null;default:1"`
	CloudAccountId int    `json:"CloudAccountId"`
	ResType        string `json:"ResType"`
	InstanceId     string `json:"InstanceId"`
	InstanceName   string `json:"InstanceName"`
	Connections    string `json:"Connections"`
	Region         string `json:"Region"`
	Engine         string `json:"Engine"`
	Cpu            string `json:"Cpu"`
	Disk           string `json:"Disk"`
	Bandwidth      string `json:"Bandwidth"`
	Memory         string `json:"Memory"`
	RenewStatus    int    `json:"RenewStatus" gorm:"not null;default:1"`
	CreateTime     string `json:"CreateTime"`
	ExpiredTime    string `json:"ExpiredTime"`
}

func (OtherRes) TableName() string {
	return "other_res"
}

func GetOtherResCount(userId uint, queryKeyword string, queryResType string,
	queryCloudAccount uint, queryExpiredTime string, queryManageUser uint) (total uint) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	total = 0
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps {
		querySql = "select count(distinct other_res.id) from other_res where other_res.data_status > 0 "
	} else {
		querySql = "select count(distinct other_res.id) from other_res inner join role_resources as rr " +
			"on other_res.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"where other_res.data_status > 0 and rr.resource_type = 'other' and ur.user_id = ? "
		args = append(args, userId)
	}
	if queryResType != "所有" {
		querySql = fmt.Sprintf("%s and other_res.res_type = ? ", querySql)
		args = append(args, queryResType)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (instance_name like ? "+
			"or connections like ? or instance_id like ?) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	if queryExpiredTime != "" {
		querySql = fmt.Sprintf("%s and ( expired_time = '' or expired_time <= ? ) ", querySql)
		args = append(args, queryExpiredTime)
	} else {
		querySql = fmt.Sprintf("%s and ( expired_time = '' or expired_time >= ? ) ", querySql)
		args = append(args, nowTime)
	}
	if queryManageUser != 0 {
		querySql = fmt.Sprintf("%s and manage_user = ? ", querySql)
		args = append(args, queryManageUser)
	}
	database.Mysql().Raw(querySql, args...).Count(&total)
	return
}

func GetOtherRes(userId uint, queryKeyword string, queryResType string, queryCloudAccount uint, queryExpiredTime string,
	queryManageUser uint,offset uint, limit uint) (otherRes []OtherDetail) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps {
		querySql = "select other_res.*, ca.name as cloud_account_name from other_res left join cloud_account as ca " +
			"on other_res.cloud_account_id = ca.id where other_res.data_status > 0 " +
			"and ( expired_time = '' or expired_time >= ? ) "
		args = append(args, nowTime)
	} else {
		querySql = "select other_res.*, ca.name as cloud_account_name from other_res inner join role_resources as rr " +
			"on other_res.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"left join cloud_account as ca on other_res.cloud_account_id = ca.id" +
			"where other_res.data_status > 0 and rr.resource_type = 'other' and ur.user_id = ? " +
			"and ( expired_time = '' or expired_time >= ? ) "
		args = append(args, userId, nowTime)
	}
	if queryResType != "所有" {
		querySql = fmt.Sprintf("%s and other_res.res_type = ? ", querySql)
		args = append(args, queryResType)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (instance_name like ? "+
			"or connections like ? or instance_id like ?) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryExpiredTime != "" {
		querySql = fmt.Sprintf("%s and ( expired_time = '' or expired_time <= ? ) ", querySql)
		args = append(args, queryExpiredTime)
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	if queryManageUser != 0 {
		querySql = fmt.Sprintf("%s and manage_user = ? ", querySql)
		args = append(args, queryManageUser)
	}
	querySql += " group by other_res.id order by other_res.id desc limit ?, ? "
	args = append(args, offset, limit)
	database.Mysql().Raw(querySql, args...).Scan(&otherRes)
	return
}

func AddCloudOtherRes(data other2.AddOtherResForm) (err error) {
	if data.Cost == "" {
		data.Cost = "0.00"
	}
	if data.Id != 0 { // update
		err = database.Mysql().Exec("update other_res set cloud_account_id = ?, res_type = ?, "+
			"instance_id = ?, instance_name = ?, connections = ?, region = ?, engine = ?, cpu = ?, disk = ?, bandwidth = ?, memory = ?, "+
			"create_time = ?, expired_time = ? where id = ?", data.CloudAccountId, data.ResType,
			data.InstanceId, data.InstanceName, data.Connections, data.Region, data.Engine, data.Cpu, data.Disk,
			data.Bandwidth, data.Memory, data.CreateTime, data.ExpiredTime, data.Id).Error
	} else { // create
		err = database.Mysql().Exec("insert into other_res (data_status, cloud_account_id, res_type, "+
			"instance_id, instance_name, connections, region, engine, cpu, disk, bandwidth, memory, create_time, "+
			"expired_time) value (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", 1, data.CloudAccountId, data.ResType,
			data.InstanceId, data.InstanceName, data.Connections, data.Region, data.Engine, data.Cpu, data.Disk,
			data.Bandwidth, data.Memory, data.CreateTime, data.ExpiredTime).Error
	}
	return
}

func GetOtherResListByRoleId(roleId int) (otherRes []OtherRes, err error) {
	err = database.Mysql().Raw("SELECT ur.resource_id as id FROM role_resources as ur "+
		"WHERE ur.resource_type = 'other' and ur.role_id = ?", roleId).Scan(&otherRes).Error
	return
}

func GetAllOtherResList() (otherRes []OtherRes, err error) {
	err = database.Mysql().Raw("SELECT id, res_type, instance_name " +
		" FROM other_res where data_status > 0").Scan(&otherRes).Error
	return
}

func AddRoleOtherRes(roleId int64, otherIdList []string) (err error) {
	resIds := strings.Join(otherIdList, ",")
	err = database.Mysql().Exec("delete from role_resources where role_id = ? "+
		"and resource_type='other' and resource_id not in (?)", roleId, resIds).Error
	for _, resId := range otherIdList {
		exist := 0
		database.Mysql().Raw("select count(*) from role_resources "+
			"where resource_type='other' and role_id = ? and resource_id = ?", roleId, resId).Count(&exist)
		if exist == 0 {
			err = database.Mysql().Exec("insert into role_resources(role_id, resource_type, "+
				"resource_id) values(?, ?, ?)", roleId, "other", resId).Error
		}
	}
	return
}

func DeleteCloudOtherRes(id uint) (err error) {
	err = database.Mysql().Exec("update other_res set data_status = 0 where id = ?", id).Error
	return err
}
