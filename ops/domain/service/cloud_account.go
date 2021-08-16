package service

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	cloud_account2 "github.com/chujieyang/commonops/ops/value_objects/cloud_account"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/11/21 4:38 PM
 * @Desc:
 */

type cloudAccount struct {}

var cloudAccountInstance = &cloudAccount{}

func GetCloudAccountService() *cloudAccount {
	return cloudAccountInstance
}

func (c *cloudAccount) GetCloudAccountDataByPage(page, size uint) (total uint, data []models.CloudAccount) {
	offset := (page - 1) * size
	total = models.GetCloudAccountsCount()
	data = models.GetCloudAccounts(offset, size)
	return
}

func (c *cloudAccount) AddCloudAccount(data cloud_account2.CloudAccountForm) error {
	return models.AddCloudAccount(data)
}

func (c *cloudAccount) UpdateCloudAccount(data cloud_account2.CloudAccountForm) error {
	return models.UpdateCloudAccount(data)
}

func (c *cloudAccount) DeleteCloudAccount(id uint) error {
	return models.DeleteCloudAccount(id)
}
