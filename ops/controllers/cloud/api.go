package cloud

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/chujieyang/commonops/ops/forms"
	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/untils"
	"github.com/chujieyang/commonops/ops/untils/ws"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"k8s.io/client-go/tools/remotecommand"
)

type CloudAccountResp struct {
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	Accounts []models.CloudAccount `json:"accounts"`
}

type CloudServerResp struct {
	Total   int                `json:"total"`
	Page    int                `json:"page"`
	Servers []models.EcsDetail `json:"servers"`
}

type RdsResp struct {
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Rds   []models.RdsDetail `json:"rds"`
}

type KvResp struct {
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Kv    []models.KvDetail `json:"kv"`
}

type SlbResp struct {
	Total int                `json:"total"`
	Page  int                `json:"page"`
	Slb   []models.SlbDetail `json:"slb"`
}

type CloudOtherResResp struct {
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
	OtherRes []models.OtherDetail `json:"otherRes"`
}

/*
 [api get]: 获取云账号列表
*/
func IGetCloudAccounts(c *gin.Context) {
	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, sizeErr := strconv.Atoi(c.DefaultQuery("size", "10"))
	if pageErr != nil || sizeErr != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "page或size类型转换异常"})
		return
	}
	offset := (page - 1) * size
	total := models.GetCloudAccountsCount()
	accountList := models.GetCloudAccounts(offset, size)
	resp := CloudAccountResp{
		Total:    total,
		Page:     page,
		Accounts: accountList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: resp})
}

/*
 [api post]: 添加云账号信息
*/
func IPostCloudAccounts(c *gin.Context) {
	var query forms.CloudAccountForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.AddCloudAccount(query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

/*
 [api put]: 更新云账号信息
*/
func IPutCloudAccounts(c *gin.Context) {
	var query forms.CloudAccountForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.UpdateCloudAccount(query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

/*
 [api delete]: 删除云账号信息
*/
func IDeleteCloudAccounts(c *gin.Context) {
	var query forms.CloudAccountForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.DeleteCloudAccount(query.Id)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

/*
 [api get]: 获取公有云的ecs服务器列表
*/
func IGetCloudServers(c *gin.Context) {
	var req forms.CloudServerQueryForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.Size
	total := models.GetServersCount(untils.GetCurrentUserId(c), req.QueryExpiredTime, req.QueryKeyword,
		req.QueryCloudAccount)
	serverList := models.GetServers(untils.GetCurrentUserId(c), offset, req.Size,
		req.QueryExpiredTime, req.QueryKeyword, req.QueryCloudAccount)
	resp := CloudServerResp{
		Total:   total,
		Page:    req.Page,
		Servers: serverList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: resp})
}

/*
 [api post]: 新增服务器资源信息
*/
func IPostCloudServers(c *gin.Context) {
	var req forms.ServerInfoForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.AddCloudServer(req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

/*
 [api put]: 完善服务器的扩展信息
*/
func IPutCloudServers(c *gin.Context) {
	var req forms.ExtraInfoForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.UpdateCloudServer(req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "修改数据成功！"})
}

/*
 [api delete]: 删除服务器资源信息
*/
func IDeleteCloudServers(c *gin.Context) {
	var query forms.ResDeleteForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.DeleteCloudServer(query.Id)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

/*
 [api get]: 获取公有云的ecs服务器详情
*/
func IGetCloudServerDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "需要传入整型类型的服务器id"})
		return
	}
	server := models.GetServerDetail(id)
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: server})
}

/*
 [api get]: 获取公有云的rds详情
*/
func IGetCloudRdsDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "需要传入整型类型的服务器id"})
		return
	}
	server := models.GetRdsDetail(id)
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: server})
}

/*
 [api get]: 获取公有云的kv详情
*/
func IGetCloudKvDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "需要传入整型类型的服务器id"})
		return
	}
	server := models.GetKvDetail(id)
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: server})
}

/*
 [api get]: 获取公有云的slb详情
*/
func IGetCloudSlbDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: "需要传入整型类型的服务器id"})
		return
	}
	server := models.GetSlbDetail(id)
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: server})
}

/*
 [api get]: 获取公有云的 rds 列表
*/
func IGetCloudRds(c *gin.Context) {
	var req forms.RdsQueryForm
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.Size
	total := models.GetRdsCount(untils.GetCurrentUserId(c), req.QueryExpiredTime, req.QueryKeyword,
		req.QueryCloudAccount)
	rdsList := models.GetRdsByPage(untils.GetCurrentUserId(c), offset, req.Size, req.QueryExpiredTime, req.QueryKeyword,
		req.QueryCloudAccount)
	resp := RdsResp{
		Total: total,
		Page:  req.Page,
		Rds:   rdsList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: resp})
}

/*
 [api put]: 完善服务器的扩展信息
*/
func IPutCloudRds(c *gin.Context) {
	var req forms.ExtraInfoForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.UpdateCloudRds(req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "修改数据成功！"})
}

/*
 [api post]: 新增rds信息
*/
func IPostCloudRds(c *gin.Context) {
	var query forms.RdsInfoForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.AddCloudRds(query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

func IDeleteCloudRds(c *gin.Context) {
	var query forms.ResDeleteForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.DeleteCloudRds(query.Id)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

/*
 [api get]: 获取公有云的 kv 列表
*/
func IGetCloudKv(c *gin.Context) {
	var req forms.KvQueryForm
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.Size
	total := models.GetKvCount(untils.GetCurrentUserId(c), req.QueryExpiredTime, req.QueryKeyword,
		req.QueryCloudAccount)
	kvList := models.GetKvByPage(untils.GetCurrentUserId(c), offset, req.Size, req.QueryExpiredTime, req.QueryKeyword,
		req.QueryCloudAccount)
	resp := KvResp{
		Total: total,
		Page:  req.Page,
		Kv:    kvList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: resp})
}

/*
 [api put]: 完善服务器的扩展信息
*/
func IPutCloudKv(c *gin.Context) {
	var req forms.ExtraInfoForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.UpdateCloudKv(req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "修改数据成功！"})
}

/*
 [api post]: 新增redis信息
*/
func IPostCloudKv(c *gin.Context) {
	var query forms.KvInfoForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.AddCloudKv(query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

func IDeleteCloudKv(c *gin.Context) {
	var query forms.ResDeleteForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.DeleteCloudKv(query.Id)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success"})
}

/*
 [api get]: 获取公有云的 slb 列表
*/
func IGetCloudSlb(c *gin.Context) {
	var req forms.SlbQueryForm
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.Size
	total := models.GetSlbCount(untils.GetCurrentUserId(c), req.QueryKeyword,
		req.QueryCloudAccount)
	slbList := models.GetSlbByPage(untils.GetCurrentUserId(c), offset, req.Size, req.QueryKeyword,
		req.QueryCloudAccount)
	resp := SlbResp{
		Total: total,
		Page:  req.Page,
		Slb:   slbList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: resp})
}

func IDeleteCloudSlb(c *gin.Context) {
	var query forms.ResDeleteForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.DeleteCloudSlb(query.Id)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "删除slb成功"})
}

/*
 [api get]: 获取其他资源列表
*/
func IGetCloudOtherRes(c *gin.Context) {
	var req forms.CloudOtherResQueryForm
	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	offset := (req.Page - 1) * req.Size
	total := models.GetOtherResCount(untils.GetCurrentUserId(c), req.QueryKeyword, req.QueryResType,
		req.QueryCloudAccount, req.QueryExpiredTime, req.QueryManageUser)
	serverList := models.GetOtherRes(untils.GetCurrentUserId(c), offset, req.Size, req.QueryKeyword,
		req.QueryResType, req.QueryCloudAccount, req.QueryExpiredTime, req.QueryManageUser)
	resp := CloudOtherResResp{
		Total:    total,
		Page:     req.Page,
		OtherRes: serverList,
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: resp})
}

/*
 [api post]: 新增其他资源信息
*/
func IPostCloudOtherRes(c *gin.Context) {
	var req forms.AddOtherResForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.AddCloudOtherRes(req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
 [api delete]: 删除其他资源信息
*/
func IDeleteCloudOtherRes(c *gin.Context) {
	var query forms.ResDeleteForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.DeleteCloudOtherRes(query.Id)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
}

func WsSsh(c *gin.Context) {
	var query forms.SshForm
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	var wsConn *ws.WsConnection

	if wsConn, err = ws.InitWebsocket(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}

	ecsInfo, err := models.GetEcsSshInfo(query.ServerId)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}

	client, err := ws.NewSshClient(ecsInfo.PublicIpAddress, ecsInfo.SshPort, ecsInfo.SshUser, untils.DesDecode(ecsInfo.SshPwd))
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	rw := io.ReadWriter(&k8s_structs.StreamHandler{
		WsConn:      wsConn,
		ResizeEvent: make(chan remotecommand.TerminalSize)})
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("session error: ", err.Error())
		return
	}
	defer session.Close()
	session.Stdout = rw
	session.Stderr = rw
	session.Stdin = rw
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", 50, 120, modes)
	if err != nil {
		fmt.Println("d1 error: ", err.Error())
	}
	err = session.Shell()
	if err != nil {
		fmt.Println("d2 error: ", err.Error())
	}
	err = session.Wait()
	if err != nil {
		fmt.Println("d3 error: ", err.Error())
	}
	c.JSON(http.StatusOK, untils.RespData{Code: 0, Msg: "success", Data: nil})
	return
}
