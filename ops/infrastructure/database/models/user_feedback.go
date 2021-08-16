package models

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"time"
)

type UserFeedback struct {
	Id int              `json:"id" gorm:"column:id"`
	DataStatus int8     `json:"DataStatus" gorm:"column:data_status;not null;default:1"`
	UserId int	`json:"userId" gorm:"column:user_id"`
	Username string	`json:"username" gorm:"column:username"`
	Content string	`json:"content" gorm:"column:content"`
	Score int `json:"score" gorm:"column:score;not null;default:0"`
	CreateTime string `json:"createTime" gorm:"create_time"`
}

func (UserFeedback) TableName() string {
	return "user_feedback"
}

func SaveUserFeedback(userId uint, username string, content string, score int) (err error) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = database.Mysql().Exec("insert into user_feedback(user_id, username, content, score, create_time) " +
		"values (?, ?, ?, ?, ?)", userId, username, content, score, nowTime).Error
	return
}

func GetUserFeedbackCount() (total int) {
	total = 0
	database.Mysql().Model(&UserFeedback{}).Where("data_status > 0").Count(&total)
	return
}

func GetUserFeedbackList(offset int, pageSize int) (userFeedback []UserFeedback) {
	database.Mysql().Where("data_status > 0").Order("id desc").Offset(offset).Limit(pageSize).Find(&userFeedback)
	return
}