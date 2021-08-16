package dms

import (
	"database/sql"
	"fmt"
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/value_objects/dms"
	"net/http"
	"strings"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/services/dms_service"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 添加实例
func IPostDmsInstance(c *gin.Context) {
	var req dms.DmsInstanceForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	instanceId := req.InstanceId
	if strings.TrimSpace(req.InstanceId) == "" {
		instanceId = utils.GetUUID()
	}
	if err := service.GetDmsService().SaveDmsInstance(req, instanceId); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

// 删除实例
func IDeleteDmsInstance(c *gin.Context) {
	var req dms.DmsInstanceForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	if err := service.GetDmsService().DeleteDmsInstance(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetDmsInstanceData(c *gin.Context) {
	var req dms.DmsInstanceQueryForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	total, instanceList := service.GetDmsService().GetDmsInstanceDataByPage(req.Query, req.Page, req.PageSize)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: dms.DmsInstanceDataResp{
		Total:        total,
		Page:         req.Page,
		InstanceData: instanceList,
	}})
}

func IGetAllDmsInstancesData(c *gin.Context) {
	instanceList, err := service.GetDmsService().GetAllInstanceData()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: instanceList})
}

func IDeleteDmsDatabaseByDatabaseId(c *gin.Context) {
	var req dms.DmsInstanceDbDeleteForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	if err := service.GetDmsService().DeleteDmsDbInstance(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPostDmsDatabase(c *gin.Context) {
	var req dms.DmsInstanceDbAddForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	instance, err := service.GetDmsService().GetDmsInstanceByInstanceId(req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	if err = service.GetDmsService().AddDmsDatabase(req, instance); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetDmsAuthData(c *gin.Context) {
	var req dms.DmsUserAuthQueryForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	total, authList := service.GetDmsService().GetDmsAuthDataByPage(req.EmpId, req.InstanceId, req.OperType, req.Page, req.Size)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: dms.DmsAuthDataResp{
		Total:    total,
		Page:     req.Page,
		AuthData: authList,
	}})
}

func IGetDmsDatabaseByInstanceId(c *gin.Context) {
	instanceId := c.Query("instanceId")
	data, err := service.GetDmsService().GetDmsAllDatabaseData(instanceId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func IPostDmsUserAuth(c *gin.Context) {
	isAdminOrOps := models.IsUserSuperAdminOrOps(utils.GetCurrentUserId(c))
	if !isAdminOrOps {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "你没有该操作权限！"})
		return
	}
	var req dms.DmsAddUserAuthForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	if err := service.GetDmsService().AddDmsUserAuth(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IDeleteDmsUserAuth(c *gin.Context) {
	var req dms.DmsDeleteUserAuthForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	err := service.GetDmsService().DeleteDmsUserAuth(req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetUserDmsInstanceData(c *gin.Context) {
	empId, _ := c.Get("empId")
	data, err := service.GetDmsService().GetUserDmsInstanceData(empId.(string))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func IGetUserDmsDatabaseData(c *gin.Context) {
	empId, _ := c.Get("empId")
	instanceId := c.Query("instanceId")
	if strings.TrimSpace(instanceId) == "" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "instanceId不能为空"})
		return
	}
	data, err := service.GetDmsService().GetUserDmsDatabaseDataByInstanceId(empId.(string), instanceId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

// 用户提交执行SQL
func IPostDmsUserExecSQL(c *gin.Context) {
	empId, _ := c.Get("empId")
	var req dms.DmsUserExecSqlForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	if req.SelectedNodeType != "database" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "用户权限需要指定到库！"})
		return
	}

	dbInfo, err := models.GetDmsDatabaseInfoByDatabaseId(req.SelectedNodeId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}

	instanceOperInfo, err := models.GetDmsInstanceOperUserInfoByInstanceId(dbInfo.InstanceId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}

	dbType := ""
	if dbInfo.DbType == "2" {
		dbType = "mysql"
	}
	if dbInfo.DbType == "3" {
		dbType = "sqlserver"
	}
	sqlType, tableNames, err := dms_service.DmsSQLParser(dbType, strings.TrimSpace(req.SqlInput))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}

	if sqlType == "select" {
		hasAuth, err := models.CheckUserDMSAuthDatabaseQuery(empId.(string), req.SelectedNodeId, tableNames)
		if !hasAuth || err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "你目前没有查询该库或表的对应权限！"})
			return
		}
	}
	if sqlType == "insert" || sqlType == "update" || sqlType == "delete" {
		if sqlType == "update" { // 不允许提交join相关的更新
			err = models.CheckUserSqlExistJoin(strings.TrimSpace(req.SqlInput))
			if err != nil {
				c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
				return
			}
		}
		hasAuth, err := models.CheckUserDMSAuthDatabaseModify(empId.(string), req.SelectedNodeId, tableNames)
		if !hasAuth || err != nil  {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "你目前没有修改该库或表的对应权限！"})
			return
		}
	}

	// SQL 执行计划分析
	dbConnectionString := ""
	switch dbType {
	case "mysql", "polardb":
		dbConnectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10s",
			instanceOperInfo.OperUser, instanceOperInfo.OperPwd, dbInfo.Host, dbInfo.Port, dbInfo.SchemaName)
	case "sqlserver":
		dbConnectionString = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&connection+timeout=30",
			instanceOperInfo.OperUser, instanceOperInfo.OperPwd, dbInfo.Host, dbInfo.Port, dbInfo.SchemaName)
	default:
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: fmt.Sprintf("不支持的数据库类型: %s", dbInfo.DbType)})
		return
	}

	execStatus := 1
	exceptionOutput := ""
	var effectRows int64
	var queryColumns []string
	var queryResult []map[string]string
	var execResult sql.Result
	var duration int64
	var operType int
	switch sqlType {
	case "select":
		operType = 2
		queryColumns, queryResult, duration, err = dms_service.DmsQuery(dbType, dbConnectionString, dbInfo.SchemaName, req.SqlInput)
	case "update", "insert", "delete":
		operType = 3
		execResult, duration, err = dms_service.DmsExec(sqlType, dbType, dbConnectionString, dbInfo.SchemaName, req.SqlInput)
		if execResult != nil {
			effectRows, _ = execResult.RowsAffected()
		}
	default:
		operType = 0
		queryResult, duration, err = nil, 0, errors.New(fmt.Sprintf("当前仅支持select、update、insert, delete操作，该操作: %s 被拒绝执行", sqlType))
	}
	if err != nil {
		execStatus = 0
		exceptionOutput = err.Error()
	}

	// 写数据到记录表
	tx := database.Mysql().Begin()
	log := models.DmsLog{
		EmpId:           empId.(string),
		Username:        utils.GetCurrentUsername(c),
		DatabaseId:      dbInfo.DatabaseId,
		DatabaseName:    dbInfo.SchemaName,
		StartTime:       utils.GetCurrentTime(),
		SqlContent:      strings.TrimSpace(req.SqlInput),
		ExecStatus:      execStatus,
		Duration:        int(duration),
		EffectRows:      int(effectRows),
		Result:          "",
		ExceptionOutput: exceptionOutput,
		HasExecuted:     1,
		SqlType:         sqlType,
	}
	err = tx.Create(&log).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	tx.Commit()

	err = models.DecreaseUserOperatorCount(empId.(string), operType, dbInfo.DatabaseId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"sqlType":         sqlType,
		"execStatus":      execStatus,
		"effectRows":      effectRows,
		"queryColumns":    queryColumns,
		"queryResult":     queryResult,
		"exceptionOutput": exceptionOutput,
		"duration":        duration,
	}})
}

func IGetDmsQueryLogData(c *gin.Context) {
	empId, _ := c.Get("empId")
	var req dms.DmsUserQueryLog
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	count, log := service.GetDmsService().GetDmsQueryLogDataByPage(empId.(string), req.Page, req.PageSize)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"total": count,
		"page":  req.Page,
		"log":   log,
	}})
}

