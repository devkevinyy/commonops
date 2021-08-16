package service

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/value_objects/other"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 1:52 PM
 * @Desc:
 */

type otherRes struct {}

var otherResInstance = &otherRes{}

func GetOtherResService() *otherRes {
	return otherResInstance
}

func (o *otherRes) GetOtherResDataByPage(userId uint, queryKeyword string, resType string, cloudAccountId uint,
	expiredTime string, manageUser uint, page uint, size uint) (total uint, data []models.OtherDetail) {
	offset := (page - 1) * size
	total = models.GetOtherResCount(userId, queryKeyword, resType, cloudAccountId, expiredTime, manageUser)
	data = models.GetOtherRes(userId, queryKeyword, resType, cloudAccountId, expiredTime, manageUser, offset, size)
	return
}

func (o *otherRes) AddCloudOtherRes(params other.AddOtherResForm) error {
	return models.AddCloudOtherRes(params)
}

func (o *otherRes) DeleteCloudOtherRes(id uint) error {
	return models.DeleteCloudOtherRes(id)
}