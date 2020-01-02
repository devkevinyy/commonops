package database

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/untils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
	"log"
)

var MysqlClient *gorm.DB
var err error


func init(){
	log.Println("mysql host: ", conf.MysqlHost)
	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.MysqlUser, conf.MysqlPwd,
		conf.MysqlHost, conf.MysqlDb)
	MysqlClient, err = gorm.Open("mysql", dbConnectionString)
	if err != nil {
		untils.Log.Error("[mysql init]", zap.String("msg", err.Error()))
	}
	MysqlClient.DB().SetMaxIdleConns(10)
	MysqlClient.DB().SetMaxOpenConns(100)
	//defer MysqlClient.Close()
}