package k8s

import (
	"net/http"

	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/gin-gonic/gin"
)

func GetIngress(c *gin.Context) {
	var req k8s_structs.NamespaceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	data, err := cluster.GetIngress(req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}
