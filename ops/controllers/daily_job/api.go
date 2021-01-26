package daily_job

import (
	"github.com/chujieyang/commonops/ops/forms"
	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/untils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type DailyJobResp struct {
	Total int 						`json:"total"`
	Page int 						`json:"page"`
	Jobs []models.DailyJob 			`json:"jobs"`
}

func IPostDailyJob(c *gin.Context) {
	var req forms.AddDailyJobForm
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	userId, _ := c.Get("userId")
	userName, _ := c.Get("username")
	req.CreatorUserId = int(userId.(float64))
	req.CreatorUserName = userName.(string)
	err = models.AddDailyJob(req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "工单提交成功", Data: nil})
}

func IPutDailyJob(c *gin.Context) {
	var req forms.UpdateDailyJobForm
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	userId, _ := c.Get("userId")
	userName, _ := c.Get("username")
	req.UserId = int(userId.(float64))
	req.UserName = userName.(string)
	err = models.UpdateDailyJob(req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "操作成功", Data: nil})
}

func IPutDailyJobExecutorUser(c *gin.Context) {
	var req forms.UpdateDailyJobExecutorUserForm
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.UpdateDailyJobExecutorUser(req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "修改成功", Data: nil})
}

func IGetDailyJobs(c *gin.Context) {
	var req forms.JobsQueryForm
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.Size
	total := models.GetJobsCount(untils.GetCurrentUserId(c), req.QueryKeyword, req.QueryCreateTime)
	slbList := models.GetJobsByPage(untils.GetCurrentUserId(c), offset, req.Size, req.QueryKeyword, req.QueryCreateTime)
	resp := DailyJobResp{
		Total: total,
		Page: req.Page,
		Jobs: slbList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg:"success", Data: resp})
}

func IGetDailyJobDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg:"需要传入整型类型的任务id"})
		return
	}
	info := models.GetJobDetail(id)
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg:"success", Data: info})
}
