package k8s

import (
	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfigDict(c *gin.Context) {
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
	data, err := cluster.GetConfigDict(req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}

func GetSecretDict(c *gin.Context) {
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
	data, err := cluster.GetSecretDict(req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}

func DeleteConfigMap(c *gin.Context) {
	var req k8s_structs.ResourceForm
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
	err = cluster.DeleteConfigDict(req.Namespace, req.ResName)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(nil))
}

func DeleteSecret(c *gin.Context) {
	var req k8s_structs.ResourceForm
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
	err = cluster.DeleteSecretDict(req.Namespace, req.ResName)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(nil))
}