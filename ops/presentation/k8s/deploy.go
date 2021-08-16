package k8s

import (
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/k8s"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetDeployments(c *gin.Context) {
	var req k8s.NamespaceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	data, err := cluster.GetDeployments(req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func UpdateDeploymentWithImage(clusterId string, namespace string, deployName string, imageName string) (err error) {
	apiServer, token := models.GetK8sInfoByClusterId(clusterId)
	k8sClient, err := service.NewKubernetesService(apiServer, token)
	if err != nil {
		return
	}
	err = k8sClient.UpdateDeploymentWithImage(namespace, deployName, imageName)
	if err != nil {
		return
	}
	return
}

// 重启 Deployment
func RestartDeployments(c *gin.Context) {
	var req k8s.NamespaceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = cluster.RestartDeployment(req.Namespace, req.ResName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}
