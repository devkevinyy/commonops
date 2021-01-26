package models

type SysConfig struct {
	Id    int8   `json:"Id" gorm:"primary key;auto increment;snot null;default: 1"`
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}
