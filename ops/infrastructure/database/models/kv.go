package models

import (
	"fmt"
	kv2 "github.com/chujieyang/commonops/ops/value_objects/kv"
	"strconv"
	"strings"
	"time"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type Kv struct {
	gorm.Model
	DataStatus          int8   `json:"DataStatus" gorm:"not null;default: 1"`
	CloudAccountId      uint    `json:"CloudAccountId"`
	InstanceClass       string `json:"InstanceClass" xml:"InstanceClass"`
	PackageType         string `json:"PackageType" xml:"PackageType"`
	ChargeType          string `json:"ChargeType" xml:"ChargeType"`
	ConnectionDomain    string `json:"ConnectionDomain" xml:"ConnectionDomain"`
	SearchKey           string `json:"SearchKey" xml:"SearchKey"`
	CreateTime          string `json:"CreateTime" xml:"CreateTime"`
	HasRenewChangeOrder string `json:"HasRenewChangeOrder" xml:"HasRenewChangeOrder"`
	InstanceType        string `json:"InstanceType" xml:"InstanceType"`
	DestroyTime         string `json:"DestroyTime" xml:"DestroyTime"`
	RegionId            string `json:"RegionId" xml:"RegionId"`
	PrivateIp           string `json:"PrivateIp" xml:"PrivateIp"`
	InstanceId          string `json:"InstanceId" xml:"InstanceId"`
	InstanceStatus      string `json:"InstanceStatus" xml:"InstanceStatus"`
	Bandwidth           int    `json:"Bandwidth" xml:"Bandwidth"`
	NetworkType         string `json:"NetworkType" xml:"NetworkType"`
	VpcId               string `json:"VpcId" xml:"VpcId"`
	NodeType            string `json:"NodeType" xml:"NodeType"`
	Connections         int    `json:"Connections" xml:"Connections"`
	ArchitectureType    string `json:"ArchitectureType" xml:"ArchitectureType"`
	ReplacateId         string `json:"ReplacateId" xml:"ReplacateId"`
	EngineVersion       string `json:"EngineVersion" xml:"EngineVersion"`
	Capacity            int    `json:"Capacity" xml:"Capacity"`
	VSwitchId           string `json:"VSwitchId" xml:"VSwitchId"`
	InstanceName        string `json:"InstanceName" xml:"InstanceName"`
	Port                int    `json:"Port" xml:"Port"`
	ZoneId              string `json:"ZoneId" xml:"ZoneId"`
	EndTime             string `json:"EndTime" xml:"EndTime"`
	QPS                 int    `json:"QPS" xml:"QPS"`
	UserName            string `json:"UserName" xml:"UserName"`
	Config              string `json:"Config" xml:"Config"`
	IsRds               bool   `json:"IsRds" xml:"IsRds"`
	ConnectionMode      string `json:"ConnectionMode" xml:"ConnectionMode"`
}

func (Kv) TableName() string {
	return "kv"
}

func GetCloudAccountKeyAndSecretByKvInstanceId(instanceId string) (cloudAccount CloudAccount) {
	server := GetKvDetailByInstanceId(instanceId)
	database.Mysql().Model(&cloudAccount).Where("id=?", server.CloudAccountId).Find(&cloudAccount)
	return
}

func GetKvDetailByInstanceId(instanceId string) (kv Kv) {
	database.Mysql().Model(&kv).Where("instance_id=?", instanceId).Find(&kv)
	return
}

func GetKvCount(userId uint, queryExpiredTime string, queryKeyword string,
	queryCloudAccount uint) (total uint) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	total = 0
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps {
		querySql = "select count(distinct kv.id) from kv where kv.data_status > 0 "
	} else {
		querySql = "select count(distinct kv.id) from kv inner join role_resources as rr " +
			"on kv.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"where kv.data_status > 0 and rr.resource_type = 'kv' and ur.user_id = ? "
		args = append(args, userId)
	}
	if queryExpiredTime != "" {
		querySql = fmt.Sprintf("%s and ( end_time = '' or end_time <= ? ) ", querySql)
		args = append(args, queryExpiredTime)
	} else {
		querySql = fmt.Sprintf("%s and ( end_time = '' or end_time >= ? ) ", querySql)
		args = append(args, nowTime)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (instance_name like ? or connection_domain like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	database.Mysql().Raw(querySql, args...).Count(&total)
	return
}

func GetKvByPage(userId uint, queryExpiredTime string, queryKeyword string,
	queryCloudAccount uint, offset uint, limit uint) (kv []KvDetail) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps {
		querySql = "select kv.*, ca.name as cloud_account_name from kv " +
			"left join cloud_account as ca on kv.cloud_account_id = ca.id where kv.data_status > 0 "
	} else {
		querySql = "select kv.*, ca.name as cloud_account_name from kv inner join role_resources as rr " +
			"on kv.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"left join cloud_account as ca on kv.cloud_account_id = ca.id " +
			"where kv.data_status > 0 and rr.resource_type = 'kv' and ur.user_id = ? "
		args = append(args, userId)
	}
	if queryExpiredTime != "" {
		querySql = fmt.Sprintf("%s and ( end_time = '' or end_time <= ? ) ", querySql)
		args = append(args, queryExpiredTime)
	} else {
		querySql = fmt.Sprintf("%s and ( end_time = '' or end_time >= ? ) ", querySql)
		args = append(args, nowTime)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (instance_name like ? "+
			"or connection_domain like ? or instance_id like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	querySql += " group by kv.id order by id desc limit ?, ? "
	args = append(args, offset, limit)
	database.Mysql().Raw(querySql, args...).Scan(&kv)
	return
}

func SaveOrUpdateKv(instanceId string, cloudAccountId uint, kv Kv) (err error) {
	count := 0
	kv.CloudAccountId = cloudAccountId
	database.Mysql().Model(&kv).Where("instance_id=?", instanceId).Count(&count)
	if count > 0 {
		err = database.Mysql().Model(&kv).Where("instance_id=?", instanceId).Omit("renew_site_id",
			"manage_user", "cost", "cost_type", "bank_account", "renew_status").Updates(&kv).Error
	} else {
		err = database.Mysql().Create(&kv).Error
	}
	AddKvDiffCache(instanceId)
	return
}

func AddKvDiffCache(instanceId string) {
	database.Mysql().Create(&DiffCache{
		Type:       "kv",
		InstanceId: instanceId,
	})
}

func GetKvDetail(id uint) (kv KvDetail) {
	database.Mysql().Raw("select kv.*, cloud_account.name as cloud_account_name "+
		"from kv left join cloud_account "+
		"on kv.cloud_account_id = cloud_account.id where kv.id=?", id).Scan(&kv)
	return
}

func GetKvListByRoleId(roleId int) (kv []Kv, err error) {
	err = database.Mysql().Raw("SELECT ur.resource_id as id FROM role_resources as ur "+
		"WHERE ur.resource_type = 'kv' and ur.role_id = ?", roleId).Scan(&kv).Error
	return
}

func GetAllKvList() (kv []Kv, err error) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = database.Mysql().Raw("SELECT id, instance_name, connection_domain "+
		" FROM kv where data_status > 0 and (end_time = '' or end_time >= ?)", nowTime).Scan(&kv).Error
	return
}

func AddRoleKv(roleId int64, kvIdList []string) (err error) {
	resIds := strings.Join(kvIdList, ",")
	err = database.Mysql().Exec("delete from role_resources where role_id = ? "+
		"and resource_type='kv' and resource_id not in (?)", roleId, resIds).Error
	for _, resId := range kvIdList {
		exist := 0
		database.Mysql().Raw("select count(*) from role_resources "+
			"where resource_type='kv' and role_id = ? and resource_id = ?", roleId, resId).Count(&exist)
		if exist == 0 {
			err = database.Mysql().Exec("insert into role_resources(role_id, resource_type, "+
				"resource_id) values(?, ?, ?)", roleId, "kv", resId).Error
		}
	}
	return
}

func AddCloudKv(data kv2.KvInfoForm) (err error) {
	port, _ := strconv.Atoi(data.Port)
	kv := Kv{
		CloudAccountId:   data.CloudAccountId,
		InstanceId:       data.InstanceId,
		InstanceName:     data.InstanceName,
		ConnectionDomain: data.ConnectionString,
		Port:             port,
		Capacity:         data.Capacity * 1024,
		InstanceType:     data.InstanceType,
		Bandwidth:        data.Bandwidth,
		CreateTime:       data.CreateTime,
		EndTime:          data.EndTime,
		InstanceStatus:   "Normal",
		DataStatus:       2, // 表示用户自定义添加
	}
	err = database.Mysql().Create(&kv).Error
	return
}

func DeleteCloudKv(id uint) (err error) {
	err = database.Mysql().Exec("update kv set data_status = 0 where id = ?", id).Error
	return
}

func UpdateCloudKv(data kv2.ExtraInfoForm) (err error) {
	updateSql := "update kv set id=id"
	if data.KvBaseInfoForm.KvCapacity != "" {
		memoryInt, err1 := strconv.Atoi(data.KvBaseInfoForm.KvCapacity)
		if err1 != nil {
			err = err1
			return
		}
		updateSql = fmt.Sprintf("%s, capacity='%d'", updateSql, memoryInt*1024)
	}

	if data.KvBaseInfoForm.KvInstanceName != "" {
		updateSql = fmt.Sprintf("%s, instance_name='%s'", updateSql, data.KvBaseInfoForm.KvInstanceName)
	}

	if data.KvBaseInfoForm.KvBandwidth != "" {
		updateSql = fmt.Sprintf("%s, bandwidth='%s'", updateSql, data.KvBaseInfoForm.KvBandwidth)
	}

	if data.KvBaseInfoForm.KvExpiredTime != "" {
		updateSql = fmt.Sprintf("%s, end_time='%s'", updateSql, data.KvBaseInfoForm.KvExpiredTime)
	}

	updateSql = fmt.Sprintf("%s where id in (?)", updateSql)
	err = database.Mysql().Exec(updateSql, strings.Split(data.Id, ",")).Error
	return
}

func CleanKvDiffCaches() error {
	return database.Mysql().Exec("delete from diff_caches where type = 'kv' ").Error
}
