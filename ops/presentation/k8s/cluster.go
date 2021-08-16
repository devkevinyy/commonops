package k8s

import (
	"github.com/chujieyang/commonops/ops/value_objects/k8s"
	"net/http"

	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
)

func PostK8sCluster(c *gin.Context) {
	var req k8s.K8sForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.AddK8sCluster(req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func DeleteK8sCluster(c *gin.Context) {
	var req k8s.K8sForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.DeleteK8s(req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func GetK8sCluster(c *gin.Context) {
	var req k8s.K8sForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	totalCount := models.GetK8sCount()
	data := models.GetK8sList()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"totalCount": totalCount,
		"k8sData":    data,
	}})
}
