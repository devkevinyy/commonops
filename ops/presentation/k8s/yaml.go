package k8s

import (
	"encoding/json"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/k8s"
	"github.com/gin-gonic/gin"
	"net/http"
	"sigs.k8s.io/yaml"
)

func CreateResourceByYaml(c *gin.Context) {
	var req k8s.YamlResource
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
	data, err := cluster.ApplyYaml(req.Namespace, req.YamlContent)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
	获取资源的yaml文件
*/
func GetResourceYaml(c *gin.Context) {
	var req k8s.ResourceForm
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
	data, err := cluster.GetYamlFile(req.Namespace, req.ResType, req.ResName)
	jsonBytes, _ := json.Marshal(data)
	dataBytes, _ := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: string(dataBytes)})
}

func UpdateResourceByYaml(c *gin.Context) {
	var req k8s.YamlResource
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
	data, err := cluster.UpdateYaml(req.Namespace, req.ResName, req.YamlContent)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}
