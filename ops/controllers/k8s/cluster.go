package k8s

import (
	"net/http"

	"github.com/chujieyang/commonops/ops/forms"
	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/untils"
	"github.com/gin-gonic/gin"
)

func PostK8sCluster(c *gin.Context) {
	var req forms.K8sForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	err = models.AddK8sCluster(req)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

func DeleteK8sCluster(c *gin.Context) {
	var req forms.K8sForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	err = models.DeleteK8s(req)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

func GetK8sCluster(c *gin.Context) {
	var req forms.K8sForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	totalCount := models.GetK8sCount()
	data := models.GetK8sList()
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"totalCount": totalCount,
		"k8sData":    data,
	}})
}
