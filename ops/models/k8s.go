package models

import (
	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/forms"
	"github.com/chujieyang/commonops/ops/untils"
)

type K8s struct {
	Id int              `json:"id" gorm:"column:id"`
	ClusterId string `json:"clusterId" gorm:"column:cluster_id"`
	DataStatus int8     `json:"DataStatus" gorm:"column:data_status;not null;default:1"`
	Name string	`json:"name" gorm:"column:name"`
	Description string     	`json:"description" gorm:"column:description"`
	KubeConfigFilePath string  `json:"kubeconfigFilePath" gorm:"column:kubeconfig_file_path"`
}

func (K8s) TableName() string {
	return "k8s"
}

func AddK8sCluster(data forms.K8sForm) (err error) {
	k8s := K8s{
		Name: data.Name,
		Description: data.Description,
		KubeConfigFilePath: data.KubeConfigFilePath,
		ClusterId: untils.GetUUID(),
	}
	err = database.MysqlClient.Create(&k8s).Error
	return
}

func UpdateK8sCluster(data forms.K8sForm) (err error) {
	err = database.MysqlClient.Exec("update k8s set name=?, description=?, app_deploy_type=?, " +
		"kubeconfig_file_path = ? where id=?", data.Name, data.Description, data.KubeConfigFilePath,
		data.Id).Error
	return
}

func DeleteK8s(data forms.K8sForm) (err error) {
	err = database.MysqlClient.Exec("update k8s set data_status = 0 where id=?", data.Id).Error
	return
}

func GetK8sCount() (total int) {
	total = 0
	database.MysqlClient.Model(&K8s{}).Where("data_status > 0 ").Count(&total)
	return
}

func GetK8sList() (k8sData []K8s) {
	database.MysqlClient.Select("id, cluster_id, name, description").Where("data_status > 0").Order("id desc").Find(&k8sData)
	return
}

func GetKubeConfigPathByClusterId(clusterId string) string {
	var cluster K8s
	database.MysqlClient.Raw("select kubeconfig_file_path from k8s where cluster_id = ?", clusterId).Scan(&cluster)
	return cluster.KubeConfigFilePath
}
