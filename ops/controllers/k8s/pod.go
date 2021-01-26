package k8s

import (
	"encoding/json"
	"net/http"

	"github.com/chujieyang/commonops/ops/database"
	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/chujieyang/commonops/ops/untils/k8s_utils"
	"github.com/chujieyang/commonops/ops/untils/ws"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/tools/remotecommand"
)

func GetPods(c *gin.Context) {
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
	data, err := cluster.GetPods(req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}

func GetPodLogs(c *gin.Context) {
	var req k8s_structs.PodContainerLogForm
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
	data, err := cluster.GetPodContainerLogs(req.Namespace, req.PodName, req.ContainerName)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(data))
}

func GetContainerTerminal(c *gin.Context) {
	var (
		wsConn *ws.WsConnection
		err    error
	)

	if wsConn, err = ws.InitWebsocket(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}

	msg, err := wsConn.WsRead()
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}

	inputMap := map[string]interface{}{}
	msgMap := map[string]interface{}{}
	err = json.Unmarshal(msg.Data, &msgMap)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}

	if msgMap["type"].(string) == "input" {
		inputData := msgMap["input"].(string)
		err = json.Unmarshal([]byte(inputData), &inputMap)
		if err != nil {
			c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
			return
		}
		if inputMap["action"] == "init_connection" {
			handler := &k8s_structs.StreamHandler{
				WsConn:      wsConn,
				ResizeEvent: make(chan remotecommand.TerminalSize),
			}
			handler.WsConn.InitLangEnv()

			k8s_utils.RemoteCommandContainerExec(inputMap["clusterId"].(string), inputMap["namespace"].(string),
				inputMap["podName"].(string), inputMap["containerName"].(string), handler)
			return
		}
	}
	c.JSON(http.StatusOK, k8s_structs.RespError(-1, "ws error"))
}

func GetPodMetrics(c *gin.Context) {
	var req k8s_structs.MetricsForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	metricsList, err := database.GetPodMetrics(req.ClusterId, req.MetricName, database.ResourceSelector{
		Namespace:    req.Namespace,
		ResourceName: req.PodName,
	})
	if err != nil {
		c.JSON(http.StatusOK, k8s_structs.RespError(-1, err.Error()))
		return
	}
	c.JSON(http.StatusOK, k8s_structs.RespSuccess(metricsList))
}
