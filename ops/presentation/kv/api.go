package kv

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/utils"
	kv2 "github.com/chujieyang/commonops/ops/value_objects/kv"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:22 AM
 * @Desc:
 */

/*
 [api get]: 获取公有云的kv详情
*/
func IGetCloudKvDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "需要传入整型类型的服务器id"})
		return
	}
	kvDetail := service.GetKvService().GetKvDetail(uint(id))
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: kvDetail})
}

/*
 [api get]: 获取公有云的 kv 列表
*/
func IGetCloudKv(c *gin.Context) {
	var req kv2.KvQueryForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	total, kvList := service.GetKvService().GetKvDataByPage(utils.GetCurrentUserId(c), req.QueryExpiredTime, req.QueryKeyword,
		req.QueryCloudAccount, req.Page, req.Size)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: kv2.KvResp{
		Total: total,
		Page:  req.Page,
		Kv:    kvList,
	}})
}

/*
 [api put]: 完善服务器的扩展信息
*/
func IPutCloudKv(c *gin.Context) {
	var req kv2.ExtraInfoForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err:= service.GetKvService().UpdateKv(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "修改数据成功！"})
}

/*
 [api post]: 新增redis信息
*/
func IPostCloudKv(c *gin.Context) {
	var req kv2.KvInfoForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetKvService().AddKv(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}

func IDeleteCloudKv(c *gin.Context) {
	var req kv2.ResDeleteForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetKvService().DeleteKv(req.Id); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}