package slb

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/utils"
	slb2 "github.com/chujieyang/commonops/ops/value_objects/slb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:38 AM
 * @Desc:
 */

/*
 [api get]: 获取公有云的slb详情
*/
func IGetCloudSlbDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "需要传入整型类型的服务器id"})
		return
	}
	slbDetail := service.GetSlbService().GetSlbDetail(uint(id))
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: slbDetail})
}

/*
 [api get]: 获取公有云的 slb 列表
*/
func IGetCloudSlb(c *gin.Context) {
	var req slb2.SlbQueryForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	total, slbList := service.GetSlbService().GetSlbDataByPage(utils.GetCurrentUserId(c), req.QueryKeyword,
		req.QueryCloudAccount, req.Page, req.Size)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data:  slb2.SlbResp{
		Total: total,
		Page:  req.Page,
		Slb:   slbList,
	}})
}

func IDeleteCloudSlb(c *gin.Context) {
	var req slb2.ResDeleteForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetSlbService().DeleteSlb(req.Id); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "删除slb成功"})
}