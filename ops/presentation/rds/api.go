package rds

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/utils"
	rds2 "github.com/chujieyang/commonops/ops/value_objects/rds"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:00 AM
 * @Desc:
 */

/*
 [api get]: 获取公有云的rds详情
*/
func IGetCloudRdsDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "需要传入整型类型的服务器id"})
		return
	}
	detail := service.GetRdsService().GetRdsDetail(uint(id))
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: detail})
}

/*
 [api get]: 获取公有云的 rds 列表
*/
func IGetCloudRds(c *gin.Context) {
	var req rds2.RdsQueryForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	total, rdsList := service.GetRdsService().GetRdsDataByPage(uint(utils.GetCurrentUserId(c)), req.QueryExpiredTime, req.QueryKeyword,
		req.QueryCloudAccount, req.Page, req.Size)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: rds2.RdsResp{
		Total: total,
		Page:  req.Page,
		Rds:   rdsList,
	}})
}

/*
 [api put]: 完善服务器的扩展信息
*/
func IPutCloudRds(c *gin.Context) {
	var req rds2.ExtraInfoForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err := service.GetRdsService().UpdateCloudRds(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "修改数据成功！"})
}

/*
 [api post]: 新增rds信息
*/
func IPostCloudRds(c *gin.Context) {
	var req rds2.RdsInfoForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetRdsService().AddCloudRds(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}

func IDeleteCloudRds(c *gin.Context) {
	var req rds2.ResDeleteForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err := service.GetRdsService().DeleteCloudRds(req.Id); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}