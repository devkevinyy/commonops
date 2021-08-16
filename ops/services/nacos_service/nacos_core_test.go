package nacos_service

import (
	"encoding/json"
	"log"
	"testing"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 7/20/21 11:12 AM
 * @Desc:
 */

func Test_GetConfig(t *testing.T) {
	nacos, err := NewNacosClient("127.0.0.1:8849", "nacos", "nacos")
	if err != nil {
		t.Fatal(err)
	}
	// nsList, err := nacos.GetNamespace()
	// fmt.Println(err, nsList)

	// data, err := nacos.GetNsConfigs("dev", 1, 10)
	// fmt.Println(err, data)

	// data, err := nacos.GetConfig("dev", "com.ops", "svc")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(data)
	//err = nacos.CopyConfig("dev", "com.fulu.ops", "mysql", "prod", "com.fulu.ops1", "mysql")
	//if err != nil {
	//	t.Fatal(err)
	//}

	allConfigs, _ := nacos.GetAllConfigs()
	configBytes, _ := json.Marshal(allConfigs)
	log.Println(string(configBytes))

	// destConfigList := []map[string]string{
	// 	{"namespace": "dev", "dataId": "com.ops", "group": "svc"},
	// 	{"namespace": "prod", "dataId": "bbbbbbbbbbbb", "group": "bb01"},
	// 	{"namespace": "prod", "dataId": "qqqqqqq", "group": "qq01"},
	// 	{"namespace": "prod", "dataId": "wwwww", "group": "www01"},
	// }
	// err = nacos.AppendStaticConfigToSelectAllConfigs("dev", "aaaa", "aa01a", destConfigList)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
