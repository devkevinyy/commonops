package data

import (
	"github.com/chujieyang/commonops/ops/cron"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IGetAliyunEcsSync(c *gin.Context) {
	cron.SyncAliYunEcsData()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "操作成功", Data: nil})
}

func IGetAliyunRdsSync(c *gin.Context) {
	cron.SyncAliYunRdsData()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "操作成功", Data: nil})
}

func IGetAliyunKvSync(c *gin.Context) {
	cron.SyncAliYunKvData()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "操作成功", Data: nil})
}

func IGetAliyunSlbSync(c *gin.Context) {
	cron.SyncAliYunSlbData()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "操作成功", Data: nil})
}

func IGetAliyunStatisData(c *gin.Context) {
	data := models.GetEcsRdsKvSlbCount()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "查询成功", Data: data})
}

