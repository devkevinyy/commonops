package service

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/opslog"
	kv2 "github.com/chujieyang/commonops/ops/value_objects/kv"
	"math"
	"time"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:25 AM
 * @Desc:
 */

type kv struct {}

var kvInstance = &kv{}

func GetKvService() *kv {
	return kvInstance
}

func (k *kv) GetKvDetail(id uint) models.KvDetail {
	return models.GetKvDetail(id)
}

func (k *kv) GetKvDataByPage(userId uint, expiredTime string, queryKeyword string, cloudAccountId uint, page uint, size uint) (total uint, kvList []models.KvDetail) {
	offset := (page - 1) * size
	total = models.GetKvCount(userId, expiredTime, queryKeyword, cloudAccountId)
	kvList = models.GetKvByPage(userId, expiredTime, queryKeyword, cloudAccountId, offset, size)
	return
}

func (k *kv) UpdateKv(params kv2.ExtraInfoForm) error {
	return models.UpdateCloudKv(params)
}

func (k *kv) AddKv(params kv2.KvInfoForm) error {
	return models.AddCloudKv(params)
}

func (k *kv) DeleteKv(id uint) error {
	return models.DeleteCloudKv(id)
}

func (k *kv) KvDiffCacheClean() error {
	return models.CleanKvDiffCaches()
}

func (k *kv) GetKvList(accessKey string, keySecret string) (data []r_kvstore.KVStoreInstance) {
	for _, region := range conf.RegionList {
		client, err := r_kvstore.NewClientWithAccessKey(region, accessKey, keySecret)
		if err != nil {
			opslog.Error().Println(err)
			continue
		}
		request := r_kvstore.CreateDescribeInstancesRequest()
		request.PageNumber = requests.NewInteger(1)
		request.PageSize = requests.NewInteger(50)
		resp, err := client.DescribeInstances(request)
		if err != nil {
			opslog.Error().Println(err)
			continue
		}
		if resp.TotalCount == 0 {
			continue
		} else {
			data = append(data, []r_kvstore.KVStoreInstance(resp.Instances.KVStoreInstance)...)
			totalPage := int(math.Ceil(float64(resp.TotalCount) / float64(resp.PageSize)))
			if totalPage > 1 {
				for i := 1; i < totalPage; i++ {
					currentPage := i + 1
					request.PageNumber = requests.NewInteger(currentPage)
					resp, err := client.DescribeInstances(request)
					if err != nil {
						opslog.Error().Println(err)
						break
					}
					data = append(data, []r_kvstore.KVStoreInstance(resp.Instances.KVStoreInstance)...)
				}
			}
		}
	}
	return data
}

/*
	获取阿里云 kv-store 的各项监控数据
*/
func (k *kv) GetKvMonitor(instanceId string, timeDimension string, metricDimension string,
	region string, accessKey string, accessSecret string) (response *cms.DescribeMetricDataResponse) {
	client, err := cms.NewClientWithAccessKey(region, accessKey, accessSecret)
	if err != nil {
		opslog.Error().Printf("aliyun sdk client init 失败: %s \n", err.Error())
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
	request.Namespace = "acs_kvstore"
	request.Period = period
	request.MetricName = metricDimension
	response, err = client.DescribeMetricData(request)
	if err != nil {
		opslog.Error().Println(err)
	}
	return response
}

