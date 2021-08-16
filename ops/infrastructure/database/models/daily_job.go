package models

import (
	"errors"
	"fmt"
	"github.com/chujieyang/commonops/ops/opslog"
	daily_job2 "github.com/chujieyang/commonops/ops/value_objects/daily_job"
	"strings"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/chujieyang/commonops/ops/utils"
)

type DailyJob struct {
	Id                   int             `json:"id" gorm:"column:id;PRIMARY_KEY"`
	JobName              string          `json:"job_name" gorm:"column:job_name;type:varchar(512)"`
	JobType              string          `json:"job_type" gorm:"column:job_type"`
	ImportantDegree      string          `json:"important_degree" gorm:"column:important_degree"`
	OpenDeployAutoConfig string          `json:"open_deploy_auto_config" gorm:"column:open_deploy_auto_config;type:text"`
	TaskContent          string          `json:"task_content" gorm:"column:task_content;type:text"`
	Remark               string          `json:"remark" gorm:"column:remark;type:text"`
	CreatorUserId    int            `json:"creator_user_id" gorm:"column:creator_user_id;type:int"`
	CreatorUserName  string         `json:"creator_user_name" gorm:"column:creator_user_name"`
	ExecutorUserId   int            `json:"executor_user_id" gorm:"column:executor_user_id;type:int"`
	ExecutorUserName string         `json:"executor_user_name" gorm:"column:executor_user_name"`
	RefuseReason     string         `json:"refuse_reason" gorm:"column:refuse_reason;type:text"`
	Status           int8           `json:"status" gorm:"column:status"`
	CreateTime       utils.JSONTime `json:"create_time" gorm:"column:create_time"`
	AcceptTime       utils.JSONTime `json:"accept_time" gorm:"column:accept_time;default:null"`
	EndTime          utils.JSONTime `json:"end_time" gorm:"column:end_time;default:null"`
}

func AddDailyJob(jobForm daily_job2.AddDailyJobForm) (err error) {
	dailyJob := DailyJob{
		JobName:              jobForm.JobName,
		JobType:              jobForm.JobType,
		ImportantDegree:      jobForm.ImportantDegree,
		OpenDeployAutoConfig: jobForm.OpenDeployAutoConfig,
		Remark:               jobForm.Remark,
		TaskContent:          jobForm.TaskContent,
		CreatorUserId:        jobForm.CreatorUserId,
		CreatorUserName:      jobForm.CreatorUserName,
		Status:               1,
		CreateTime:           utils.GetNowTime(),
	}
	tx := database.Mysql().Begin()
	err = tx.Create(&dailyJob).Error
	if err != nil {
		tx.Rollback()
		opslog.Error().Println(err)
		return
	}
	tx.Commit()
	return
}

func UpdateDailyJob(jobForm daily_job2.UpdateDailyJobForm) (err error) {
	dailyJob := DailyJob{
		ExecutorUserId:   jobForm.UserId,
		ExecutorUserName: jobForm.UserName,
	}
	switch jobForm.Action {
	case "getJob":
		dailyJob.Status = int8(2)
		dailyJob.AcceptTime = utils.GetNowTime()
		total := 0
		database.Mysql().Model(&dailyJob).Where("status = 1 and id = ?", jobForm.Id).Count(&total)
		if total != 1 {
			return errors.New("该任务处于不能被领取的状态")
		}
		break
	case "finishJob":
		var jobInfo DailyJob
		dailyJob.Status = int8(3)
		dailyJob.EndTime = utils.GetNowTime()
		database.Mysql().Model(&dailyJob).Where("status = 2 and id = ?", jobForm.Id).Find(&jobInfo)
		if jobInfo.Id == 0 {
			return errors.New("该任务处于不能被设置为完成的状态")
		}
		dailyJob.AcceptTime = jobInfo.AcceptTime
		break
	case "deleteJob":
		dailyJob.Status = int8(-1)
		total := 0
		database.Mysql().Model(&dailyJob).Where("creator_user_id = ? and id = ?",
			jobForm.UserId, jobForm.Id).Count(&total)
		if total != 1 {
			return errors.New("用户只能删除自己创建的任务")
		}
		break
	case "refuseJob":
		dailyJob.Status = int8(0)
		dailyJob.RefuseReason = jobForm.RefuseReason
		total := 0
		database.Mysql().Model(&dailyJob).Where("status = 1 and id = ?",
			jobForm.Id).Count(&total)
		if total != 1 {
			return errors.New("该任务处于不能被拒绝的状态")
		}
		break
	default:
		return errors.New("不支持的操作")
	}

	err = database.Mysql().Exec("update daily_jobs set status = ?, executor_user_id = ?, "+
		"executor_user_name = ?, accept_time = ?, end_time = ?, refuse_reason = ? where id=?",
		dailyJob.Status, dailyJob.ExecutorUserId, dailyJob.ExecutorUserName, dailyJob.AcceptTime,
		dailyJob.EndTime, dailyJob.RefuseReason, jobForm.Id).Error
	return
}

func UpdateDailyJobExecutorUser(data daily_job2.UpdateDailyJobExecutorUserForm) (err error) {
	infoList := strings.Split(data.ChangeUserId, "-")
	var jobInfo DailyJob
	err = database.Mysql().Raw("select * from daily_jobs where id= ? ",
		data.JobId).Scan(&jobInfo).Error
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	updateSql := "update daily_jobs set executor_user_id = ?, executor_user_name = ? where id = ?"
	err = database.Mysql().Exec(updateSql, infoList[0], infoList[1], data.JobId).Error
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

func GetJobsCount(userId uint, queryKeyword string, queryCreateTime string) (total uint) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	total = 0
	var querySql = "select count(*) from daily_jobs where status >= 0 "
	var args []interface{}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and ( job_name like ? or creator_user_name like ? "+
			"or executor_user_name like ?) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCreateTime != "" {
		querySql = fmt.Sprintf("%s and create_time between ? and ? ", querySql)
		args = append(args, fmt.Sprintf("%s 00:00:00", queryCreateTime), fmt.Sprintf("%s 23:59:59", queryCreateTime))
	}
	if !isSuperAdminOrOps {
		querySql = fmt.Sprintf("%s and creator_user_id = ? ", querySql)
		args = append(args, userId)
	}
	database.Mysql().Raw(querySql, args...).Count(&total)
	return
}

func GetJobsByPage(userId uint, queryKeyword string, queryCreateTime string, offset uint, limit uint) (jobs []DailyJob) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	var querySql = "select * from daily_jobs where status >= 0 "
	var args []interface{}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (job_name like ? or creator_user_name like ? "+
			"or executor_user_name like ? ) ", querySql)
		args = append(args, "%"+queryKeyword+"%", "%"+queryKeyword+"%", "%"+queryKeyword+"%")
	}
	if queryCreateTime != "" {
		querySql = fmt.Sprintf("%s and create_time between ? and ? ", querySql)
		args = append(args, fmt.Sprintf("%s 00:00:00", queryCreateTime), fmt.Sprintf("%s 23:59:59", queryCreateTime))
	}
	if !isSuperAdminOrOps {
		querySql = fmt.Sprintf("%s and creator_user_id = ? ", querySql)
		args = append(args, userId)
	}
	querySql += " order by status asc, id desc limit ?, ? "
	args = append(args, offset, limit)
	database.Mysql().Raw(querySql, args...).Scan(&jobs)
	return
}

func GetJobDetail(id int) (job DailyJob) {
	database.Mysql().Find(&job, id)
	return
}
