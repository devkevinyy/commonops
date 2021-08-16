package other

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 12:07 PM
 * @Desc:
 */

type CloudOtherResResp struct {
	Total    uint                  `json:"total"`
	Page     uint                  `json:"page"`
	OtherRes interface{} `json:"otherRes"`
}

type CloudOtherResQueryForm struct {
	Page              uint    `form:"page" json:"page" binding:"required"`
	Size              uint    `form:"size" json:"size" binding:"required"`
	QueryKeyword      string `form:"queryKeyword" json:"queryKeyword"`
	QueryResType      string `form:"queryResType" json:"queryResType"`
	QueryCloudAccount uint    `form:"queryCloudAccount" json:"queryCloudAccount"`
	QueryExpiredTime  string `form:"queryExpiredTime" json:"queryExpiredTime"`
	QueryManageUser   uint    `form:"queryManageUser" json:"queryManageUser"`
}

type AddOtherResForm struct {
	Id             int    `json:"id"`
	CloudAccountId int    `json:"cloudAccountId"`
	ResType        string `json:"resType"`
	InstanceId     string `json:"instanceId"`
	InstanceName   string `json:"instanceName"`
	Connections    string `json:"connections"`
	Region         string `json:"region"`
	Engine         string `json:"engine"`
	Cpu            string `json:"cpu"`
	Disk           string `json:"disk"`
	Bandwidth      string `json:"bandwidth"`
	Memory         string `json:"memory"`
	CreateTime     string `json:"createTime"`
	ExpiredTime    string `json:"expiredTime"`
	ManageUser     int    `json:"manageUser"`
	Cost           string `json:"cost"`
	CostType       string `json:"costType"`
}

type ResDeleteForm struct {
	Id uint `json:"id"`
}