package cloud_account

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/exception"
	"github.com/chujieyang/commonops/ops/value_objects/cloud_account"
	"net/http"
	"strconv"

	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
)

/*
 [api get]: 获取云账号列表
*/
func IGetCloudAccounts(c *gin.Context) {
	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, sizeErr := strconv.Atoi(c.DefaultQuery("size", "10"))
	if pageErr != nil || sizeErr != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: exception.ArgsException.Error()})
		return
	}
	total, accountList := service.GetCloudAccountService().GetCloudAccountDataByPage(uint(page), uint(size))
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data:  cloud_account.CloudAccountResp{
		Total:    total,
		Page:     uint(page),
		Accounts: accountList,
	}})
}

/*
 [api post]: 添加云账号信息
*/
func IPostCloudAccounts(c *gin.Context) {
	var req cloud_account.CloudAccountForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err := service.GetCloudAccountService().AddCloudAccount(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}

/*
 [api put]: 更新云账号信息
*/
func IPutCloudAccounts(c *gin.Context) {
	var req cloud_account.CloudAccountForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err := service.GetCloudAccountService().UpdateCloudAccount(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}

/*
 [api delete]: 删除云账号信息
*/
func IDeleteCloudAccounts(c *gin.Context) {
	var req cloud_account.CloudAccountForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetCloudAccountService().DeleteCloudAccount(req.Id); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}
