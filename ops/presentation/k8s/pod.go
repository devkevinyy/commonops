package k8s

import (
	"encoding/json"
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/k8s"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/tools/remotecommand"
)

func GetPods(c *gin.Context) {
	var req k8s.NamespaceForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	data, err := cluster.GetPods(req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func GetPodLogs(c *gin.Context) {
	var req k8s.PodContainerLogForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	cluster, err := getContextCluster(c)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	data, err := cluster.GetPodContainerLogs(req.Namespace, req.PodName, req.ContainerName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func GetContainerTerminal(c *gin.Context) {
	var (
		wsConn *service.WsConnection
		err    error
	)

	if wsConn, err = service.NewWebSocketService(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}

	msg, err := wsConn.WsRead()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}

	inputMap := map[string]interface{}{}
	msgMap := map[string]interface{}{}
	err = json.Unmarshal(msg.Data, &msgMap)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}

	if msgMap["type"].(string) == "input" {
		inputData := msgMap["input"].(string)
		err = json.Unmarshal([]byte(inputData), &inputMap)
		if err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		if inputMap["action"] == "init_connection" {
			handler := &service.StreamHandler{
				WsConn:      wsConn,
				ResizeEvent: make(chan remotecommand.TerminalSize),
			}
			handler.WsConn.InitLangEnv()

			apiServer, token := models.GetK8sInfoByClusterId(inputMap["clusterId"].(string))
			service.RemoteCommandContainerExec(apiServer, token, inputMap["namespace"].(string),
				inputMap["podName"].(string), inputMap["containerName"].(string), handler)
			return
		}
	}
	c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "ws error", Data: nil})
}

func GetPodMetrics(c *gin.Context) {
	var req k8s.MetricsForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	metricsList, err := models.GetPodMetrics(req.ClusterId, req.MetricName, models.ResourceSelector{
		Namespace:    req.Namespace,
		ResourceName: req.PodName,
	})
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: metricsList})
}
