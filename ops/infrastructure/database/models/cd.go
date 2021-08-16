package models

import (
	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type CdProcessTemplate struct {
	gorm.Model
	DataStatus    int8   `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	EmpId         string `json:"EmpId" gorm:"column:emp_id;type:varchar(32)"`
	JobName       string `json:"JobName" gorm:"column:job_name;type:varchar(256)"`
	TemplateName  string `json:"TemplateName" gorm:"column:template_name;type:varchar(512)"`
	ClusterId     string `json:"ClusterId" gorm:"column:cluster_id;type:varchar(128)"`
	Namespace     string `json:"Namespace" gorm:"column:namespace;type:varchar(32)"`
	DeployYaml    string `json:"DeployYaml" gorm:"column:deploy_yaml;type:text"`
	ServiceYaml   string `json:"ServiceYaml" gorm:"column:service_yaml;type:text"`
	ConfigmapYaml string `json:"ConfigmapYaml" gorm:"column:configmap_yaml;type:text"`
	IngressYaml   string `json:"IngressYaml" gorm:"column:ingress_yaml;type:text"`
}

func (CdProcessTemplate) TableName() string {
	return "cd_process_template"
}

func IsTemplateNameExist(templateName string) (count int, err error) {
	err = database.Mysql().Raw("select count(*) from cd_process_template where data_status = 1 and template_name = ?", templateName).Count(&count).Error
	return
}

func AddNewCdProcessTemplate(empId string, jobName string, templateName string, clusterId string, namespace string, deployYAML string,
	serviceYAML string, configmapYAML string, ingressYAML string) (id int, err error) {
	cdProcessTemplate := &CdProcessTemplate{
		EmpId:         empId,
		JobName:       jobName,
		TemplateName:  templateName,
		ClusterId:     clusterId,
		Namespace:     namespace,
		DeployYaml:    deployYAML,
		ServiceYaml:   serviceYAML,
		ConfigmapYaml: configmapYAML,
		IngressYaml:   ingressYAML,
	}
	err = database.Mysql().Create(cdProcessTemplate).Error
	id = int(cdProcessTemplate.ID)
	return
}

func GetTemplateInfoById(templateId int) (info CdProcessTemplate, err error) {
	err = database.Mysql().Raw("select * from cd_process_template where id = ?", templateId).Scan(&info).Error
	return
}

func GetAllCdProcessTemplateName(empId string) (cdProcessTemplate []CdProcessTemplate, err error) {
	err = database.Mysql().Raw("select id, job_name, template_name from cd_process_template where data_status = 1 and emp_id = ?", empId).Scan(&cdProcessTemplate).Error
	return
}

type CdProcessLog struct {
	gorm.Model
	DataStatus int8   `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	EmpId      string `json:"EmpId" gorm:"column:emp_id;type:varchar(32)"`
	TemplateId int    `json:"TemplateId" gorm:"column:template_id;type:int"`
	ImageName  string `json:"ImageName" gorm:"column:image_name;type:varchar(256)"`
	Success    int    `json:"Success" gorm:"column:success;type:int"`
	Result     string `json:"Result" gorm:"column:result;type:varchar(1024)"`
}

func (CdProcessLog) TableName() string {
	return "cd_process_log"
}

type CdProcessLogDetail struct {
	Id           int    `json:"id"`
	TemplateId   int    `json:"templateId"`
	ImageName    string    `json:"imageName"`
	Success      int    `json:"success"`
	Result       string `json:"result"`
	TemplateName string `json:"templateName"`
	Namespace    string `json:"namespace"`
	JobName      string `json:"jobName"`
}

func AddCdProcessLog(empId string, templateId int, imageName string) (id int, err error) {
	cdProcessLog := &CdProcessLog{
		EmpId:      empId,
		TemplateId: templateId,
		ImageName:  imageName,
		Success:    0,
		Result:     "",
	}
	err = database.Mysql().Create(cdProcessLog).Error
	id = int(cdProcessLog.ID)
	return
}

func UpdateCdProcessLogResult(logId int, isSuccess int, result string) (err error) {
	querySql := "update cd_process_log set success = ?, result = ? where id = ?"
	err = database.Mysql().Exec(querySql, isSuccess, result, logId).Error
	return
}

func GetCdProcessLogCount(empId string) (count int, err error) {
	querySql := "select count(*) from cd_process_log where data_status = 1 and emp_id = ?"
	err = database.Mysql().Raw(querySql, empId).Count(&count).Error
	return
}

func GetCdProcessLogByPage(empId string, offset int, limit int) (log []CdProcessLogDetail, err error) {
	querySql := "select l.id, l.template_id, l.image_name, l.success, l.result, t.template_name, " +
		" t.cluster_id, t.namespace, t.job_name from cd_process_log as l inner join cd_process_template as t on l.template_id = t.id " +
		" where l.data_status = 1 and l.emp_id = ? order by id desc limit ?, ?"
	err = database.Mysql().Raw(querySql, empId, offset, limit).Scan(&log).Error
	return
}
