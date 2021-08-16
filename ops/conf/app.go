package conf

import (
	"log"

	"github.com/chujieyang/commonops/ops/opslog"
	"gopkg.in/ini.v1"
)

var Host string
var Port string
var SecretSalt string
var DesKey string
var DmsConcurrentOperationNum int

var MysqlHost string
var MysqlUser string
var MysqlPwd string
var MysqlDb string

var JenkinsHost string
var JenkinsAuthUser string
var JenkinsAuthToken string

func init() {
	opslog.Info().Println("config.ini 配置初始化...")
	cfg, err := ini.Load("./conf/config.ini")
	if err != nil {
		log.Fatalf("Fail to read file: %v \n", err)
	}

	Host = cfg.Section("app").Key("Host").String()
	Port = cfg.Section("app").Key("Port").String()
	SecretSalt = cfg.Section("app").Key("SecretSalt").String()
	DesKey = cfg.Section("app").Key("DesKey").String()
	DmsConcurrentOperationNum = cfg.Section("app").Key("DmsConcurrentOperationNum").MustInt(2)

	MysqlHost = cfg.Section("database").Key("MysqlHost").String()
	MysqlUser = cfg.Section("database").Key("MysqlUser").String()
	MysqlPwd = cfg.Section("database").Key("MysqlPwd").String()
	MysqlDb = cfg.Section("database").Key("MysqlDb").String()

	JenkinsHost = cfg.Section("jenkins").Key("JenkinsHost").String()
	JenkinsAuthUser = cfg.Section("jenkins").Key("JenkinsAuthUser").String()
	JenkinsAuthToken = cfg.Section("jenkins").Key("JenkinsAuthToken").String()
}
