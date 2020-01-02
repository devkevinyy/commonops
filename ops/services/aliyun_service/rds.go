package aliyun_service

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/untils"
	"go.uber.org/zap"
	"log"
	"math"
)

type DbInfo struct {
	DBInstance rds.DBInstance
	DBAttribute rds.DBInstanceAttribute
}

func RdsDiffCacheClean() (err error) {
	err = database.MysqlClient.Exec("delete from diff_caches where type = 'rds' ").Error
	return
}

func GetRdsList(accessKey string, keySecret string)(data []DbInfo){
	for _, region := range conf.RegionList {
		untils.Log.Info("开始获取阿里云rds数据", zap.String("accessKey", accessKey),
			zap.String("region", region))
		client, err := rds.NewClientWithAccessKey(region, accessKey, keySecret)
		if err != nil {
			untils.Log.Error(err.Error())
			continue
		}
		request := rds.CreateDescribeDBInstancesRequest()
		request.PageNumber = requests.NewInteger(1)
		request.PageSize = requests.NewInteger(50)
		resp, err := client.DescribeDBInstances(request)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		if resp.TotalRecordCount == 0 {
			untils.Log.Info("该账号在该地区无rds数据.")
			continue
		} else {
			untils.Log.Info("数据总数", zap.Int("msg", resp.TotalRecordCount))
			for _, instance := range resp.Items.DBInstance {
				dbInfo := DbInfo{
					DBInstance: instance,
					DBAttribute: GetRdsInstanceAttribute(client, instance.DBInstanceId),
				}
				data = append(data, dbInfo)
			}
			totalPage := int(math.Ceil(float64(resp.TotalRecordCount) / float64(resp.PageRecordCount)))
			untils.Log.Info("数据总页数", zap.Int("msg", totalPage))
			if totalPage > 1 {
				for i := 1; i<totalPage; i++ {
					currentPage := i + 1
					untils.Log.Info("开始获取数据页数", zap.Int("msg", currentPage))
					request.PageNumber = requests.NewInteger(currentPage)
					resp, err := client.DescribeDBInstances(request)
					if err != nil {
						fmt.Print(err.Error())
						break
					}
					for _, instance := range resp.Items.DBInstance {
						dbInfo := DbInfo{
							DBInstance: instance,
							DBAttribute: GetRdsInstanceAttribute(client, instance.DBInstanceId),
						}
						data = append(data, dbInfo)
					}
				}
			}
		}
	}
	return data
}

func GetRdsInstanceAttribute(client *rds.Client, instanceId string) (dbAttribute rds.DBInstanceAttribute) {
	request := rds.CreateDescribeDBInstanceAttributeRequest()
	request.DBInstanceId = instanceId
	response, err := client.DescribeDBInstanceAttribute(request)
	if err != nil {
		log.Println(err.Error())
		return
	}
	dbAttribute = response.Items.DBInstanceAttribute[0]
	return
}