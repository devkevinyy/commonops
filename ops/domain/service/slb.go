package service

import (
	"fmt"
	"math"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	aliSlb "github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/opslog"
	"go.uber.org/zap"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:40 AM
 * @Desc:
 */

type slb struct{}

var slbInstance = &slb{}

func GetSlbService() *slb {
	return slbInstance
}

func (s *slb) GetSlbDetail(id uint) models.SlbDetail {
	return models.GetSlbDetail(id)
}

func (s *slb) GetSlbDataByPage(userId uint, queryKeyword string, cloudAccountId uint, page uint, size uint) (total uint, data []models.SlbDetail) {
	offset := (page - 1) * size
	total = models.GetSlbCount(userId, queryKeyword, cloudAccountId)
	data = models.GetSlbByPage(userId, queryKeyword, cloudAccountId, offset, size)
	return
}

func (s *slb) DeleteSlb(id uint) error {
	return models.DeleteCloudSlb(id)
}

/*
	获取阿里云 slb 的各项监控数据
*/
func (s *slb) GetSlbMonitor(instanceId string, timeDimension string, metricDimension string,
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
	request.Namespace = "acs_slb_dashboard"
	request.Period = period
	request.MetricName = metricDimension

	response, err = client.DescribeMetricData(request)
	if err != nil {
		opslog.Error().Println(err)
	}
	return response
}

func (s *slb) SlbDiffCacheClean() error {
	return models.CleanSlbDiffCaches()
}

func (s *slb) GetSlbList(accessKey string, keySecret string) (data []aliSlb.LoadBalancer) {
	for _, region := range conf.RegionList {
		opslog.Info().Println("开始获取阿里云slb数据", zap.String("accessKey", accessKey),
			zap.String("region", region))
		client, err := aliSlb.NewClientWithAccessKey(region, accessKey, keySecret)
		if err != nil {
			opslog.Error().Println(err)
			continue
		}
		request := aliSlb.CreateDescribeLoadBalancersRequest()
		request.PageNumber = requests.NewInteger(1)
		request.PageSize = requests.NewInteger(50)
		resp, err := client.DescribeLoadBalancers(request)
		if err != nil {
			opslog.Error().Println(err)
			continue
		}
		if resp.TotalCount == 0 {
			opslog.Info().Println("该账号在该地区无slb数据.")
			continue
		} else {
			opslog.Info().Println("数据总数", zap.Int("msg", resp.TotalCount))
			data = append(data, []aliSlb.LoadBalancer(resp.LoadBalancers.LoadBalancer)...)
			totalPage := int(math.Ceil(float64(resp.TotalCount) / float64(resp.PageSize)))
			opslog.Info().Println("数据总页数", zap.Int("msg", totalPage))
			if totalPage > 1 {
				for i := 1; i < totalPage; i++ {
					currentPage := i + 1
					opslog.Info().Println("开始获取数据页数", zap.Int("msg", currentPage))
					request.PageNumber = requests.NewInteger(currentPage)
					resp, err := client.DescribeLoadBalancers(request)
					if err != nil {
						opslog.Error().Println(err.Error())
						break
					}
					data = append(data, []aliSlb.LoadBalancer(resp.LoadBalancers.LoadBalancer)...)
				}
			}
		}
	}
	return data
}
