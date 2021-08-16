package service

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	aliRds "github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/opslog"
	rds2 "github.com/chujieyang/commonops/ops/value_objects/rds"
	"math"
	"time"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:02 AM
 * @Desc:
 */

type rds struct {}

var rdsInstance = &rds{}

func GetRdsService() *rds {
	return rdsInstance
}

func (r *rds) GetRdsDetail(id uint) models.RdsDetail {
	return models.GetRdsDetail(id)
}

func (r *rds) GetRdsDataByPage(userId uint, expiredTime string, queryKeyword string, cloudAccountId uint,
	page uint, size uint) (total uint, dataList []models.RdsDetail) {
	offset := (page - 1) * size
	total = models.GetRdsCount(userId, expiredTime, queryKeyword, uint(cloudAccountId))
	dataList = models.GetRdsByPage(userId, expiredTime, queryKeyword, uint(cloudAccountId), offset, size)
	return
}

func (r *rds) UpdateCloudRds(params rds2.ExtraInfoForm) error {
	return  models.UpdateCloudRds(params)
}

func (r *rds) AddCloudRds(params rds2.RdsInfoForm) error {
	return  models.AddCloudRds(params)
}

func (r *rds) DeleteCloudRds(id uint) error {
	return models.DeleteCloudRds(id)
}

/*
	获取阿里云 rds 的各项监控数据
*/
func (r *rds) GetRdsMonitor(instanceId string, timeDimension string, metricDimension string,
	region string, accessKey string, accessSecret string) (response *cms.DescribeMetricDataResponse) {
	client, err := cms.NewClientWithAccessKey(region, accessKey, accessSecret)
	if err != nil {
		opslog.Info().Println("aliyun sdk client init 失败: ", err.Error())
		return
	}
	request := cms.CreateDescribeMetricDataRequest()

	period := "60"
	now := time.Now().Unix()
	start := now - 3600
	switch timeDimension {
	case "1h":
		start = now - 3600
	case "6h":
		start = now - (3600 * 6)
	case "12h":
		start = now - (3600 * 12)
	case "1d":
		start = now - (3600 * 24)
	case "3d":
		start = now - (3600 * 3 * 24)
	case "7d":
		start = now - (3600 * 7 * 24)
	case "14d":
		start = now - (3600 * 14 * 24)
	}

	request.Dimensions = fmt.Sprintf("{'instanceId': '%s'}", instanceId)
	request.StartTime = time.Unix(start, 0).Format("2006-01-02 15:04:05")
	request.EndTime = time.Unix(now, 0).Format("2006-01-02 15:04:05")
	request.Namespace = "acs_rds_dashboard"
	request.Period = period
	request.MetricName = metricDimension

	response, err = client.DescribeMetricData(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	return response
}

func (r *rds) RdsDiffCacheClean() error {
	return models.CleanRdsDiffCaches()
}

func (r *rds) GetRdsList(accessKey string, keySecret string) (data []rds2.DbInfo) {
	for _, region := range conf.RegionList {
		client, err := aliRds.NewClientWithAccessKey(region, accessKey, keySecret)
		if err != nil {
			opslog.Error().Println(err.Error())
			continue
		}
		request := aliRds.CreateDescribeDBInstancesRequest()
		request.PageNumber = requests.NewInteger(1)
		request.PageSize = requests.NewInteger(50)
		resp, err := client.DescribeDBInstances(request)
		if err != nil {
			opslog.Error().Println(err)
			continue
		}
		if resp.TotalRecordCount == 0 {
			continue
		} else {
			for _, instance := range resp.Items.DBInstance {
				dbInfo := rds2.DbInfo{
					DBInstance:  instance,
					DBAttribute: r.GetRdsInstanceAttribute(client, instance.DBInstanceId),
				}
				data = append(data, dbInfo)
			}
			totalPage := int(math.Ceil(float64(resp.TotalRecordCount) / float64(resp.PageRecordCount)))
			if totalPage > 1 {
				for i := 1; i < totalPage; i++ {
					currentPage := i + 1
					request.PageNumber = requests.NewInteger(currentPage)
					resp, err := client.DescribeDBInstances(request)
					if err != nil {
						opslog.Error().Println(err)
						break
					}
					for _, instance := range resp.Items.DBInstance {
						dbInfo := rds2.DbInfo{
							DBInstance:  instance,
							DBAttribute: r.GetRdsInstanceAttribute(client, instance.DBInstanceId),
						}
						data = append(data, dbInfo)
					}
				}
			}
		}
	}
	return data
}

func (r *rds) GetRdsInstanceAttribute(client *aliRds.Client, instanceId string) (dbAttribute aliRds.DBInstanceAttribute) {
	request := aliRds.CreateDescribeDBInstanceAttributeRequest()
	request.DBInstanceId = instanceId
	response, err := client.DescribeDBInstanceAttribute(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	dbAttribute = response.Items.DBInstanceAttribute[0]
	return
}