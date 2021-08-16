package models

import (
	"fmt"
	"strings"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type Slb struct {
	gorm.Model
	DataStatus         int8   `json:"DataStatus" gorm:"not null;default: 1"`
	CloudAccountId     uint    `json:"CloudAccountId"`
	Count              uint    `json:"Count" xml:"Count"`
	SlaveZoneId        string `json:"SlaveZoneId" xml:"SlaveZoneId"`
	LoadBalancerStatus string `json:"LoadBalancerStatus" xml:"LoadBalancerStatus"`
	VSwitchId          string `json:"VSwitchId" xml:"VSwitchId"`
	MasterZoneId       string `json:"MasterZoneId" xml:"MasterZoneId"`
	PayType            string `json:"PayType" xml:"PayType"`
	RegionIdAlias      string `json:"RegionIdAlias" xml:"RegionIdAlias"`
	CreateTime         string `json:"CreateTime" xml:"CreateTime"`
	Address            string `json:"Address" xml:"Address"`
	LoadBalancerId     string `json:"LoadBalancerId" xml:"LoadBalancerId"`
	AddressIPVersion   string `json:"AddressIPVersion" xml:"AddressIPVersion"`
	RegionId           string `json:"RegionId" xml:"RegionId"`
	ResourceGroupId    string `json:"ResourceGroupId" xml:"ResourceGroupId"`
	LoadBalancerName   string `json:"LoadBalancerName" xml:"LoadBalancerName"`
	InternetChargeType string `json:"InternetChargeType" xml:"InternetChargeType"`
	AddressType        string `json:"AddressType" xml:"AddressType"`
	VpcId              string `json:"VpcId" xml:"VpcId"`
	NetworkType        string `json:"NetworkType" xml:"NetworkType"`
}

func (Slb) TableName() string {
	return "slb"
}

func GetCloudAccountKeyAndSecretBySlbInstanceId(instanceId string) (cloudAccount CloudAccount) {
	server := GetSlbDetailByInstanceId(instanceId)
	database.Mysql().Model(&cloudAccount).Where("id=?", server.CloudAccountId).Find(&cloudAccount)
	return
}

func GetSlbDetailByInstanceId(instanceId string) (slb Slb) {
	database.Mysql().Model(&slb).Where("load_balancer_id=?", instanceId).Find(&slb)
	return
}

func GetSlbCount(userId uint, queryKeyword string, queryCloudAccount uint) (total uint) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	total = 0
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps {
		querySql = "select count(distinct slb.id) from slb where slb.data_status > 0 "
	} else {
		querySql = "select count(distinct slb.id) from slb inner join role_resources as rr " +
			"on slb.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"where slb.data_status > 0 and rr.resource_type = 'slb' and ur.user_id = ? "
		args = append(args, userId)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and address like ? or load_balancer_name like ? ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	database.Mysql().Raw(querySql, args...).Count(&total)
	return
}

func GetSlbByPage(userId uint, queryKeyword string, queryCloudAccount uint, offset uint, limit uint) (slb []SlbDetail) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	querySql := ""
	var args []interface{}
	if isSuperAdminOrOps {
		querySql = "select slb.*, ca.name as cloud_account_name from slb " +
			"left join cloud_account as ca on slb.cloud_account_id = ca.id where slb.data_status > 0 "
	} else {
		querySql = "select slb.*, ca.name as cloud_account_name from slb inner join role_resources as rr " +
			"on slb.id = rr.resource_id inner join user_roles as ur on rr.role_id = ur.role_id " +
			"left join cloud_account as ca on slb.cloud_account_id = ca.id " +
			"where slb.data_status > 0 and rr.resource_type = 'slb' and ur.user_id = ? "
		args = append(args, userId)
	}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (slb.address like ? "+
			"or slb.load_balancer_name like ? or load_balancer_id like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCloudAccount != 0 {
		querySql = fmt.Sprintf("%s and cloud_account_id = ? ", querySql)
		args = append(args, queryCloudAccount)
	}
	querySql += " group by slb.id order by slb.id desc limit ?, ? "
	args = append(args, offset, limit)
	database.Mysql().Raw(querySql, args...).Scan(&slb)
	return
}

func SaveOrUpdateSlb(instanceId string, cloudAccountId uint, slb Slb) (err error) {
	count := 0
	slb.CloudAccountId = cloudAccountId
	database.Mysql().Model(&slb).Where("load_balancer_id=?", instanceId).Count(&count)
	if count > 0 {
		err = database.Mysql().Model(&slb).Where("load_balancer_id=?", instanceId).Omit("renew_site_id",
			"manage_user", "cost", "cost_type", "bank_account", "renew_status").Updates(&slb).Error
	} else {
		err = database.Mysql().Create(&slb).Error
	}
	AddSlbDiffCache(instanceId)
	return
}

func AddSlbDiffCache(instanceId string) {
	database.Mysql().Create(&DiffCache{
		Type:       "slb",
		InstanceId: instanceId,
	})
}

func GetSlbDetail(id uint) (slb SlbDetail) {
	database.Mysql().Raw("select slb.*, cloud_account.name as cloud_account_name "+
		"from slb left join cloud_account "+
		"on slb.cloud_account_id = cloud_account.id where slb.id=?", id).Scan(&slb)
	return
}

func GetSlbListByRoleId(roleId int) (slb []Slb, err error) {
	err = database.Mysql().Raw("SELECT ur.resource_id as id FROM role_resources as ur "+
		"WHERE ur.resource_type = 'slb' and ur.role_id = ?", roleId).Scan(&slb).Error
	return
}

func GetAllSlbList() (slb []Slb, err error) {
	err = database.Mysql().Raw("SELECT id, address, load_balancer_name " +
		" FROM slb where data_status > 0").Scan(&slb).Error
	return
}

func AddRoleSlb(roleId int64, slbIdList []string) (err error) {
	resIds := strings.Join(slbIdList, ",")
	err = database.Mysql().Exec("delete from role_resources where role_id = ? "+
		"and resource_type='slb' and resource_id not in (?)", roleId, resIds).Error
	for _, resId := range slbIdList {
		exist := 0
		database.Mysql().Raw("select count(*) from role_resources "+
			"where resource_type='slb' and role_id = ? and resource_id = ?", roleId, resId).Count(&exist)
		if exist == 0 {
			err = database.Mysql().Exec("insert into role_resources(role_id, resource_type, "+
				"resource_id) values(?, ?, ?)", roleId, "slb", resId).Error
		}
	}
	return
}

func DeleteCloudSlb(id uint) (err error) {
	err = database.Mysql().Exec("update slb set data_status = 0 where id = ?", id).Error
	return
}

func CleanSlbDiffCaches() error {
	return database.Mysql().Exec("delete from diff_caches where type = 'slb' ").Error
}