package models

import (
	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/forms"
)

type CloudAccount struct {
	Id int              `json:"id" gorm:"column:id"`
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

func GetCloudAccountsCount() (total int) {
	total = 0
	database.MysqlClient.Model(&CloudAccount{}).Count(&total)
	return
}

func GetAllCloudAccounts() (accounts []CloudAccount) {
	accounts = []CloudAccount{}
	database.MysqlClient.Find(&accounts)
	return
}

func GetCloudAccounts(offset int, limit int) (accounts []CloudAccount) {
	accounts = []CloudAccount{}
	database.MysqlClient.Offset(offset).Limit(limit).Find(&accounts)
	return
}

func DeleteCloudAccount(id int) (err error) {
	account := CloudAccount{
		Id: id,
	}
	err = database.MysqlClient.Delete(&account).Error
	return
}

func AddCloudAccount(data forms.CloudAccountForm) (err error) {
	cloudAccount := CloudAccount{
		CloudType: data.AccountType,
		Name: data.AccountName,
		Passwd: data.AccountPwd,
		Key: data.AccountKey,
		Secret: data.AccountSecret,
		Region: data.AccountRegion,
	}
	err = database.MysqlClient.Create(&cloudAccount).Error
	return
}

func UpdateCloudAccount(data forms.CloudAccountForm) (err error) {
	cloudAccount := CloudAccount{
		Id: data.Id,
		CloudType: data.AccountType,
		Name: data.AccountName,
		Passwd: data.AccountPwd,
		Key: data.AccountKey,
		Secret: data.AccountSecret,
		Region: data.AccountRegion,
	}
	err = database.MysqlClient.Save(&cloudAccount).Error
	return
}