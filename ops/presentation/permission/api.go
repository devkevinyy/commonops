package permission

import (
	"net/http"
	"strconv"

	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/permission"
	"github.com/chujieyang/commonops/ops/value_objects/user"
	"github.com/gin-gonic/gin"
)

// IGetPermissionsList gets all auth links
// [api get]: 获取所有权限菜单列表
func IGetPermissionsList(c *gin.Context) {
	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, sizeErr := strconv.Atoi(c.DefaultQuery("size", "10"))
	if pageErr != nil || sizeErr != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "page或size类型转换异常"})
		return
	}
	offset := (page - 1) * size
	total := models.GetPermissionsCount()
	permissionsList := models.GetPermissionsList(offset, size)
	resp := user.PermissionsListResp{
		Total:       total,
		Page:        page,
		Permissions: permissionsList,
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: resp})
}

// ICreateAuthLink creates a new auth link
// [api post]: 创建权限链接
func ICreateAuthLink(c *gin.Context) {
	var req map[string]string
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.AddAuthLink(req["name"], req["description"], req["urlPath"], req["authType"])
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

// IDeleteAuthLink deletes auth link
// [api delete]: 删除权限链接
func IDeleteAuthLink(c *gin.Context) {
	var req map[string]int
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.DeleteAuthLink(req["id"])
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

// IGetAuthLink gets detail of authLink
// [api get]: 获取权限链接详情
func IGetAuthLink(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	authLink, err := models.GetAuthLink(id)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: authLink})
}

// IPutAuthLink puts modify data of auth link
// [api put]: 修改权限链接详情
func IPutAuthLink(c *gin.Context) {
	var req permission.AuthLinkInfoForm
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err := models.UpdateAuthLink(req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "修改数据成功!"})
}
