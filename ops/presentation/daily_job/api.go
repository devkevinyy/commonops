package daily_job

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/daily_job"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func IPostDailyJob(c *gin.Context) {
	var req daily_job.AddDailyJobForm
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	userId, _ := c.Get("userId")
	userName, _ := c.Get("username")
	req.CreatorUserId = int(userId.(float64))
	req.CreatorUserName = userName.(string)
	if err := service.GetDailyJobService().AddDailyJob(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "工单提交成功", Data: nil})
}

func IPutDailyJob(c *gin.Context) {
	var req daily_job.UpdateDailyJobForm
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	userId, _ := c.Get("userId")
	userName, _ := c.Get("username")
	req.UserId = int(userId.(float64))
	req.UserName = userName.(string)
	if err := service.GetDailyJobService().UpdateDailyJob(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "操作成功", Data: nil})
}

func IPutDailyJobExecutorUser(c *gin.Context) {
	var req daily_job.UpdateDailyJobExecutorUserForm
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err := service.GetDailyJobService().UpdateDailyJobExecutorUser(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "修改成功", Data: nil})
}

func IGetDailyJobs(c *gin.Context) {
	var req daily_job.JobsQueryForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	total, slbList := service.GetDailyJobService().GetDailyJobDataByPage(utils.GetCurrentUserId(c), req.QueryKeyword,
		req.QueryCreateTime, req.Page, req.Size)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: daily_job.DailyJobResp{
		Total: total,
		Page:  req.Page,
		Jobs:  slbList,
	}})
}

func IGetDailyJobDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg:"需要传入整型类型的任务id"})
		return
	}
	info := service.GetDailyJobService().GetDailyJobDetail(id)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: info})
}

