package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/forms"
	"github.com/chujieyang/commonops/ops/untils"
	"github.com/jinzhu/gorm"
)

type Ecs struct {
	gorm.Model
	DataStatus              int8   `json:"DataStatus" gorm:"not null;default:1"`
	CloudAccountId          int    `json:"CloudAccountId"`
	ImageId                 string `json:"ImageId"`
	InstanceType            string `json:"InstanceType"`
	InstanceNetworkType     string `json:"InstanceNetworkType"`
	LocalStorageAmount      string `json:"LocalStorageAmount"`
	InstanceChargeType      string `json:"InstanceChargeType"`
	ClusterId               string `json:"ClusterId"`
	InstanceName            string `json:"InstanceName"`
	StartTime               string `json:"StartTime"`
	ZoneId                  string `json:"ZoneId"`
	InternetChargeType      string `json:"InternetChargeType"`
	InternetMaxBandwidthIn  string `json:"InternetMaxBandwidthIn"`
	HostName                string `json:"HostName"`
	Cpu                     string `json:"Cpu"`
	Status                  string `json:"Status"`
	OSName                  string `json:"OSName"`
	OSNameEn                string `json:"OSNameEn"`
	SerialNumber            string `json:"SerialNumber"`
	RegionId                string `json:"RegionId"`
	InternetMaxBandwidthOut string `json:"InternetMaxBandwidthOut"`
	ResourceGroupId         string `json:"ResourceGroupId"`
	InstanceTypeFamily      string `json:"InstanceTypeFamily"`
	InstanceId              string `json:"InstanceId"`
	DeploymentSetId         string `json:"DeploymentSetId"`
	Description             string `json:"Description"`
	ExpiredTime             string `json:"ExpiredTime"`
	OSType                  string `json:"OSType"`
	Memory                  string `json:"Memory"`
	CreationTime            string `json:"CreationTime"`
	LocalStorageCapacity    string `json:"LocalStorageCapacity"`
	InnerIpAddress          string `json:"InnerIpAddress"`
	PublicIpAddress         string `json:"PublicIpAddress"`
	PrivateIpAddress        string `json:"PrivateIpAddress"`
}

type DiffCache struct {
	Type       string
	InstanceId string
}

func (Ecs) TableName() string {
	return "ecs"
}

func SaveOrUpdateEcs(instanceId string, cloudAccountId int, ecs Ecs) bool {
	count := 0
	ecs.CloudAccountId = cloudAccountId
	result := database.MysqlClient.Model(&ecs).Where("instance_id=?", instanceId).Count(&count)
	if count > 0 {
		result = database.MysqlClient.Model(&ecs).Where("instance_id=?", instanceId).Omit("renew_site_id",
			"manage_user", "cost", "cost_type", "bank_account", "renew_status").Updates(&ecs)
	} else {
		result = database.MysqlClient.Create(&ecs)
	}
	if result.Error != nil {
		untils.Log.Error(result.Error.Error())
		return false
	}
	AddEcsDiffCache(instanceId)
	return true
}

func AddEcsDiffCache(instanceId string) {
	database.MysqlClient.Create(&DiffCache{
		Type:       "ecs",
		InstanceId: instanceId,
	})
}

func GetCloudAccountKeyAndSecretByEcsInstanceId(instanceId string) (cloudAccount CloudAccount) {
	server := GetServerDetailByInstanceId(instanceId)
	database.MysqlClient.Model(&cloudAccount).Where("id=?", server.CloudAccountId).Find(&cloudAccount)
	return
}

func GetServersCount(userId uint, queryExpiredTime string, queryKeyword string, queryCloudAccount int) (total int) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	total = 0
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps { // 是超级管理员或运维中心角色组不参与资源权限控制
		querySql = "select count(distinct ecs.id) from ecs where ecs.data_status > 0 "
	} else {
		querySql = "select count(distinct ecs.id) from ecs inner join role_resources as rr " +
			"on ecs.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"where ecs.data_status > 0 and rr.resource_type = 'ecs' and ur.user_id = ? "
		args = append(args, userId)
	}
	if queryExpiredTime != "" {
		querySql = fmt.Sprintf("%s and ( expired_time = '' or expired_time <= ? ) ", querySql)
		args = append(args, queryExpiredTime)
	} else {
		querySql = fmt.Sprintf("%s and ( expired_time = '' or expired_time >= ? ) ", querySql)
		args = append(args, nowTime)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (inner_ip_address like ? "+
			"or public_ip_address like ? or private_ip_address like ? "+
			"or instance_name like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	database.MysqlClient.Raw(querySql, args...).Count(&total)
	return
}

func GetServers(userId uint, offset int, limit int, queryExpiredTime string, queryKeyword string,
	queryCloudAccount int) (servers []EcsDetail) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps { // 是超级管理员或运维中心角色组不参与资源权限控制
		querySql = "select ecs.*, ca.name as cloud_account_name from ecs left join cloud_account as ca on ecs.cloud_account_id = ca.id " +
			"where ecs.data_status > 0 "
	} else {
		querySql = "select ecs.*, ca.name as cloud_account_name from ecs inner join role_resources as rr " +
			"on ecs.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"left join cloud_account as ca on ecs.cloud_account_id = ca.id " +
			"where ecs.data_status > 0 and rr.resource_type = 'ecs' and ur.user_id = ? "
		args = append(args, userId)
	}
	if queryExpiredTime != "" {
		querySql = fmt.Sprintf("%s and ( expired_time = '' or expired_time <= ? ) ", querySql)
		args = append(args, queryExpiredTime)
	} else {
		querySql = fmt.Sprintf("%s and ( expired_time = '' or expired_time >= ? ) ", querySql)
		args = append(args, nowTime)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (inner_ip_address like ? "+
			"or public_ip_address like ? or private_ip_address like ? "+
			"or instance_name like ? or instance_id like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	querySql += " group by ecs.id order by id desc limit ?, ? "
	args = append(args, offset, limit)
	database.MysqlClient.Raw(querySql, args...).Scan(&servers)
	return
}

func GetServerDetail(id int) (server EcsDetail) {
	database.MysqlClient.Raw("select ecs.*, cloud_account.name as cloud_account_name "+
		"from ecs left join cloud_account on ecs.cloud_account_id = cloud_account.id where ecs.id=?", id).Scan(&server)
	return
}

func GetServerDetailByInstanceId(instanceId string) (server Ecs) {
	database.MysqlClient.Model(&server).Where("instance_id=?", instanceId).Find(&server)
	return
}

func GetEcsListByRoleId(roleId int) (ecs []Ecs, err error) {
	err = database.MysqlClient.Raw("SELECT ur.resource_id as id FROM role_resources as ur "+
		"WHERE ur.resource_type = 'ecs' and ur.role_id = ?", roleId).Scan(&ecs).Error
	return
}

func GetAllEcsList() (ecs []Ecs, err error) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = database.MysqlClient.Raw("SELECT id, instance_name, public_ip_address, inner_ip_address"+
		" FROM ecs where data_status > 0 and (expired_time = '' or expired_time >= ?)", nowTime).Scan(&ecs).Error
	return
}

func AddRoleEcs(roleId int64, ecsIdList []string) (err error) {
	resIds := strings.Join(ecsIdList, ",")
	err = database.MysqlClient.Exec("delete from role_resources where role_id = ? "+
		"and resource_type='ecs' and resource_id not in (?)", roleId, resIds).Error
	for _, resId := range ecsIdList {
		exist := 0
		database.MysqlClient.Raw("select count(*) from role_resources "+
			"where resource_type='ecs' and role_id = ? and resource_id = ?", roleId, resId).Count(&exist)
		if exist == 0 {
			err = database.MysqlClient.Exec("insert into role_resources(role_id, resource_type, "+
				"resource_id) values(?, ?, ?)", roleId, "ecs", resId).Error
		}
	}
	return
}

func AddCloudServer(data forms.ServerInfoForm) (err error) {
	ecs := Ecs{
		HostName:             data.HostName,
		InstanceId:           data.InstanceId,
		InstanceName:         data.InstanceName,
		Description:          data.Description,
		InnerIpAddress:       data.InnerIpAddress,
		PublicIpAddress:      data.PublicIpAddress,
		Memory:               strconv.Itoa(data.Memory * 1024),
		LocalStorageCapacity: strconv.Itoa(data.Disk * 1024),
		CreationTime:         data.CreateTime,
		ExpiredTime:          data.ExpiredTime,
		CloudAccountId:       data.CloudAccountId,
		OSType:               data.OsType,
		Cpu:                  strconv.Itoa(data.Cpu),
		Status:               "Running",
		DataStatus:           2, // 表示用户自定义添加
	}
	err = database.MysqlClient.Create(&ecs).Error
	return
}

func UpdateCloudServer(data forms.ExtraInfoForm) (err error) {
	updateSql := "update ecs set id=id"
	if data.EcsBaseInfoForm.Memory != "" {
		memoryInt, err1 := strconv.Atoi(data.EcsBaseInfoForm.Memory)
		if err1 != nil {
			err = err1
			return
		}
		memory := fmt.Sprintf("%d", memoryInt*1024)
		updateSql = fmt.Sprintf("%s, memory='%s'", updateSql, memory)
	}

	if data.EcsBaseInfoForm.Cpu != "" {
		cpu := fmt.Sprintf("%s", data.EcsBaseInfoForm.Cpu)
		updateSql = fmt.Sprintf("%s, cpu='%s'", updateSql, cpu)
	}

	if data.EcsBaseInfoForm.InstanceId != "" {
		instanceId := fmt.Sprintf("%s", data.EcsBaseInfoForm.InstanceId)
		updateSql = fmt.Sprintf("%s, instance_id='%s'", updateSql, instanceId)
	}

	if data.EcsBaseInfoForm.InstanceName != "" {
		instanceName := fmt.Sprintf("%s", data.EcsBaseInfoForm.InstanceName)
		updateSql = fmt.Sprintf("%s, instance_name='%s'", updateSql, instanceName)
	}

	if data.EcsBaseInfoForm.InnerIpAddress != "" {
		innerIpAddress := fmt.Sprintf("%s", data.EcsBaseInfoForm.InnerIpAddress)
		updateSql = fmt.Sprintf("%s, inner_ip_address='%s'", updateSql, innerIpAddress)
	}

	if data.EcsBaseInfoForm.PublicIpAddress != "" {
		publicIpAddress := fmt.Sprintf("%s", data.EcsBaseInfoForm.PublicIpAddress)
		updateSql = fmt.Sprintf("%s, public_ip_address='%s'", updateSql, publicIpAddress)
	}

	if data.EcsBaseInfoForm.PrivateIpAddress != "" {
		privateIpAddress := fmt.Sprintf("%s", data.EcsBaseInfoForm.PrivateIpAddress)
		updateSql = fmt.Sprintf("%s, private_ip_address='%s'", updateSql, privateIpAddress)
	}

	if data.EcsBaseInfoForm.ExpiredTime != "" {
		expireTime := fmt.Sprintf("%s", data.EcsBaseInfoForm.ExpiredTime)
		updateSql = fmt.Sprintf("%s, expired_time='%s'", updateSql, expireTime)
	}

	updateSql = fmt.Sprintf("%s where id in (?)", updateSql)
	err = database.MysqlClient.Debug().Exec(updateSql, strings.Split(data.Id, ",")).Error
	return
}

func DeleteCloudServer(id int) (err error) {
	err = database.MysqlClient.Exec("update ecs set data_status = 0 where id = ?", id).Error
	return
}

func GetEcsRdsKvSlbCount() (resp map[string]interface{}) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	var ecsCount int
	var rdsCount int
	var kvCount int
	var slbCount int
	database.MysqlClient.Raw("select count(*) from ecs where data_status > 0 "+
		"and ( expired_time = '' or expired_time >= ? )", nowTime).Count(&ecsCount)
	database.MysqlClient.Raw("select count(*) from rds where data_status > 0 "+
		"and ( expire_time = '' or expire_time >= ? )", nowTime).Count(&rdsCount)
	database.MysqlClient.Raw("select count(*) from kv where data_status > 0 "+
		"and ( end_time = '' or end_time >= ? )", nowTime).Count(&kvCount)
	database.MysqlClient.Raw("select count(*) from slb where data_status > 0 ").Count(&slbCount)
	resp = map[string]interface{}{
		"ecsCount": ecsCount,
		"rdsCount": rdsCount,
		"kvCount":  kvCount,
		"slbCount": slbCount,
	}
	return
}
