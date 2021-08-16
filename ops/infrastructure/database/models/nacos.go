package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type Nacos struct {
	Id         int64  `json:"Id" gorm:"column:id;type:int;PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	DataStatus int8   `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	EndPoint   string `json:"EndPoint" gorm:"column:end_point;type:varchar(255)"`
	Alias      string `json:"Alias" gorm:"column:alias;type:varchar(255)"`
	Username   string `json:"Username" gorm:"column:username;type:varchar(255)"`
	Password   string `json:"Password" gorm:"column:password;type:varchar(255)"`
}

func (Nacos) TableName() string {
	return "nacos"
}

func AddNewNacosServer(alias, endpoint, username, password string) (id int, err error) {
	count := 0
	if err = database.Mysql().Raw("select count(*) from nacos where end_point = ? and data_status = 1", endpoint).Count(&count).Error; err != nil {
		return
	}
	if count > 0 {
		err = errors.New("已存在相同EndPoint的集群")
		return
	}
	nacos := &Nacos{
		Alias:    alias,
		EndPoint: endpoint,
		Username: username,
		Password: password,
	}
	err = database.Mysql().Create(nacos).Error
	id = int(nacos.Id)
	return
}

func GetNacosInfoById(id string) (info Nacos, err error) {
	err = database.Mysql().Raw("select * from nacos where id = ? limit 1", id).Scan(&info).Error
	return
}

func GetNacosList() (data []Nacos, err error) {
	querySql := "select id, end_point, alias from nacos where data_status = 1"
	err = database.Mysql().Raw(querySql).Scan(&data).Error
	return
}

type ConfigTemplate struct {
	Id            int64  `json:"Id" gorm:"column:id;type:int;PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	DataStatus    int8   `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	Name          string `json:"Name" gorm:"column:name;type:varchar(512)"`
	ConfigContent string `json:"ConfigContent" gorm:"column:config_content;type:text"`
	FillField     string `json:"FillField" gorm:"column:fill_field;type:text"`
	UpdateTime    string `json:"UpdateTime" gorm:"column:update_time;type:timestamp"`
}

func (ConfigTemplate) TableName() string {
	return "config_template"
}

func AddConfigTemplate(name, configContent string, fillField []string) error {
	fillFieldBytes, err := json.Marshal(fillField)
	if err != nil {
		return err
	}
	config := &ConfigTemplate{
		Name:          name,
		ConfigContent: configContent,
		FillField:     string(fillFieldBytes),
		UpdateTime:    time.Now().Format("2006-01-02 15:04:05"),
	}
	return database.Mysql().Create(config).Error
}

func GetConfigTemplateCount(name string) (count int, err error) {
	querySql := "select count(*) from config_template where data_status = 1 and name like ?"
	err = database.Mysql().Raw(querySql, "%"+name+"%").Count(&count).Error
	return
}

func GetConfigTemplateByPage(name string, offset, size int) (data []ConfigTemplate, err error) {
	querySql := "select * from config_template where data_status = 1 and name like ? order by id desc limit ?, ?"
	err = database.Mysql().Raw(querySql, "%"+name+"%", offset, size).Scan(&data).Error
	return
}

func GetAllNacosConfigRelatedToTheTemplate(tx *gorm.DB, id int) (configList []ConfigDetail, err error) {
	querySql := "select c.*, t.config_content, t.fill_field from nacos_config as c left join config_template as t on c.template_id = t.id where t.id = ? "
	if err = tx.Raw(querySql, id).Scan(&configList).Error; err != nil {
		return
	}
	return
}

func UpdateConfigTemplate(tx *gorm.DB, id int, name, configContent string, fillField string) (err error) {
	updateSql := "update config_template set name = ?, config_content = ?, fill_field = ? where id = ?"
	err = tx.Exec(updateSql, name, configContent, fillField, id).Error
	return
}

func GetConfigTemplateById(id int) (data ConfigTemplate, err error) {
	querySql := "select * from config_template where id = ? "
	err = database.Mysql().Raw(querySql, id).Scan(&data).Error
	return
}

func DeleteConfigTemplate(id int) (err error) {
	updateSql := "update config_template set data_status = 0 where id = ?"
	err = database.Mysql().Exec(updateSql, id).Error
	return
}

func GetConfigTemplatesAll() (data []ConfigTemplate, err error) {
	querySql := "select id, name, config_content, fill_field from config_template where data_status = 1 "
	err = database.Mysql().Raw(querySql).Scan(&data).Error
	return
}

func GetConfigTemplateDetail(id int) (data ConfigTemplate, err error) {
	querySql := "select * from config_template where id = ? "
	err = database.Mysql().Raw(querySql, id).Scan(&data).Error
	return
}

type NacosConfig struct {
	Id         int64  `json:"Id" gorm:"column:id;type:int;PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	DataStatus int8   `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	ClusterId  string `json:"ClusterId" gorm:"column:cluster_id;type:int"`
	Namespace  string `json:"Namespace" gorm:"column:namespace;type:varchar(255)"`
	DataId     string `json:"DataId" gorm:"column:data_id;type:varchar(255)"`
	Group      string `json:"Group" gorm:"column:config_group;type:varchar(255)"`
	TemplateId int    `json:"TemplateId" gorm:"column:template_id;type:int"`
	FillData   string `json:"FillData" gorm:"column:fill_data;type:varchar(1024)"`
}

func (NacosConfig) TableName() string {
	return "nacos_config"
}

func AddNacosConfig(tx *gorm.DB, clusterId string, templateId int, namspace, dataId, group, fillData string) (err error) {
	config := &NacosConfig{
		ClusterId:  clusterId,
		Namespace:  namspace,
		DataId:     dataId,
		Group:      group,
		TemplateId: templateId,
		FillData:   fillData,
	}
	return tx.Create(config).Error
}

func UpdateNacosConfig(tx *gorm.DB, id int, templateId int, dataId, group, fillData string) (err error) {
	updateSql := "update nacos_config set data_id = ?, config_group = ?, template_id = ?, fill_data = ? where id = ?"
	err = tx.Exec(updateSql, dataId, group, templateId, fillData, id).Error
	return
}

func DeleteNacosConfig(tx *gorm.DB, id int) (err error) {
	updateSql := "update nacos_config set data_status = 0 where id = ?"
	err = tx.Exec(updateSql, id).Error
	return
}

type ConfigDetail struct {
	Id            int64  `json:"Id"`
	DataStatus    int8   `json:"DataStatus"`
	ClusterId     string `json:"ClusterId"`
	Namespace     string `json:"Namespace"`
	DataId        string `json:"DataId"`
	ConfigGroup         string `json:"ConfigGroup"`
	TemplateId    int    `json:"TemplateId"`
	FillData      string `json:"FillData"`
	ConfigContent string `json:"ConfigContent"`
	FillField     string `json:"FillField"`
}

func GetNacosConfig(clusterId string, namspace, dataId, group string) (data ConfigDetail, err error) {
	querySql := "select c.*, t.config_content, t.fill_field from nacos_config as c left join config_template as t on c.template_id = t.id where c.cluster_id = ? and c.namespace = ? and c.data_id = ? and c.config_group = ? limit 1 "
	err = database.Mysql().Raw(querySql, clusterId, namspace, dataId, group).Scan(&data).Error
	return
}

func GetOpsNacosConfigCount(clusterId, namespace string) (count int, err error) {
	querySql := "select count(*) from nacos_config where data_status = 1 and cluster_id = ? and namespace = ?"
	err = database.Mysql().Raw(querySql, clusterId, namespace).Count(&count).Error
	return
}

type OpsNacosConfigContent struct {
	Id            int64  `json:"Id"`
	DataStatus    int8   `json:"DataStatus"`
	ClusterId     string `json:"ClusterId"`
	Namespace     string `json:"Namespace"`
	DataId        string `json:"DataId"`
	ConfigGroup   string `json:"ConfigGroup"`
	TemplateId    int    `json:"TemplateId"`
	FillData      string `json:"FillData"`
	Name          string `json:"Name"`
	ConfigContent string `json:"ConfigContent"`
	FillField     string `json:"FillField"`
}

func GetOpsNacosConfigByPage(clusterId, namespace string, offset, size int) (data []OpsNacosConfigContent, err error) {
	querySql := "select c.*, t.name, t.config_content, t.fill_field from nacos_config as c left join config_template as t on c.template_id = t.id where c.data_status = 1 and c.cluster_id = ? and c.namespace = ? order by c.id desc limit ?, ?"
	err = database.Mysql().Raw(querySql, clusterId, namespace, offset, size).Scan(&data).Error
	return
}
