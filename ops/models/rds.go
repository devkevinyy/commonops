package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/forms"
	"github.com/jinzhu/gorm"
)

type Rds struct {
	gorm.Model
	DataStatus            int8   `json:"DataStatus" gorm:"not null;default: 1"`
	CloudAccountId        int    `json:"CloudAccountId"`
	InsId                 int    `json:"InsId" xml:"InsId"`
	DBInstanceId          string `json:"DBInstanceId" xml:"DBInstanceId"`
	DBInstanceDescription string `json:"DBInstanceDescription" xml:"DBInstanceDescription"`
	PayType               string `json:"PayType" xml:"PayType"`
	DBInstanceType        string `json:"DBInstanceType" xml:"DBInstanceType"`
	RegionId              string `json:"RegionId" xml:"RegionId"`
	ExpireTime            string `json:"ExpireTime" xml:"ExpireTime"`
	DestroyTime           string `json:"DestroyTime" xml:"DestroyTime"`
	DBInstanceStatus      string `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
	Engine                string `json:"Engine" xml:"Engine"`
	DBInstanceNetType     string `json:"DBInstanceNetType" xml:"DBInstanceNetType"`
	ConnectionMode        string `json:"ConnectionMode" xml:"ConnectionMode"`
	LockMode              string `json:"LockMode" xml:"LockMode"`
	Category              string `json:"Category" xml:"Category"`
	DBInstanceStorageType string `json:"DBInstanceStorageType" xml:"DBInstanceStorageType"`
	DBInstanceClass       string `json:"DBInstanceClass" xml:"DBInstanceClass"`
	InstanceNetworkType   string `json:"InstanceNetworkType" xml:"InstanceNetworkType"`
	VpcCloudInstanceId    string `json:"VpcCloudInstanceId" xml:"VpcCloudInstanceId"`
	LockReason            string `json:"LockReason" xml:"LockReason"`
	ZoneId                string `json:"ZoneId" xml:"ZoneId"`
	MutriORsignle         bool   `json:"MutriORsignle" xml:"MutriORsignle"`
	CreateTime            string `json:"CreateTime" xml:"CreateTime"`
	EngineVersion         string `json:"EngineVersion" xml:"EngineVersion"`
	GuardDBInstanceId     string `json:"GuardDBInstanceId" xml:"GuardDBInstanceId"`
	TempDBInstanceId      string `json:"TempDBInstanceId" xml:"TempDBInstanceId"`
	MasterInstanceId      string `json:"MasterInstanceId" xml:"MasterInstanceId"`
	VpcId                 string `json:"VpcId" xml:"VpcId"`
	VSwitchId             string `json:"VSwitchId" xml:"VSwitchId"`
	ReplicateId           string `json:"ReplicateId" xml:"ReplicateId"`
	ResourceGroupId       string `json:"ResourceGroupId" xml:"ResourceGroupId"`
	ConnectionString      string `json:"ConnectionString" xml:"ConnectionString"`
	Port                  string `json:"Port" xml:"Port"`
	DBInstanceMemory      int    `json:"DBInstanceMemory" xml:"DBInstanceMemory"`
	DBInstanceStorage     int    `json:"DBInstanceStorage" xml:"DBInstanceStorage"`
}

func (Rds) TableName() string {
	return "rds"
}

func GetCloudAccountKeyAndSecretByRdsInstanceId(instanceId string) (cloudAccount CloudAccount) {
	server := GetRdsDetailByInstanceId(instanceId)
	database.MysqlClient.Model(&cloudAccount).Where("id=?", server.CloudAccountId).Find(&cloudAccount)
	return
}

func GetRdsCount(userId uint, queryExpiredTime string, queryKeyword string,
	queryCloudAccount int) (total int) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	total = 0
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps {
		querySql = "select count(distinct rds.id) from rds where rds.data_status > 0 "
	} else {
		querySql = "select count(distinct rds.id) from rds inner join role_resources as rr " +
			"on rds.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"where rds.data_status > 0 and rr.resource_type = 'rds' and ur.user_id = ? "
		args = append(args, userId)
	}
	if queryExpiredTime != "" {
		querySql = fmt.Sprintf("%s and ( expire_time = '' or expire_time <= ? ) ", querySql)
		args = append(args, queryExpiredTime)
	} else {
		querySql = fmt.Sprintf("%s and ( expire_time = '' or expire_time >= ? ) ", querySql)
		args = append(args, nowTime)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (db_instance_description like ? "+
			"or connection_string like ? or db_instance_id like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	database.MysqlClient.Raw(querySql, args...).Count(&total)
	return
}

func GetRdsByPage(userId uint, offset int, limit int, queryExpiredTime string, queryKeyword string,
	queryCloudAccount int) (rds []RdsDetail) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps {
		querySql = "select rds.*, ca.name as cloud_account_name from rds left join cloud_account as ca " +
			"on rds.cloud_account_id = ca.id where rds.data_status > 0 and ( expire_time = '' or expire_time >= ? ) "
		args = append(args, nowTime)
	} else {
		querySql = "select rds.*, ca.name as cloud_account_name from rds inner join role_resources as rr " +
			"on rds.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"left join cloud_account as ca on rds.cloud_account_id = ca.id " +
			"where rds.data_status > 0 and rr.resource_type = 'rds' and ur.user_id = ? " +
			"and ( expire_time = '' or expire_time >= ? ) "
		args = append(args, userId, nowTime)
	}
	if queryExpiredTime != "" {
		querySql = fmt.Sprintf("%s and ( expire_time = '' or expire_time <= ? ) ", querySql)
		args = append(args, queryExpiredTime)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (db_instance_description like ? "+
			"or connection_string like ? or db_instance_id like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	querySql += " group by rds.id order by id desc limit ?, ? "
	args = append(args, nowTime)
	database.MysqlClient.Raw(querySql, args...).Scan(&rds)
	return
}

func GetRdsDetailByInstanceId(instanceId string) (rds Rds) {
	database.MysqlClient.Model(&rds).Where("db_instance_id=?", instanceId).Find(&rds)
	return
}

func SaveOrUpdateRds(instanceId string, cloudAccountId int, rds Rds) (err error) {
	count := 0
	rds.CloudAccountId = cloudAccountId
	database.MysqlClient.Model(&rds).Where("db_instance_id=?", instanceId).Count(&count)
	if count > 0 {
		err = database.MysqlClient.Model(&rds).Where("db_instance_id=?", instanceId).Updates(&rds).Error
	} else {
		err = database.MysqlClient.Create(&rds).Error
	}
	AddRdsDiffCache(instanceId)
	return
}

func AddRdsDiffCache(instanceId string) {
	database.MysqlClient.Create(&DiffCache{
		Type:       "rds",
		InstanceId: instanceId,
	})
}

func GetRdsDetail(id int) (rds RdsDetail) {
	database.MysqlClient.Raw("select rds.*, cloud_account.name as cloud_account_name "+
		"from rds left join cloud_account "+
		"on rds.cloud_account_id = cloud_account.id where rds.id=?", id).Scan(&rds)
	return
}

func GetRdsListByRoleId(roleId int) (rds []Rds, err error) {
	err = database.MysqlClient.Raw("SELECT ur.resource_id as id FROM role_resources as ur "+
		"WHERE ur.resource_type = 'rds' and ur.role_id = ?", roleId).Scan(&rds).Error
	return
}

func GetAllRdsList() (rds []Rds, err error) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = database.MysqlClient.Raw("SELECT id, db_instance_description, db_instance_id"+
		" FROM rds where data_status > 0 and ( expire_time ='' or expire_time >= ? )", nowTime).Scan(&rds).Error
	return
}

func AddRoleRds(roleId int64, rdsIdList []string) (err error) {
	resIds := strings.Join(rdsIdList, ",")
	err = database.MysqlClient.Exec("delete from role_resources where role_id = ? "+
		"and resource_type='rds' and resource_id not in (?)", roleId, resIds).Error
	for _, resId := range rdsIdList {
		exist := 0
		database.MysqlClient.Raw("select count(*) from role_resources "+
			"where resource_type='rds' and role_id = ? and resource_id = ?", roleId, resId).Count(&exist)
		if exist == 0 {
			err = database.MysqlClient.Exec("insert into role_resources(role_id, resource_type, "+
				"resource_id) values(?, ?, ?)", roleId, "rds", resId).Error
		}
	}
	return
}

func AddCloudRds(data forms.RdsInfoForm) (err error) {
	rds := Rds{
		CloudAccountId:        data.CloudAccountId,
		DBInstanceId:          data.DBInstanceId,
		DBInstanceDescription: data.DBInstanceDescription,
		Engine:                data.Engine,
		ConnectionString:      data.ConnectionString,
		Port:                  data.Port,
		DBInstanceMemory:      data.DBInstanceMemory * 1024,
		DBInstanceStorage:     data.DBInstanceStorage,
		CreateTime:            data.CreateTime,
		ExpireTime:            data.ExpireTime,
		DBInstanceStatus:      "Running",
		DataStatus:            2, // 表示用户自定义添加
	}
	err = database.MysqlClient.Create(&rds).Error
	return
}

func DeleteCloudRds(id int) (err error) {
	err = database.MysqlClient.Exec("update rds set data_status = 0 where id = ?", id).Error
	return
}

func UpdateCloudRds(data forms.ExtraInfoForm) (err error) {
	var updateSql string

	memory := "db_instance_memory"
	if data.RdsBaseInfoForm.DbMemory != "" {
		memoryInt, err1 := strconv.Atoi(data.RdsBaseInfoForm.DbMemory)
		if err1 != nil {
			err = err1
			return
		}
		memory = fmt.Sprintf("%d", memoryInt*1024)
	}

	dbInstanceDescription := "db_instance_description"
	if data.RdsBaseInfoForm.DbInstanceDescription != "" {
		dbInstanceDescription = fmt.Sprintf("'%s'", data.RdsBaseInfoForm.DbInstanceDescription)
	}

	dbInstanceStorage := "db_instance_storage"
	if data.RdsBaseInfoForm.DbInstanceStorage != "" {
		dbInstanceStorage = fmt.Sprintf("'%s'", data.RdsBaseInfoForm.DbInstanceStorage)
	}

	expireTime := "expire_time"
	if data.RdsBaseInfoForm.DbExpiredTime != "" {
		expireTime = fmt.Sprintf("'%s'", data.RdsBaseInfoForm.DbExpiredTime)
	}

	updateSql = "update rds set db_instance_description=?, db_instance_storage=?, " +
		"db_instance_memory=?, expire_time=? where id in (?)"
	err = database.MysqlClient.Exec(updateSql, dbInstanceDescription, dbInstanceStorage,
		memory, expireTime, strings.Split(data.Id, ",")).Error
	return
}
