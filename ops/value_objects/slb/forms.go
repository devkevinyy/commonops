package slb

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:38 AM
 * @Desc:
 */

type SlbResp struct {
	Total uint                `json:"total"`
	Page  uint                `json:"page"`
	Slb   interface{} `json:"slb"`
}

type SlbQueryForm struct {
	Page              uint    `form:"page" json:"page" binding:"required"`
	Size              uint    `form:"size" json:"size" binding:"required"`
	QueryKeyword      string `form:"queryKeyword" json:"queryKeyword"`
	QueryCloudAccount uint    `form:"queryCloudAccount" json:"queryCloudAccount"`
}

type ResDeleteForm struct {
	Id uint `json:"id"`
}