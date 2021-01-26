package k8s

import (
	"net/http"

	"github.com/chujieyang/commonops/ops/database"
	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/gin-gonic/gin"
)

func GetNodes(c *gin.Context) {
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	data, err := cluster.GetNodes()
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}

func GetNodesMetrics(c *gin.Context) {
	var req k8s_structs.MetricsForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	metricsList, err := database.GetNodeMetrics(req.ClusterId, req.MetricName, database.ResourceSelector{
		Namespace:    req.Namespace,
		ResourceName: req.NodeName,
	})
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(metricsList))
}
