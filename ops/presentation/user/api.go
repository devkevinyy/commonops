package user

import (
	"net/http"
	"strconv"

	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/opslog"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HC(c *gin.Context) {
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
 [api post]: 创建新用户
*/
func IPostUserCreate(c *gin.Context) {
	var req user.UserForm
	err := c.BindJSON(&req)
	if err != nil {
		opslog.Error().Printf("[IPostUserCreate]: %s, %s \n", "用户参数异常", err.Error())
		c.JSON(http.StatusBadRequest, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	id, err := models.AddUser(req.Username, req.Password, req.Email, req.Position)
	if err != nil {
		opslog.Error().Println("[IPostUserCreate]", zap.String("新增用户异常", err.Error()))
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	data := make(map[string]interface{})
	data["id"] = id
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
 [api post]: 修改用户状态
*/
func IPutUserActive(c *gin.Context) {
	var req user.UserStatusForm
	err := c.BindJSON(&req)
	if err != nil {
		opslog.Error().Println("[IPostUserCreate]", zap.String("用户参数异常", err.Error()))
		c.JSON(http.StatusBadRequest, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	err = models.UpdateUserActiveStatus(req.Email, req.Active)
	if err != nil {
		opslog.Error().Println("[IPostUserCreate]", zap.String("新增用户异常", err.Error()))
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
 [api get]: 获取所有用户列表
*/
func IGetUsersList(c *gin.Context) {
	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, sizeErr := strconv.Atoi(c.DefaultQuery("size", "10"))
	if pageErr != nil || sizeErr != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "page或size类型转换异常"})
		return
	}
	offset := (page - 1) * size
	total := models.GetUsersCount()
	userList := models.GetUserList(offset, size)
	resp := user.UsersListResp{
		Total: total,
		Page:  page,
		Users: userList,
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: resp})
}

/*
 [api get]: 获取所有角色列表
*/
func IGetRolesList(c *gin.Context) {
	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, sizeErr := strconv.Atoi(c.DefaultQuery("size", "10"))
	if pageErr != nil || sizeErr != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "page或size类型转换异常"})
		return
	}
	offset := (page - 1) * size
	total := models.GetRolesCount()
	userList := models.GetRolesList(offset, size)
	resp := user.RolesListResp{
		Total: total,
		Page:  page,
		Roles: userList,
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: resp})
}

/*
 [api get]: 获取所有权限菜单列表
*/
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

/*
	[api post]: 用户登录
*/
func ILogin(c *gin.Context) {
	var req user.LoginForm
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	userInfo, err := models.GetUserInfo(req.Username)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "用户不存在或被禁用"})
		return
	}
	inputPwd := utils.GenUserPassword(req.Password)
	if inputPwd != userInfo.Password {
		c.JSON(http.StatusOK, utils.RespData{Code: 1, Msg: "密码校验错误"})
		return
	}
	data := map[string]interface{}{
		"token": utils.GenJWT(map[string]interface{}{
			"userId":       userInfo.ID,
			"email":        userInfo.Email,
			"username":     userInfo.UserName,
			"position":     userInfo.Position,
			"show_intro":   userInfo.FirstIn,
			"isSuperAdmin": models.IsUserSuperAdmin(userInfo.ID),
			"empId":        userInfo.EmpId,
		}),
	}
	models.UpdateUserFirstInStatus(userInfo.ID)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
	[api get]: 用户jwt有效期刷新
*/
func ITokenRefresh(c *gin.Context) {
	email, _ := c.Get("email")
	userInfo, err := models.GetUserInfo(email.(string))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "用户不存在"})
		return
	}
	data := map[string]interface{}{
		"token": utils.GenJWT(map[string]interface{}{
			"userId":       userInfo.ID,
			"email":        userInfo.Email,
			"username":     userInfo.UserName,
			"position":     userInfo.Position,
			"show_intro":   userInfo.FirstIn,
			"isSuperAdmin": models.IsUserSuperAdmin(userInfo.ID),
			"empId":        userInfo.EmpId,
		}),
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
	[api post]: 用户修改密码
*/
func IUpdatePassword(c *gin.Context) {
	empId, _ := c.Get("empId")
	var req user.UpdatePasswordForm
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "二次确认密码不一致", Data: nil})
		return
	}
	userInfo, err := models.UpdateUserPassword(empId.(string), utils.GenUserPassword(req.Password))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: userInfo})
}

/*
	[api post]: 创建新角色
*/
func ICreateRole(c *gin.Context) {
	var req map[string]string
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.AddNewRole(req["roleName"], req["description"])
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
	[api update]: 更新角色
*/
func IUpdateRole(c *gin.Context) {
	var req map[string]interface{}
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.UpdateRole(int(req["id"].(float64)), req["roleName"].(string), req["description"].(string))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
	[api delete]: 删除角色
*/
func IDeleteRole(c *gin.Context) {
	var req map[string]int
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.DeleteRole(req["id"])
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
	[api get]: 获取用户允许访问的权限列表
*/
func IGetUserPermissions(c *gin.Context) {
	userId, _ := c.Get("userId")
	authType := c.Query("authType")
	permissionList := models.GetUserPermissions(int(userId.(float64)), authType)
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: permissionList})
}

/*
	[api get]: 获取角色组下的所有用户
*/
func IGetRoleUserList(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Query("roleId"))
	inUser, err := models.GetUsersByRoleId(roleId)
	allUser, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	data := map[string]interface{}{
		"in":  inUser,
		"all": allUser,
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
	[api post]: 管理角色下面的用户(增加/删除)
*/
func IPostRoleUserList(c *gin.Context) {
	var req user.AddRoleUsers
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.AddRoleUsers(req.RoleId, req.UserIdList)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
	[api get]: 获取角色组下的各种资源（ecs, rds, kv, slb）
*/
func IGetRoleResourceList(c *gin.Context) {
	resourceType := c.Query("resourceType")
	roleId, err := strconv.Atoi(c.Query("roleId"))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	data := map[string]interface{}{}
	switch resourceType {
	case "ecs":
		inData, err1 := models.GetEcsListByRoleId(roleId)
		allData, err2 := models.GetAllEcsList()
		if err1 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err1.Error(), Data: nil})
			return
		}
		if err2 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err2.Error(), Data: nil})
			return
		}
		data["in"] = inData
		data["all"] = allData
	case "rds":
		inData, err1 := models.GetRdsListByRoleId(roleId)
		allData, err2 := models.GetAllRdsList()
		if err1 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err1.Error(), Data: nil})
			return
		}
		if err2 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err2.Error(), Data: nil})
			return
		}
		data["in"] = inData
		data["all"] = allData
	case "kv":
		inData, err1 := models.GetKvListByRoleId(roleId)
		allData, err2 := models.GetAllKvList()
		if err1 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err1.Error(), Data: nil})
			return
		}
		if err2 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err2.Error(), Data: nil})
			return
		}
		data["in"] = inData
		data["all"] = allData
	case "slb":
		inData, err1 := models.GetSlbListByRoleId(roleId)
		allData, err2 := models.GetAllSlbList()
		if err1 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err1.Error(), Data: nil})
			return
		}
		if err2 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err2.Error(), Data: nil})
			return
		}
		data["in"] = inData
		data["all"] = allData
	case "other":
		inData, err1 := models.GetOtherResListByRoleId(roleId)
		allData, err2 := models.GetAllOtherResList()
		if err1 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err1.Error(), Data: nil})
			return
		}
		if err2 != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err2.Error(), Data: nil})
			return
		}
		data["in"] = inData
		data["all"] = allData
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
	[api post]: 管理角色下面的ECS资源(增加/删除)
*/
func IPostRoleResourcesList(c *gin.Context) {
	var req user.AddRoleResources
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.AddRoleEcs(req.RoleId, req.EcsIdList)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.AddRoleRds(req.RoleId, req.RdsIdList)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.AddRoleKv(req.RoleId, req.KvIdList)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.AddRoleSlb(req.RoleId, req.SlbIdList)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.AddRoleOtherRes(req.RoleId, req.OtherIdList)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
	[api get]: 获取角色组下的所有权限链接
*/
func IGetRoleAuthLinkList(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Query("roleId"))
	inData, err := models.GetAuthLinksByRoleId(roleId)
	allData, err := models.GetAllAuthLinks()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	data := map[string]interface{}{
		"in":  inData,
		"all": allData,
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

/*
	[api post]: 创建权限链接
*/
func ICreateAuthLink(c *gin.Context) {
	var req map[string]string
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	err = models.AddAuthLink(req["name"], req["description"], req["path"])
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

/*
	[api delete]: 删除权限链接
*/
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

/*
	[api post]: 管理角色下面的权限链接(增加/删除)
*/
func ICreateRoleAuthLink(c *gin.Context) {
	var req user.AddRoleAuthLinks
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = models.AddRoleAuthLinks(req.RoleId, req.AuthLinkIdList)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetUserFeedback(c *gin.Context) {
	page, pageErr := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, sizeErr := strconv.Atoi(c.DefaultQuery("size", "10"))
	if pageErr != nil || sizeErr != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: "page或size类型转换异常"})
		return
	}
	offset := (page - 1) * size
	total := models.GetUserFeedbackCount()
	userFeedbackList := models.GetUserFeedbackList(offset, size)
	resp := user.FeedbackListResp{
		Total:     total,
		Page:      page,
		Feedbacks: userFeedbackList,
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: resp})
}

func IPostUserFeedback(c *gin.Context) {
	var req user.FeedbackForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	userId := utils.GetCurrentUserId(c)
	username := utils.GetCurrentUsername(c)
	err = models.SaveUserFeedback(userId, username, req.Advice, req.Score)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "提交成功", Data: nil})
}
