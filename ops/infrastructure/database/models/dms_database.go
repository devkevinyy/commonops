package models

import (
	"fmt"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type DmsDatabase struct {
	gorm.Model
	DataStatus int8   `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	InstanceId string `json:"InstanceId" gorm:"column:instance_id;type:varchar(128)"`
	DatabaseId string `json:"DatabaseId" gorm:"column:database_id;type:varchar(128);unique_index:database_id_idx"`
	SchemaName string `json:"SchemaName" gorm:"column:schema_name;type:varchar(128)"`
	State      string `json:"State" gorm:"column:state;type:varchar(32)"`
	DbType     string `json:"DbType" gorm:"column:db_type;type:varchar(32)"`
	Host       string `json:"Host" gorm:"column:host;type:varchar(256)"`
	Port       int    `json:"Port" gorm:"column:port;type:smallint"`
}

func (DmsDatabase) TableName() string {
	return "dms_database"
}

func SaveDmsDatabase(instanceId string, databaseId string, schemaName string, state string, dbType string, host string,
	port int) (err error) {
	saveSql := fmt.Sprintf("INSERT INTO dms_database(instance_id, database_Id, schema_name, state, db_type, host, port) " +
		"VALUES(?,?,?,?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE instance_id = ?, schema_name = ?, state = ?, db_type = ?, host = ?, " +
		"port = ?")
	err = database.Mysql().Exec(saveSql, instanceId, databaseId, schemaName, state, dbType, host, port, instanceId,
		schemaName, state, dbType, host, port).Error
	return
}

func GetDmsAllDatabaseData(instanceId string) (data []DmsDatabase, err error) {
	err = database.Mysql().Raw("select id, instance_id, database_id, schema_name, state, db_type, host, port "+
		"from dms_database where data_status = 1 and instance_id = ?", instanceId).Scan(&data).Error
	return
}

func GetDmsDatabaseInfoByDatabaseId(databaseId string) (databaseI DmsDatabase, err error) {
	err = database.Mysql().Raw("select instance_id, database_id, schema_name, db_type, host, port from dms_database where database_id = ?", databaseId).Scan(&databaseI).Error
	return
}

func DeleteDmsDatabase(instanceId string, databaseId string) (err error) {
	err = database.Mysql().Exec("update dms_database set data_status = 0 where instance_id = ? and database_id = ?", instanceId, databaseId).Error
	return
}