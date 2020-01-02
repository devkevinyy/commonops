package k8s

import (
	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/chujieyang/commonops/ops/untils/k8s_utils"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"net/http"
)

func getContextCluster(c *gin.Context) (client *k8s_utils.K8sClientSet, err error) {
	client = c.Keys["k8sCluster"].(*k8s_utils.K8sClientSet)
	return
}

func GetNamespaces(c *gin.Context) {
	cluster, err := getContextCluster(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	data, err := cluster.GetNamespaces()
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}

func PostNamespaces(c *gin.Context) {
	var req k8s_structs.NamespaceForm
	err := c.Bind(&req)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	data, err := cluster.CreateNamespaces(req.Namespace)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}

func DeleteNamespaces(c *gin.Context) {
	var req k8s_structs.NamespaceForm
	err := c.Bind(&req)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	err = cluster.DeleteNamespaces(req.Namespace)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(nil))
}