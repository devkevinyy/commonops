package dms

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/forms"
	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/services/dms_service"
	"github.com/chujieyang/commonops/ops/untils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type DmsAuthDataResp struct {
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	AuthData []models.DmsAuth `json:"authData"`
}

type DmsInstanceDataResp struct {
	Total        int                  `json:"total"`
	Page         int                  `json:"page"`
	InstanceData []models.DmsInstance `json:"instanceData"`
}

// 添加实例
func IPostDmsInstance(c *gin.Context) {
	var req forms.DmsInstanceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	instanceId := req.InstanceId
	if strings.TrimSpace(req.InstanceId) == "" {
		instanceId = untils.GetUUID()
	}
	if err = models.SaveDmsInstance(req.Port, req.InstanceType, req.Host, "NORMAL", instanceId, req.InstanceAlias, req.OperUser, req.OperPwd); err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

// 删除实例
func IDeleteDmsInstance(c *gin.Context) {
	var req forms.DmsInstanceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	if err = models.DeleteDmsInstance(req.InstanceId); err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetDmsInstanceData(c *gin.Context) {
	var req forms.DmsInstanceQueryForm
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.PageSize
	total := models.GetDmsInstanceCount(req.Query)
	instanceList := models.GetDmsInstanceList(req.Query, offset, req.PageSize)
	resp := DmsInstanceDataResp{
		Total:        total,
		Page:         req.Page,
		InstanceData: instanceList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: resp})
}

func IGetAllDmsInstancesData(c *gin.Context) {
	instanceList, err := models.GetDmsAllInstanceData()
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: instanceList})
}

func IDeleteDmsDatabaseByDatabaseId(c *gin.Context) {
	var req forms.DmsInstanceDbDeleteForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	err = models.DeleteDmsDatabase(req.InstanceId, req.DatabaseId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPostDmsDatabase(c *gin.Context) {
	var req forms.DmsInstanceDbAddForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	instance, err := models.GetDmsInstanceByInstanceId(req.InstanceId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	err = models.SaveDmsDatabase(req.InstanceId, untils.GetUUID(), req.DbName, "NORMAL",
		instance.InstanceType, instance.Host, instance.Port)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetDmsAuthData(c *gin.Context) {
	var req forms.DmsUserAuthQueryForm
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.Size
	total := models.GetUserAuthCount(req.EmpId, req.InstanceId, req.OperType)
	authList := models.GetUserAuthList(req.EmpId, req.InstanceId, req.OperType, offset, req.Size)
	resp := DmsAuthDataResp{
		Total:    total,
		Page:     req.Page,
		AuthData: authList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: resp})
}

func IGetAllDmsInstanceData(c *gin.Context) {
	data, err := models.GetDmsAllInstanceData()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: data})
}

func IGetDmsDatabaseByInstanceId(c *gin.Context) {
	instanceId := c.Query("instanceId")
	data, err := models.GetDmsAllDatabaseData(instanceId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: data})
}

func IGetDmsTablesByDatabaseId(c *gin.Context) {
	databaseId := c.Query("databaseId")
	data, err := models.GetDmsAllTableData(databaseId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: data})
}

func IPostDmsUserAuth(c *gin.Context) {
	isAdminOrOps := models.IsUserSuperAdminOrOps(untils.GetCurrentUserId(c))
	if !isAdminOrOps {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "你没有该操作权限！"})
		return
	}
	var req forms.DmsAddUserAuthForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "1" + err.Error()})
		return
	}
	userInfo, err := models.GetUserInfoByEmpId(req.EmpId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "2" + err.Error()})
		return
	}
	if req.SelectedNodeType != "database" {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "用户权限需要指定到库！"})
		return
	}
	dbInfo, err := models.GetDmsDatabaseInfoByDatabaseId(req.SelectedNodeId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "3" + err.Error()})
		return
	}
	instanceInfo, err := models.GetDmsInstanceByInstanceId(dbInfo.InstanceId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "4" + err.Error()})
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
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "5" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IDeleteDmsUserAuth(c *gin.Context) {
	var req forms.DmsDeleteUserAuthForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	err := models.DmsDeleteUserAuth(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetUserDmsInstanceData(c *gin.Context) {
	empId, _ := c.Get("empId")
	data, err := models.GetUserInstanceData(empId.(string))
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: data})
}

func IGetUserDmsDatabaseData(c *gin.Context) {
	empId, _ := c.Get("empId")
	instanceId := c.Query("instanceId")
	if strings.TrimSpace(instanceId) == "" {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "instanceId不能为空"})
		return
	}
	data, err := models.GetUserDatabaseData(empId.(string), instanceId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: data})
}

// 用户提交执行SQL
func IPostDmsUserExecSQL(c *gin.Context) {
	empId, _ := c.Get("empId")
	var req forms.DmsUserExecSqlForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	if req.SelectedNodeType != "database" {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "用户权限需要指定到库！"})
		return
	}

	dbInfo, err := models.GetDmsDatabaseInfoByDatabaseId(req.SelectedNodeId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}

	instanceOperInfo, err := models.GetDmsInstanceOperUserInfoByInstanceId(dbInfo.InstanceId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
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
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}

	if sqlType == "select" {
		hasAuth, err := models.CheckUserDMSAuthDatabaseQuery(empId.(string), req.SelectedNodeId, tableNames)
		if !hasAuth || err != nil {
			c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "你目前没有查询该库或表的对应权限！"})
			return
		}
	}
	if sqlType == "insert" || sqlType == "update" || sqlType == "delete" {
		if sqlType == "update" { // 不允许提交join相关的更新
			err = models.CheckUserSqlExistJoin(strings.TrimSpace(req.SqlInput))
			if err != nil {
				c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
				return
			}
		}
		hasAuth, err := models.CheckUserDMSAuthDatabaseModify(empId.(string), req.SelectedNodeId, tableNames)
		if !hasAuth || err != nil  {
			c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "你目前没有修改该库或表的对应权限！"})
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
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: fmt.Sprintf("不支持的数据库类型: %s", dbInfo.DbType)})
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
	tx := database.MysqlClient.Begin()
	log := models.DmsLog{
		EmpId:           empId.(string),
		Username:        untils.GetCurrentUsername(c),
		DatabaseId:      dbInfo.DatabaseId,
		DatabaseName:    dbInfo.SchemaName,
		StartTime:       untils.GetCurrentTime(),
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
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	tx.Commit()

	err = models.DecreaseUserOperatorCount(empId.(string), operType, dbInfo.DatabaseId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
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
	var req forms.DmsUserQueryLog
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.PageSize
	count, err := models.GetDmsQueryLogCount(empId.(string))
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	log, err := models.GetDmsQueryLogByPage(empId.(string), offset, req.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"total": count,
		"page":  req.Page,
		"log":   log,
	}})
}

func IGetDmsUserApprovesData(c *gin.Context) {
	empId, _ := c.Get("empId")
	var req forms.DmsUserQueryLog
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.PageSize
	count, err := models.GetDmsUserApproveCount(empId.(string))
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	approves, err := models.GetDmsUserApproveByPage(empId.(string), offset, req.PageSize)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"total":    count,
		"page":     req.Page,
		"approves": approves,
	}})
}
