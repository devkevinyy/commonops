package service

import (
	"github.com/chujieyang/commonops/ops/exception"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	dmsForm "github.com/chujieyang/commonops/ops/value_objects/dms"
	"strings"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/16/21 3:20 PM
 * @Desc:
 */

type dms struct{}

var dmsInstance = &dms{}

func GetDmsService() *dms {
	return dmsInstance
}

func (d *dms) SaveDmsInstance(req dmsForm.DmsInstanceForm, instanceId string) error {
	return models.SaveDmsInstance(req.Port, req.InstanceType, req.Host, "NORMAL", instanceId,
		req.InstanceAlias, req.OperUser, req.OperPwd)
}

func (d *dms) DeleteDmsInstance(req dmsForm.DmsInstanceForm) error {
	return models.DeleteDmsInstance(req.InstanceId)
}

func (d *dms) GetDmsInstanceDataByPage(query string, page uint, size uint) (uint, []models.DmsInstance) {
	offset := (page - 1) * size
	total := models.GetDmsInstanceCount(query)
	instanceList := models.GetDmsInstanceList(query, offset, size)
	return total, instanceList
}

func (d *dms) GetAllInstanceData() ([]models.DmsInstance, error) {
	return models.GetDmsAllInstanceData()
}

func (d *dms) DeleteDmsDbInstance(req dmsForm.DmsInstanceDbDeleteForm) error {
	return models.DeleteDmsDatabase(req.InstanceId, req.DatabaseId)
}

func (d *dms) GetDmsInstanceByInstanceId(req dmsForm.DmsInstanceDbAddForm) (models.DmsInstance, error) {
	return models.GetDmsInstanceByInstanceId(req.InstanceId)
}

func (d *dms) AddDmsDatabase(req dmsForm.DmsInstanceDbAddForm, instance models.DmsInstance) error {
	return models.SaveDmsDatabase(req.InstanceId, utils.GetUUID(), req.DbName, "NORMAL",
		instance.InstanceType, instance.Host, instance.Port)
}

func (d *dms) GetDmsAuthDataByPage(empId string, instanceId string, operType string, page uint, size uint) (uint, []models.DmsAuth) {
	offset := (page - 1) * size
	total := models.GetUserAuthCount(empId, instanceId, operType)
	authList := models.GetUserAuthList(empId, instanceId, operType, offset, size)
	return total, authList
}

func (d *dms) GetDmsAllDatabaseData(instanceId string) ([]models.DmsDatabase, error) {
	return models.GetDmsAllDatabaseData(instanceId)
}

func (d *dms) AddDmsUserAuth(req dmsForm.DmsAddUserAuthForm) (err error) {
	userInfo, err := models.GetUserInfoByEmpId(req.EmpId)
	if err != nil {
		return
	}
	if req.SelectedNodeType != "database" {
		err = exception.DmsUserAuthNeedValidDatabaseException
		return
	}
	dbInfo, err := models.GetDmsDatabaseInfoByDatabaseId(req.SelectedNodeId)
	if err != nil {
		return
	}
	instanceInfo, err := models.GetDmsInstanceByInstanceId(dbInfo.InstanceId)
	if err != nil {
		return
	}
	approveEmpId := "0"
	if strings.TrimSpace(req.ApproveEmpId) != "" {
		approveEmpId = req.ApproveEmpId
	}
	err = models.AddUserDmsAuth(req.EmpId, userInfo.UserName, 2,
		instanceInfo.InstanceType, instanceInfo.InstanceId, instanceInfo.InstanceAlias,
		req.SelectedNodeId, dbInfo.SchemaName, req.OperType, req.ValidTime, req.OperCount,
		req.AllowTables, approveEmpId)
	if err != nil {
		return
	}
	return
}

func (d *dms)  DeleteDmsUserAuth(req dmsForm.DmsDeleteUserAuthForm) error {
	return models.DmsDeleteUserAuth(req.Id)
}

func (d *dms) GetUserDmsInstanceData(empId string) ([]models.DmsAuth, error) {
	return models.GetUserInstanceData(empId)
}

func (d *dms) GetUserDmsDatabaseDataByInstanceId(empId string, instanceId string) ([]models.DmsAuth, error) {
	return models.GetUserDatabaseData(empId, instanceId)
}

func (d *dms) GetDmsQueryLogDataByPage(empId string, page uint, size uint) (uint, []models.DmsLog) {
	offset := (page - 1) * size
	count := models.GetDmsQueryLogCount(empId)
	log := models.GetDmsQueryLogByPage(empId, offset, size)
	return count, log
}
