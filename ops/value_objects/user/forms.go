package user

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 3:06 PM
 * @Desc:
 */

type UsersListResp struct {
	Total int 						`json:"total"`
	Page int 						`json:"page"`
	Users interface{}			`json:"users"`
}

type RolesListResp struct {
	Total int 						`json:"total"`
	Page int 						`json:"page"`
	Roles interface{}			`json:"roles"`
}

type PermissionsListResp struct {
	Total int 							`json:"total"`
	Page int 							`json:"page"`
	Permissions interface{}	`json:"permissions"`
}

type FeedbackListResp struct {
	Total int 							`json:"total"`
	Page int 							`json:"page"`
	Feedbacks interface{}	    `json:"feedbacks"`
}

type UserForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email"`
	Position string `form:"position" json:"position"`
}

type UserStatusForm struct {
	Email  string `form:"email" json:"email"`
	Active bool   `form:"active" json:"active"`
}

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UpdatePasswordForm struct {
	Password        string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required"`
}

type AddRoleResources struct {
	RoleId      int64    `json:"roleId" binding:"required"`
	EcsIdList   []string `json:"ecsIdList"`
	RdsIdList   []string `json:"rdsIdList"`
	KvIdList    []string `json:"kvIdList"`
	SlbIdList   []string `json:"slbIdList"`
	OtherIdList []string `json:"otherIdList"`
}

type AddRoleAuthLinks struct {
	RoleId         int64    `json:"roleId" binding:"required"`
	AuthLinkIdList []string `json:"authLinkIdList"`
}

type FeedbackForm struct {
	Advice string `form:"advice" json:"advice"`
	Score  int    `form:"score" json:"score"`
}

type AddRoleUsers struct {
	RoleId     int64    `json:"roleId" binding:"required"`
	UserIdList []string `json:"userIdList"`
}

