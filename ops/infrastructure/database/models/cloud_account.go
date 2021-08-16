package models

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database"
	cloud_account2 "github.com/chujieyang/commonops/ops/value_objects/cloud_account"
)

type CloudAccount struct {
	Id uint              `json:"id" gorm:"column:id"`
	CloudType string	`json:"accountType" gorm:"column:cloud_type"`
	Name string     	`json:"accountName" gorm:"column:name"`
	Passwd string   	`json:"accountPwd" gorm:"column:passwd"`
	Key string      	`json:"accountKey" gorm:"column:key"`
	Secret string   	`json:"accountSecret" gorm:"column:secret"`
	Region string   	`json:"accountRegion" gorm:"column:region_id"`
}


func (CloudAccount) TableName() string {
	return "cloud_account"
}

func GetCloudAccountsCount() (total uint) {
	total = 0
	database.Mysql().Model(&CloudAccount{}).Count(&total)
	return
}

func GetAllCloudAccounts() (accounts []CloudAccount) {
	accounts = []CloudAccount{}
	database.Mysql().Find(&accounts)
	return
}

func GetCloudAccounts(offset uint, limit uint) (accounts []CloudAccount) {
	accounts = []CloudAccount{}
	database.Mysql().Offset(offset).Limit(limit).Find(&accounts)
	return
}

func DeleteCloudAccount(id uint) (err error) {
	account := CloudAccount{
		Id: id,
	}
	err = database.Mysql().Delete(&account).Error
	return
}

func AddCloudAccount(data cloud_account2.CloudAccountForm) (err error) {
	cloudAccount := CloudAccount{
		CloudType: data.AccountType,
		Name: data.AccountName,
		Passwd: data.AccountPwd,
		Key: data.AccountKey,
		Secret: data.AccountSecret,
		Region: data.AccountRegion,
	}
	err = database.Mysql().Create(&cloudAccount).Error
	return
}

func UpdateCloudAccount(data cloud_account2.CloudAccountForm) (err error) {
	cloudAccount := CloudAccount{
		Id: data.Id,
		CloudType: data.AccountType,
		Name: data.AccountName,
		Passwd: data.AccountPwd,
		Key: data.AccountKey,
		Secret: data.AccountSecret,
		Region: data.AccountRegion,
	}
	err = database.Mysql().Save(&cloudAccount).Error
	return
}

func GetCloudAccountInfo(id int) (cloudAccount CloudAccount, err error) {
	err = database.Mysql().Raw("select * from cloud_account where id = ?", id).Scan(&cloudAccount).Error
	return
}

func GetCloudAccountInfoByName(name string) (cloudAccount CloudAccount, err error) {
	err = database.Mysql().Raw("select * from cloud_account where name = ?", name).Scan(&cloudAccount).Error
	return
}