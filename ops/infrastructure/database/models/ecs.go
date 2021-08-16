package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chujieyang/commonops/ops/opslog"
	ecs2 "github.com/chujieyang/commonops/ops/value_objects/ecs"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/jinzhu/gorm"
)

type Ecs struct {
	gorm.Model
	DataStatus              int8   `json:"DataStatus" gorm:"not null;default:1"`
	CloudAccountId          uint   `json:"CloudAccountId"`
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
	SshPort                 int    `json:"SshPort"`
	SshUser                 string `json:"SshUser"`
	SshPwd                  string `json:"SshPwd"`
	Tags                    string `json:"Tags"`
}

type DiffCache struct {
	Type       string
	InstanceId string
}

func (Ecs) TableName() string {
	return "ecs"
}

func SaveOrUpdateEcs(instanceId string, cloudAccountId uint, ecs Ecs) bool {
	count := 0
	ecs.CloudAccountId = cloudAccountId
	result := database.Mysql().Model(&ecs).Where("instance_id=?", instanceId).Count(&count)
	if count > 0 {
		result = database.Mysql().Model(&ecs).Where("instance_id=?", instanceId).Updates(&ecs)
	} else {
		result = database.Mysql().Create(&ecs)
	}
	if result.Error != nil {
		opslog.Error().Println(result.Error)
		return false
	}
	AddEcsDiffCache(instanceId)
	return true
}

func AddEcsDiffCache(instanceId string) {
	database.Mysql().Create(&DiffCache{
		Type:       "ecs",
		InstanceId: instanceId,
	})
}

func GetCloudAccountKeyAndSecretByEcsInstanceId(instanceId string) (cloudAccount CloudAccount) {
	server := GetServerDetailByInstanceId(instanceId)
	database.Mysql().Model(&cloudAccount).Where("id=?", server.CloudAccountId).Find(&cloudAccount)
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
			"or instance_name like ? or instance_id like ? or tags like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	database.Mysql().Raw(querySql, args...).Count(&total)
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
			"or instance_name like ? or instance_id like ? or tags like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	querySql += " group by ecs.id order by id desc limit ?, ? "
	args = append(args, offset, limit)
	database.Mysql().Raw(querySql, args...).Scan(&servers)
	return
}

func GetServerDetail(id uint) (server EcsDetail) {
	database.Mysql().Raw("select ecs.*, cloud_account.name as cloud_account_name "+
		"from ecs left join cloud_account on ecs.cloud_account_id = cloud_account.id where ecs.id=?", id).Scan(&server)
	return
}

func GetServerDetailByInstanceId(instanceId string) (server Ecs) {
	database.Mysql().Model(&server).Where("instance_id=?", instanceId).Find(&server)
	return
}

func GetEcsListByRoleId(roleId int) (ecs []Ecs, err error) {
	err = database.Mysql().Raw("SELECT ur.resource_id as id FROM role_resources as ur "+
		"WHERE ur.resource_type = 'ecs' and ur.role_id = ?", roleId).Scan(&ecs).Error
	return
}

func GetAllEcsList() (ecs []Ecs, err error) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = database.Mysql().Raw("SELECT id, instance_name, public_ip_address, inner_ip_address"+
		" FROM ecs where data_status > 0 and (expired_time = '' or expired_time >= ?)", nowTime).Scan(&ecs).Error
	return
}

func AddRoleEcs(roleId int64, ecsIdList []string) (err error) {
	resIds := strings.Join(ecsIdList, ",")
	err = database.Mysql().Exec("delete from role_resources where role_id = ? "+
		"and resource_type='ecs' and resource_id not in (?)", roleId, resIds).Error
	for _, resId := range ecsIdList {
		exist := 0
		database.Mysql().Raw("select count(*) from role_resources "+
			"where resource_type='ecs' and role_id = ? and resource_id = ?", roleId, resId).Count(&exist)
		if exist == 0 {
			err = database.Mysql().Exec("insert into role_resources(role_id, resource_type, "+
				"resource_id) values(?, ?, ?)", roleId, "ecs", resId).Error
		}
	}
	return
}

func AddCloudServer(data ecs2.ServerInfoForm) (err error) {
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
		SshPort:              data.SshPort,
		SshUser:              data.SshUser,
		SshPwd:               utils.DesEncode(data.SshPwd),
		Tags:                 strings.Join(data.Tags, ","),
	}
	err = database.Mysql().Create(&ecs).Error
	return
}

func UpdateCloudServer(data ecs2.ExtraInfoForm) (err error) {
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

	if data.EcsBaseInfoForm.SshPort != "" {
		updateSql = fmt.Sprintf("%s, ssh_port='%s'", updateSql, data.EcsBaseInfoForm.SshPort)
	}

	if data.EcsBaseInfoForm.SshUser != "" {
		updateSql = fmt.Sprintf("%s, ssh_user='%s'", updateSql, data.EcsBaseInfoForm.SshUser)
	}

	if data.EcsBaseInfoForm.SshPwd != "" {
		sshInfo, _ := GetEcsSshInfo(data.Id)
		if sshInfo.SshPwd != data.EcsBaseInfoForm.SshPwd {
			updateSql = fmt.Sprintf("%s, ssh_pwd='%s'", updateSql, utils.DesEncode(data.EcsBaseInfoForm.SshPwd))
		}
	}

	updateSql = fmt.Sprintf("%s, tags='%s'", updateSql, strings.Join(data.EcsBaseInfoForm.Tags, ","))

	updateSql = fmt.Sprintf("%s where id in (?)", updateSql)
	err = database.Mysql().Exec(updateSql, strings.Split(data.Id, ",")).Error
	return
}

func DeleteCloudServer(id uint) (err error) {
	err = database.Mysql().Exec("update ecs set data_status = 0 where id = ?", id).Error
	return
}

func GetEcsRdsKvSlbCount() (resp map[string]interface{}) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	var ecsCount int
	var rdsCount int
	var kvCount int
	var slbCount int
	database.Mysql().Raw("select count(*) from ecs where data_status > 0 "+
		"and ( expired_time = '' or expired_time >= ? )", nowTime).Count(&ecsCount)
	database.Mysql().Raw("select count(*) from rds where data_status > 0 "+
		"and ( expire_time = '' or expire_time >= ? )", nowTime).Count(&rdsCount)
	database.Mysql().Raw("select count(*) from kv where data_status > 0 "+
		"and ( end_time = '' or end_time >= ? )", nowTime).Count(&kvCount)
	database.Mysql().Raw("select count(*) from slb where data_status > 0 ").Count(&slbCount)
	resp = map[string]interface{}{
		"ecsCount": ecsCount,
		"rdsCount": rdsCount,
		"kvCount":  kvCount,
		"slbCount": slbCount,
	}
	return
}

func GetEcsSshInfo(id string) (ecs Ecs, err error) {
	err = database.Mysql().Raw("select public_ip_address, ssh_port, ssh_user, ssh_pwd from ecs where data_status > 0 and id = ?", id).Scan(&ecs).Error
	return
}

func CleanEcsDiffCaches() error {
	return database.Mysql().Exec("delete from diff_caches where type = 'ecs' ").Error
}

func GetEcsTreeData() (ecs []Ecs, err error) {
	err = database.Mysql().Raw("SELECT id, instance_name, description, public_ip_address, inner_ip_address, tags from ecs where data_status = 2").Scan(&ecs).Error
	return
}