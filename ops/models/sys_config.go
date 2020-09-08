package models

import "github.com/chujieyang/commonops/ops/database"

type SysConfig struct {
	Id    int8   `json:"Id" gorm:"primary key;auto increment;snot null;default: 1"`
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

func UpdatePrometheusValue(value string) (err error) {
	err = database.MysqlClient.Exec("update sys_config set value = ? where name = 'prometheus_host'", value).Error
	return
}

func GetPrometheusValue() (config SysConfig, err error) {
	err = database.MysqlClient.Raw("select * from sys_config where name = 'prometheus_host'").Scan(&config).Error
	return
}
