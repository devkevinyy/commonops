package models

import (
	"errors"
	"fmt"
	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/forms"
	"github.com/chujieyang/commonops/ops/untils"
	"strings"
)

type DailyJob struct {
	Id int           `json:"id" gorm:"column:id;PRIMARY_KEY"`
	JobName string 		`json:"job_name" gorm:"column:job_name;type:varchar(512)"`
	JobType string 		`json:"job_type" gorm:"column:job_type"`
	ImportantDegree string 		`json:"important_degree" gorm:"column:important_degree"`
	OpenDeployAutoConfig string  `json:"open_deploy_auto_config" gorm:"column:open_deploy_auto_config;type:text"`
	TaskContent string  `json:"task_content" gorm:"column:task_content;type:text"`
	Remark string 		`json:"remark" gorm:"column:remark;type:text"`
	CreatorUserId int `json:"creator_user_id" gorm:"column:creator_user_id;type:int"`
	CreatorUserName string `json:"creator_user_name" gorm:"column:creator_user_name"`
	ExecutorUserId int `json:"executor_user_id" gorm:"column:executor_user_id;type:int"`
	ExecutorUserName string `json:"executor_user_name" gorm:"column:executor_user_name"`
	Status int8  `json:"status" gorm:"column:status"`
	CreateTime untils.JSONTime `json:"create_time" gorm:"column:create_time"`
	AcceptTime untils.JSONTime `json:"accept_time" gorm:"column:accept_time;default:null"`
	EndTime untils.JSONTime `json:"end_time" gorm:"column:end_time;default:null"`
}

func AddDailyJob(jobForm forms.AddDailyJobForm) (err error) {
	dailyJob := DailyJob{
		JobName: jobForm.JobName,
		JobType: jobForm.JobType,
		ImportantDegree: jobForm.ImportantDegree,
		OpenDeployAutoConfig: jobForm.OpenDeployAutoConfig,
		Remark: jobForm.Remark,
		TaskContent: jobForm.TaskContent,
		CreatorUserId: jobForm.CreatorUserId,
		CreatorUserName: jobForm.CreatorUserName,
		Status: 1,
		CreateTime: untils.GetNowTime(),
	}
	tx := database.MysqlClient.Begin()
	err = tx.Create(&dailyJob).Error
	if err != nil {
		tx.Rollback()
		untils.Log.Error(err.Error())
		return
	}
	tx.Commit()
	return
}

func UpdateDailyJob(jobForm forms.UpdateDailyJobForm) (err error) {
	dailyJob := DailyJob{
		ExecutorUserId: jobForm.UserId,
		ExecutorUserName: jobForm.UserName,
	}
	switch jobForm.Action {
	case "getJob":
		dailyJob.Status = int8(2)
		dailyJob.AcceptTime = untils.GetNowTime()
		total := 0
		database.MysqlClient.Model(&dailyJob).Where("status = 1 and id = ?", jobForm.Id).Count(&total)
		if total != 1 {
			return errors.New("该任务处于不能被领取的状态")
		}
		break
	case "finishJob":
		var jobInfo DailyJob
		dailyJob.Status = int8(3)
		dailyJob.EndTime = untils.GetNowTime()
		database.MysqlClient.Model(&dailyJob).Where("status = 2 and id = ?", jobForm.Id).Find(&jobInfo)
		if jobInfo.Id == 0 {
			return errors.New("该任务处于不能被设置为完成的状态")
		}
		dailyJob.AcceptTime = jobInfo.AcceptTime
		break
	case "deleteJob":
		dailyJob.Status = int8(-1)
		total := 0
		database.MysqlClient.Model(&dailyJob).Where("creator_user_id = ? and id = ?",
			jobForm.UserId, jobForm.Id).Count(&total)
		if total != 1 {
			return errors.New("用户只能删除自己创建的任务")
		}
		break
	default:
		return errors.New("不支持的操作")
	}

	err = database.MysqlClient.Exec("update daily_jobs set status = ?, executor_user_id = ?, " +
		"executor_user_name = ?, accept_time = ?, end_time = ? where id=?",
		dailyJob.Status, dailyJob.ExecutorUserId, dailyJob.ExecutorUserName, dailyJob.AcceptTime,
		dailyJob.EndTime, jobForm.Id).Error
	return
}

func UpdateDailyJobExecutorUser(data forms.UpdateDailyJobExecutorUserForm) (err error) {
 	infoList := strings.Split(data.ChangeUserId, "-")
 	var jobInfo DailyJob
 	err = database.MysqlClient.Raw("select * from daily_jobs where id= ? ",
 		data.JobId).Scan(&jobInfo).Error
 	if err != nil {
 		untils.Log.Error(err.Error())
 		return
	}
	updateSql := "update daily_jobs set executor_user_id = ?, executor_user_name = ? where id = ?"
	err = database.MysqlClient.Exec(updateSql, infoList[0], infoList[1], data.JobId).Error
	if err != nil {
		untils.Log.Error(err.Error())
	}
	return
}

func GetJobsCount(userId uint, queryKeyword string, queryCreateTime string) (total int) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	total = 0
	var querySql = "select count(*) from daily_jobs where status > 0 "
	var args []interface{}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and ( job_name like ? or creator_user_name like ? " +
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
	database.MysqlClient.Raw(querySql, args...).Count(&total)
	return
}

func GetJobsByPage(userId uint, offset int, limit int, queryKeyword string, queryCreateTime string) (jobs []DailyJob) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	var querySql = "select * from daily_jobs where status > 0 "
	var args []interface{}
	if queryKeyword != "" {
		querySql = fmt.Sprintf("%s and (job_name like ? or creator_user_name like ? " +
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
	database.MysqlClient.Raw(querySql, args...).Scan(&jobs)
	return
}

func GetJobDetail(id int) (job DailyJob) {
	database.MysqlClient.Find(&job, id)
	return
}
