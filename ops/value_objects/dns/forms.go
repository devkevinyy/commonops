package dns

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 3:02 PM
 * @Desc:
 */

type DnsDomainReq struct {
	PageNum        int    `json:"pageNum" form:"pageNum"`
	PageSize       int    `json:"pageSize" form:"pageSize"`
	CloudAccountId int    `json:"cloudAccountId" form:"cloudAccountId"`
	DomainName     string `json:"domainName" form:"domainName"`
	RR             string `json:"rr" form:"rr"`       // 主机记录
	RType          string `json:"rType" form:"rType"` // 记录类型
	RValue         string `json:"rValue" form:"rValue"`
	RecordId       string `json:"recordId" form:"recordId"`
	Status         string `json:"status" form:"status"`
}

type DnsDomainCnameReq struct {
	RR     string `json:"rr" form:"rr" binding:"required"`
	RValue string `json:"rValue" form:"rValue" binding:"required"`
	Sign   string `json:"sign" form:"sign" binding:"required"`
}

type DnsDomainCnameQueryReq struct {
	PageNum  int    `json:"pageNum" form:"pageNum" binding:"required"`
	PageSize int    `json:"pageSize" form:"pageSize" binding:"required"`
	RR       string `json:"rr" form:"rr" binding:"required"`
	Sign     string `json:"sign" form:"sign" binding:"required"`
}

type AliDnsListResp struct {
	TotalCount int64 `json:"TotalCount"`
}