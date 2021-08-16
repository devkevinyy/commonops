package models

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/jinzhu/gorm"
)

type DmsTable struct {
	gorm.Model
	DataStatus int8    `json:"DataStatus" gorm:"column:data_status;type:tinyint;not null;default:1"`
	TableId  string    `json:"TableId" gorm:"column:table_id;type:varchar(128);unique_index:table_id_idx"`
	DatabaseId  string    `json:"DatabaseId" gorm:"column:database_id;type:varchar(128)"`
	Tablename  string    `json:"TableName"  gorm:"column:table_name;type:varchar(32)"`
	TableSchemaName  string    `json:"TableSchemaName" gorm:"column:table_schema_name;type:varchar(128)"`
	Engine  string    `json:"Engine" gorm:"column:engine;type:varchar(128)"`
	Encoding   string    `json:"Encoding" gorm:"column:encoding;type:varchar(32)"`
	TableType  string    `json:"TableType" gorm:"column:table_type;type:varchar(256)"`
	NumRows  string    `json:"NumRows" gorm:"column:num_rows;type:varchar(32)"`
	StoreCapacity  string    `json:"StoreCapacity" gorm:"column:store_capacity;type:varchar(256)"`
}

func (DmsTable) TableName() string {
	return "dms_table"
}

func SaveDmsTable(databaseId string, tableId string, tablename string, tableSchemaName string,
	engine string, encoding string, tableType string, numRows int64, storeCapacity int64) (err error) {
	saveSql := fmt.Sprintf("INSERT INTO dms_table(database_id, table_id, table_name, table_schema_name," +
		"engine, encoding, table_type, num_rows, store_capacity) " +
		"VALUES(?,?,?,?,?,?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE database_id = ?, table_name = ?, table_schema_name = ?," +
		"engine = ?, encoding = ?, table_type = ?, num_rows = ?, store_capacity = ?")
	err = database.Mysql().Exec(saveSql, databaseId, tableId, tablename, tableSchemaName,
		engine, encoding, tableType, numRows, storeCapacity, databaseId, tablename, tableSchemaName,
		engine, encoding, tableType, numRows, storeCapacity).Error
	return
}

func GetDmsAllTableData(databaseId string) (data []DmsTable, err error) {
	err = database.Mysql().Raw("select table_id, table_name, table_schema_name, engine, encoding, " +
		"table_type, num_rows, store_capacity from dms_table where database_id = ?", databaseId).Scan(&data).Error
	return
}