package models

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/chujieyang/commonops/ops/utils"
	k8s2 "github.com/chujieyang/commonops/ops/value_objects/k8s"
)

type K8s struct {
	Id          int    `json:"id" gorm:"column:id"`
	ClusterId   string `json:"clusterId" gorm:"column:cluster_id"`
	DataStatus  int8   `json:"DataStatus" gorm:"column:data_status;not null;default:1"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Token       string `json:"token" gorm:"column:token"`
	ApiServer   string `json:"apiServer" gorm:"column:api_server"`
}

func (K8s) TableName() string {
	return "k8s"
}

func AddK8sCluster(data k8s2.K8sForm) (err error) {
	k8s := K8s{
		Name:        data.Name,
		Description: data.Description,
		Token:       data.Token,
		ClusterId:   utils.GetUUID(),
		ApiServer:   data.ApiServer,
	}
	err = database.Mysql().Create(&k8s).Error
	return
}

func DeleteK8s(data k8s2.K8sForm) (err error) {
	err = database.Mysql().Exec("update k8s set data_status = 0 where id=?", data.Id).Error
	return
}

func GetK8sCount() (total int) {
	total = 0
	database.Mysql().Model(&K8s{}).Where("data_status > 0 ").Count(&total)
	return
}

func GetK8sList() (k8sData []K8s) {
	database.Mysql().Select("id, cluster_id, name, description").Where("data_status > 0").Order("id desc").Find(&k8sData)
	return
}

func GetK8sInfoByClusterId(clusterId string) (apiServer string, token string) {
	var cluster K8s
	database.Mysql().Raw("select * from k8s where cluster_id = ?", clusterId).Scan(&cluster)
	return cluster.ApiServer, cluster.Token
}
