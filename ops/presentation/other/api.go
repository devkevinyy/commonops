package other

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/utils"
	other2 "github.com/chujieyang/commonops/ops/value_objects/other"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 12:06 PM
 * @Desc:
 */

/*
 [api get]: 获取其他资源列表
*/
func IGetCloudOtherRes(c *gin.Context) {
	var req other2.CloudOtherResQueryForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	total, resList := service.GetOtherResService().GetOtherResDataByPage(utils.GetCurrentUserId(c), req.QueryKeyword,
		req.QueryResType, req.QueryCloudAccount, req.QueryExpiredTime, req.QueryManageUser, req.Page, req.Size)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data:  other2.CloudOtherResResp{
		Total:    total,
		Page:     req.Page,
		OtherRes: resList,
	}})
}

/*
 [api post]: 新增其他资源信息
*/
func IPostCloudOtherRes(c *gin.Context) {
	var req other2.AddOtherResForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetOtherResService().AddCloudOtherRes(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
 [api delete]: 删除其他资源信息
*/
func IDeleteCloudOtherRes(c *gin.Context) {
	var req other2.ResDeleteForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetOtherResService().DeleteCloudOtherRes(req.Id); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}