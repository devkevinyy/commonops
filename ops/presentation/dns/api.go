package dns

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/value_objects/dns"
	"net/http"

	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
)

/*
 [api post]: 新增域名
*/
func IPostDnsDomain(c *gin.Context) {
	var req dns.DnsDomainReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	cloudAccount, err := models.GetCloudAccountInfo(req.CloudAccountId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	dnsClient, err := service.GetAliDnsService(cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	data, err := dnsClient.AddCloudAccountDomain(req.DomainName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
 [api get]: 获取域名列表
*/
func IGetDnsDomainList(c *gin.Context) {
	var req dns.DnsDomainReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	cloudAccount, err := models.GetCloudAccountInfo(req.CloudAccountId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	dnsClient, err := service.GetAliDnsService(cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	data, err := dnsClient.GetCloudAccountDomainList(req.DomainName, fmt.Sprintf("%d", req.PageNum), fmt.Sprintf("%d", req.PageSize))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
 [api get]: 获取域名解析历史列表
*/
func IGetDnsDomainHistoryList(c *gin.Context) {
	var req dns.DnsDomainReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	cloudAccount, err := models.GetCloudAccountInfo(req.CloudAccountId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	dnsClient, err := service.GetAliDnsService(cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	data, err := dnsClient.GetCloudAccountDomainHistoryList(req.DomainName, "1", "10")
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
 [api get]: 获取域名解析记录列表
*/
func IGetDnsDomainRecordsList(c *gin.Context) {
	var req dns.DnsDomainReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	cloudAccount, err := models.GetCloudAccountInfo(req.CloudAccountId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	dnsClient, err := service.GetAliDnsService(cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	data, err := dnsClient.GetCloudAccountDomainRecordsList(req.DomainName, req.RR,
		fmt.Sprintf("%d", req.PageNum), fmt.Sprintf("%d", req.PageSize))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
 [api post]: 新增域名解析记录
*/
func IPostDnsDomainRecord(c *gin.Context) {
	var req dns.DnsDomainReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	cloudAccount, err := models.GetCloudAccountInfo(req.CloudAccountId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	dnsClient, err := service.GetAliDnsService(cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	data, err := dnsClient.AddCloudAccountDomainRecord(req.DomainName, req.RR, req.RType, req.RValue)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
 [api post]: 修改域名解析记录
*/
func IUpdateDnsDomainRecord(c *gin.Context) {
	var req dns.DnsDomainReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	cloudAccount, err := models.GetCloudAccountInfo(req.CloudAccountId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	dnsClient, err := service.GetAliDnsService(cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	data, err := dnsClient.UpdateCloudAccountDomainRecord(req.RecordId, req.RR, req.RType, req.RValue)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
 [api delete]: 删除域名解析记录
*/
func IDeleteDnsDomainRecord(c *gin.Context) {
	var req dns.DnsDomainReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	cloudAccount, err := models.GetCloudAccountInfo(req.CloudAccountId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	dnsClient, err := service.GetAliDnsService(cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	data, err := dnsClient.DeleteCloudAccountDomainRecord(req.RecordId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
 [api post]: 设置域名解析记录状态
*/
func IPostDnsDomainRecordStatus(c *gin.Context) {
	var req dns.DnsDomainReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	cloudAccount, err := models.GetCloudAccountInfo(req.CloudAccountId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	dnsClient, err := service.GetAliDnsService(cloudAccount.Region, cloudAccount.Key, cloudAccount.Secret)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	data, err := dnsClient.SetCloudAccountDomainRecordStatus(req.RecordId, req.Status)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}
