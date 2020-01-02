package k8s

import (
	"encoding/json"
	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"net/http"
	"sigs.k8s.io/yaml"
)

func CreateResourceByYaml(c *gin.Context) {
	var req k8s_structs.YamlResource
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
	data, err := cluster.ApplyYaml(req.Namespace, req.YamlContent)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}


/*
	获取资源的yaml文件
 */
func GetResourceYaml(c *gin.Context) {
	var req k8s_structs.ResourceForm
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
	data, err := cluster.GetYamlFile(req.Namespace, req.ResType, req.ResName)
	jsonBytes, _ := json.Marshal(data)
	dataBytes, _ := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(string(dataBytes)))
}