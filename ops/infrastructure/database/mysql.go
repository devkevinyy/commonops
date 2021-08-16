package database

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/opslog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var mysqlClient *gorm.DB
var err error

func init() {
	opslog.Info().Printf("mysql host: %s \n", conf.MysqlHost)
	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.MysqlUser,
		conf.MysqlPwd, conf.MysqlHost, conf.MysqlDb)
	if mysqlClient, err = gorm.Open("mysql", dbConnectionString); err != nil {
		opslog.Error().Fatalf("mysql init exception: %s \n", err.Error())
	}
	mysqlClient.DB().SetMaxIdleConns(10)
	mysqlClient.DB().SetMaxOpenConns(100)
}

func Mysql() *gorm.DB {
	return mysqlClient
}