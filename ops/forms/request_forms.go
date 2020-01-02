package forms



type UserForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email string `form:"email" json:"email"`
	Position string `form:"position" json:"position"`
}

type UserStatusForm struct {
	Email string `form:"email" json:"email"`
	Active bool `form:"active" json:"active"`
}

type LoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type CloudServerQueryForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
	QueryExpiredTime string `form:"queryExpiredTime" json:"queryExpiredTime"`
	QueryKeyword string `form:"queryKeyword" json:"queryKeyword"`
	QueryCloudAccount int `form:"queryCloudAccount" json:"queryCloudAccount"`
}

type RdsQueryForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
	QueryExpiredTime string `form:"queryExpiredTime" json:"queryExpiredTime"`
	QueryKeyword string `form:"queryKeyword" json:"queryKeyword"`
	QueryCloudAccount int `form:"queryCloudAccount" json:"queryCloudAccount"`
}

type KvQueryForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
	QueryExpiredTime string `form:"queryExpiredTime" json:"queryExpiredTime"`
	QueryKeyword string `form:"queryKeyword" json:"queryKeyword"`
	QueryCloudAccount int `form:"queryCloudAccount" json:"queryCloudAccount"`
}

type SlbQueryForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
	QueryKeyword string `form:"queryKeyword" json:"queryKeyword"`
	QueryCloudAccount int `form:"queryCloudAccount" json:"queryCloudAccount"`
}

type JobsQueryForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
	QueryKeyword string `form:"queryKeyword" json:"queryKeyword"`
	QueryCreateTime string `form:"queryCreateTime" json:"queryCreateTime"`
}

type CloudAccountForm struct {
	Id int `form:"id" json:"id"`
	AccountType string `form:"accountType" json:"accountType"`
	AccountName string `form:"accountName" json:"accountName"`
	AccountPwd string `form:"accountPwd" json:"accountPwd"`
	AccountKey string `form:"accountKey" json:"accountKey"`
	AccountSecret string `form:"accountSecret" json:"accountSecret"`
	AccountRegion string `form:"accountRegion" json:"accountRegion"`
	BankAccount int `form:"bankAccount" json:"bankAccount"`
}

type UpdatePasswordForm struct {
	Password string	`form:"password" json:"password" binding:"required"`
	ConfirmPassword string	`form:"confirm_password" json:"confirm_password" binding:"required"`
}

type AddRoleUsers struct {
	RoleId int64 `json:"roleId" binding:"required"`
	UserIdList []string `json:"userIdList"`
}

type AddRoleAuthLinks struct {
	RoleId int64 `json:"roleId" binding:"required"`
	AuthLinkIdList []string `json:"authLinkIdList"`
}

type AddRoleResources struct {
	RoleId int64 `json:"roleId" binding:"required"`
	EcsIdList []string `json:"ecsIdList"`
	RdsIdList []string `json:"rdsIdList"`
	KvIdList []string `json:"kvIdList"`
	SlbIdList []string `json:"slbIdList"`
	OtherIdList []string `json:"otherIdList"`
}

type AddDailyJobForm struct {
	JobName string 	`json:"jobName"`
	JobType string 	`json:"jobType"`
	ImportantDegree string 	`json:"importantDegree"`
	OpenDeployAutoConfig string `json:"open_deploy_auto_config"`
	Remark string 		`json:"remark"`
	TaskContent string 		`json:"task_content"`
	CreatorUserId int 		`json:"creator_user_id"`
	CreatorUserName string 		`json:"creator_user_name"`
}

type UpdateDailyJobForm struct {
	Id int `json:"id"`
	Action string `json:"action"`
	UserId int 		`json:"user_id"`
	UserName string  `json:"user_name"`
	RefuseReason string `json:"refuseReason"`
}

type UpdateDailyJobExecutorUserForm struct {
	JobId int `json:"jobId"`
	ChangeUserId string `json:"changeUserId"`
}

type ServerInfoForm struct {
	HostName string `json:"hostName"`
	InstanceId string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	Description string `json:"description"`
	InnerIpAddress string `json:"innerIpAddress"`
	PublicIpAddress string `json:"publicIpAddress"`
	Memory int `json:"memory"`
	Disk int `json:"disk"`
	CreateTime string `json:"createTime"`
	ExpiredTime string `json:"expiredTime"`
	CloudAccountId int `json:"cloudAccountId"`
	OsType string `json:"osType"`
	Cpu int `json:"cpu"`
}

type EcsBaseInfoForm struct {
	InstanceId string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	InnerIpAddress string `json:"innerIpAddress"`
	PublicIpAddress string `json:"publicIpAddress"`
	PrivateIpAddress string `json:"privateIpAddress"`
	Cpu string `json:"cpu"`
	Memory string `json:"memory"`
	ExpiredTime string `json:"expiredTime"`
}

type RdsBaseInfoForm struct {
	DbInstanceDescription string `json:"dbInstanceDescription"`
	DbInstanceStorage string `json:"dbInstanceStorage"`
	DbMemory string `json:"dbMemory"`
	DbExpiredTime string `json:"dbExpiredTime"`
}

type KvBaseInfoForm struct {
	KvInstanceName string `json:"kvInstanceName"`
	KvBandwidth string `json:"kvBandwidth"`
	KvCapacity string `json:"kvCapacity"`
	KvExpiredTime string `json:"kvExpiredTime"`
}

type ExtraInfoForm struct {
	Id string `json:"id"`
	EcsBaseInfoForm
	RdsBaseInfoForm
	KvBaseInfoForm
}

type RdsInfoForm struct {
	DBInstanceId string `json:"dbInstanceId"`
	DBInstanceDescription string `json:"dbInstanceDescription"`
	Engine string `json:"engine"`
	ConnectionString string `json:"connectionString"`
	Port string `json:"port"`
	DBInstanceMemory int `json:"dbInstanceMemory"`
	DBInstanceStorage int `json:"dbInstanceStorage"`
	CreateTime string `json:"createTime"`
	ExpireTime string `json:"expireTime"`
	CloudAccountId int `json:"cloudAccountId"`
}

type KvInfoForm struct {
	InstanceId string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	InstanceType string `json:"instanceType"`
	ConnectionString string `json:"connectionString"`
	Port string `json:"port"`
	Capacity int `json:"capacity"`
	Bandwidth int `json:"bandwidth"`
	CreateTime string `json:"createTime"`
	EndTime string `json:"endTime"`
	CloudAccountId int `json:"cloudAccountId"`
}

type ResDeleteForm struct {
	Id int `json:"id"`
}

type SiteInfoQueryForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
	QueryKeyword string `form:"queryKeyword" json:"queryKeyword"`
}

type SiteInfoForm struct {
	Id int `json:"id"`
	Organize string `json:"organize"`
	DeptName string `json:"deptName"`
	SiteId string  `json:"siteId"`
}

type CloudOtherResQueryForm struct {
	Page int `form:"page" json:"page" binding:"required"`
	Size int `form:"size" json:"size" binding:"required"`
	QueryKeyword string `form:"queryKeyword" json:"queryKeyword"`
	QueryResType string `form:"queryResType" json:"queryResType"`
	QueryCloudAccount int `form:"queryCloudAccount" json:"queryCloudAccount"`
	QueryExpiredTime string `form:"queryExpiredTime" json:"queryExpiredTime"`
	QueryManageUser int `form:"queryManageUser" json:"queryManageUser"`
}

type AddOtherResForm struct {
	Id                       int		`json:"id"`
	CloudAccountId           int		`json:"cloudAccountId"`
	ResType                  string     `json:"resType"`
	InstanceId               string     `json:"instanceId"`
	InstanceName             string     `json:"instanceName"`
	Connections              string     `json:"connections"`
	Region                   string     `json:"region"`
	Engine                   string     `json:"engine"`
	Cpu                      string     `json:"cpu"`
	Disk                     string     `json:"disk"`
	Bandwidth                string     `json:"bandwidth"`
	Memory                   string     `json:"memory"`
	CreateTime               string     `json:"createTime"`
	ExpiredTime              string     `json:"expiredTime"`
	ManageUser               int        `json:"manageUser"`
	Cost                     string     `json:"cost"`
	CostType                 string     `json:"costType"`
}

type ExpiredResQuery struct {
	StartDate string `form:"startDate" json:"startDate"`
	EndDate string `form:"endDate" json:"endDate"`
	ShowAll string   `form:"showAll" json:"showAll"`
	ResType string   `form:"resType" json:"resType"`
}

type ResExpiredQueryForm struct {
	StartDate string `form:"startDate" json:"startDate"`
	EndDate string `form:"endDate" json:"endDate"`
}

type FeedbackForm struct {
	Advice string `form:"advice" json:"advice"`
	Score int `form:"score" json:"score"`
}

type JobBuildForm struct {
	Name string `json:"name" form:"name"`
	BuildId int `json:"buildId" form:"buildId"`
	Start int `json:"start" form:"start"`
}

type LogQueryFrom struct {
	Project string `json:"project" form:"project"`
	LogStore string `json:"logStore" form:"logStore"`
	QueryService string `json:"queryService" form:"queryService"`
	QueryOperation string `json:"queryOperation" form:"queryOperation"`
	QueryTag string `json:"queryTag" form:"queryTag"`
	QueryExp string `json:"queryExp" form:"queryExp"`
	From int64 `json:"form" form:"from"`
	To int64 `json:"to" form:"to"`
	Line int64 `json:"line" form:"line"`
	Offset int64 `json:"offset" form:"offset"`
}

type DistributeTraceForm struct {
	Project string `json:"project" form:"project"`
	LogStore string `json:"logStore" form:"logStore"`
	TraceId string `json:"traceId" form:"traceId"`
	From int64 `json:"form" form:"from"`
	To int64 `json:"to" form:"to"`
}

type K8sForm struct {
	Id    int		`form:"id" json:"id"`
	ClusterId string `form:"clusterId" json:"clusterId"`
	Name string	    `form:"name" json:"name"`
	Description string  `form:"description" json:"description"`
	KubeConfigFilePath string   `form:"kubeconfig_file_path" json:"kubeconfig_file_path"`
}