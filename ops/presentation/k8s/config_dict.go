package k8s

import (
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/k8s"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetConfigDict(c *gin.Context) {
	var req k8s.NamespaceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	data, err := cluster.GetConfigDict(req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code:0, Msg: "success", Data:data})
}

func GetSecretDict(c *gin.Context) {
	var req k8s.NamespaceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	data, err := cluster.GetSecretDict(req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code:0, Msg: "success", Data:data})
}

func DeleteConfigMap(c *gin.Context) {
	var req k8s.ResourceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK,  utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK,  utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	err = cluster.DeleteConfigDict(req.Namespace, req.ResName)
	if err != nil {
		c.JSON(http.StatusOK,  utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code:0, Msg: "success", Data:nil})
}

func DeleteSecret(c *gin.Context) {
	var req k8s.ResourceForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	err = cluster.DeleteSecretDict(req.Namespace, req.ResName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code:-1, Msg: err.Error(), Data:nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code:0, Msg: "success", Data:nil})
}