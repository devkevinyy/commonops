package ecs

import (
	"io"
	"net/http"
	"strconv"

	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	ecs2 "github.com/chujieyang/commonops/ops/value_objects/ecs"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"k8s.io/client-go/tools/remotecommand"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/11/21 5:14 PM
 * @Desc:
 */

/*
 [api get]: 获取公有云的ecs服务器列表
*/
func IGetCloudServers(c *gin.Context) {
	var req ecs2.CloudServerQueryForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	req.UserId = utils.GetCurrentUserId(c)
	total, serverList := service.GetEcsService().GetEcsDataByPage(req)
	resp := ecs2.CloudServerResp{
		Total:   total,
		Page:    req.Page,
		Servers: serverList,
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: resp})
}

/*
 [api post]: 新增服务器资源信息
*/
func IPostCloudServers(c *gin.Context) {
	var req ecs2.ServerInfoForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetEcsService().AddCloudServer(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}

/*
 [api post]: 批量执行SSH操作
*/
func IPostCloudServersBatchSSH(c *gin.Context) {
	var req ecs2.BatchSshFrom
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	result, err := service.GetEcsService().BatchSshExec(req.Ids, req.Command)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: result})
}

/*
 [api put]: 完善服务器的扩展信息
*/
func IPutCloudServers(c *gin.Context) {
	var req ecs2.ExtraInfoForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err := service.GetEcsService().UpdateCloudServer(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "修改数据成功！"})
}

/*
 [api delete]: 删除服务器资源信息
*/
func IDeleteCloudServers(c *gin.Context) {
	var query ecs2.ResDeleteForm
	if err := c.Bind(&query); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if err := service.GetEcsService().DeleteCloudServer(query.Id); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success"})
}

/*
 [api get]: 获取公有云的ecs服务器详情
*/
func IGetCloudServerDetail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	server := service.GetEcsService().GetServerDetail(uint(id))
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: server})
}

func WsSsh(c *gin.Context) {
	var err error
	var wsConn *service.WsConnection
	var query ecs2.SshForm
	if err = c.Bind(&query); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if wsConn, err = service.NewWebSocketService(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}

	ecsInfo, err := models.GetEcsSshInfo(query.ServerId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}

	client, err := service.NewSshClient(ecsInfo.PublicIpAddress, ecsInfo.SshPort, ecsInfo.SshUser, utils.DesDecode(ecsInfo.SshPwd))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	rw := io.ReadWriter(&service.StreamHandler{
		WsConn:      wsConn,
		ResizeEvent: make(chan remotecommand.TerminalSize)})
	session, err := client.NewSession()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	defer func() {
		if err := session.Close(); err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
	}()
	session.Stdout = rw
	session.Stderr = rw
	session.Stdin = rw
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("xterm", 50, 120, modes); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = session.Shell(); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = session.Wait(); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
	return
}

/*
 [api get]: 获取所有服务器资源的树结构数据
*/
func IGetCloudServersTreeData(c *gin.Context) {
	data, err := service.GetEcsService().GetServersTreeData()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}
