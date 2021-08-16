package daily_job

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 2:55 PM
 * @Desc:
 */

type AddDailyJobForm struct {
	JobName              string `json:"jobName"`
	JobType              string `json:"jobType"`
	ImportantDegree      string `json:"importantDegree"`
	OpenDeployAutoConfig string `json:"open_deploy_auto_config"`
	Remark               string `json:"remark"`
	TaskContent          string `json:"task_content"`
	CreatorUserId        int    `json:"creator_user_id"`
	CreatorUserName      string `json:"creator_user_name"`
}

type UpdateDailyJobForm struct {
	Id           int    `json:"id"`
	Action       string `json:"action"`
	UserId       int    `json:"user_id"`
	UserName     string `json:"user_name"`
	RefuseReason string `json:"refuseReason"`
}

type UpdateDailyJobExecutorUserForm struct {
	JobId        int    `json:"jobId"`
	ChangeUserId string `json:"changeUserId"`
}

type JobsQueryForm struct {
	Page            uint   `form:"page" json:"page" binding:"required"`
	Size            uint   `form:"size" json:"size" binding:"required"`
	QueryKeyword    string `form:"queryKeyword" json:"queryKeyword"`
	QueryCreateTime string `form:"queryCreateTime" json:"queryCreateTime"`
}

type DailyJobResp struct {
	Total uint              `json:"total"`
	Page  uint              `json:"page"`
	Jobs  interface{} `json:"jobs"`
}

