package service

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/value_objects/daily_job"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/16/21 3:07 PM
 * @Desc:
 */

type dailyJob struct {}

var dailyJobInstance = &dailyJob{}

func GetDailyJobService() *dailyJob {
	return dailyJobInstance
}

func (d *dailyJob) AddDailyJob(req daily_job.AddDailyJobForm) error {
	return models.AddDailyJob(req)
}

func (d *dailyJob) UpdateDailyJob(req daily_job.UpdateDailyJobForm) error {
	return models.UpdateDailyJob(req)
}

func (d *dailyJob) UpdateDailyJobExecutorUser(req daily_job.UpdateDailyJobExecutorUserForm) error {
	return models.UpdateDailyJobExecutorUser(req)
}

func (d *dailyJob) GetDailyJobDataByPage(userId uint, queryKeyword string, queryCreateTime string, page uint, size uint) (uint, []models.DailyJob) {
	offset := (page - 1) * size
	total := models.GetJobsCount(userId, queryKeyword, queryCreateTime)
	slbList := models.GetJobsByPage(userId, queryKeyword, queryCreateTime, offset, size)
	return total, slbList
}

func (d *dailyJob) GetDailyJobDetail(id int) models.DailyJob {
	return models.GetJobDetail(id)
}
