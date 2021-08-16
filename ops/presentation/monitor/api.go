package monitor

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"net/http"

	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
)

func IGetCloudEcsMonitor(c *gin.Context) {
	instanceId := c.DefaultQuery("instanceId", "")
	timeDimension := c.DefaultQuery("timeDimension", "1h")
	metricDimension := c.DefaultQuery("metricDimension", "CPUUtilization")
	if instanceId == "" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "需要传入实例id"})
		return
	}
	cloudAccount := models.GetCloudAccountKeyAndSecretByEcsInstanceId(instanceId)
	if cloudAccount.CloudType != "阿里云" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "暂无阿里云之外的监控信息"})
		return
	}
	data := service.GetEcsService().GetEcsMonitor(instanceId, timeDimension, metricDimension,
		cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func IGetCloudRdsMonitor(c *gin.Context) {
	instanceId := c.DefaultQuery("instanceId", "")
	timeDimension := c.DefaultQuery("timeDimension", "1h")
	metricDimension := c.DefaultQuery("metricDimension", "CpuUsage")
	if instanceId == "" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "需要传入实例id"})
		return
	}
	cloudAccount := models.GetCloudAccountKeyAndSecretByRdsInstanceId(instanceId)
	if cloudAccount.CloudType != "阿里云" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "暂无阿里云之外的监控信息"})
		return
	}
	data := service.GetRdsService().GetRdsMonitor(instanceId, timeDimension, metricDimension,
		cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func IGetCloudKvMonitor(c *gin.Context) {
	instanceId := c.DefaultQuery("instanceId", "")
	timeDimension := c.DefaultQuery("timeDimension", "1h")
	metricDimension := c.DefaultQuery("metricDimension", "ShardingConnectionUsage")
	if instanceId == "" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "需要传入实例id"})
		return
	}
	cloudAccount := models.GetCloudAccountKeyAndSecretByKvInstanceId(instanceId)
	if cloudAccount.CloudType != "阿里云" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "暂无阿里云之外的监控信息"})
		return
	}
	data := service.GetKvService().GetKvMonitor(instanceId, timeDimension, metricDimension,
		cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func IGetCloudSlbMonitor(c *gin.Context) {
	instanceId := c.DefaultQuery("instanceId", "")
	timeDimension := c.DefaultQuery("timeDimension", "1h")
	metricDimension := c.DefaultQuery("metricDimension", "ActiveConnection")
	if instanceId == "" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "需要传入实例id"})
		return
	}
	cloudAccount := models.GetCloudAccountKeyAndSecretBySlbInstanceId(instanceId)
	if cloudAccount.CloudType != "阿里云" {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "暂无阿里云之外的监控信息"})
		return
	}
	data := service.GetSlbService().GetSlbMonitor(instanceId, timeDimension, metricDimension,
		cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}
