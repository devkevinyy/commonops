package aliyun_service

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"time"
)


/*
	获取阿里云 ecs 的各项监控数据
 */
func GetEcsMonitor(instanceId string, timeDimension string, metricDimension string,
					region string, accessKey string, accessSecret string) (response *cms.DescribeMetricDataResponse) {
	client, err := cms.NewClientWithAccessKey(region, accessKey, accessSecret)
	if err != nil {
		fmt.Println("aliyun sdk client init 失败: ", err.Error())
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
	request.Namespace = "acs_ecs_dashboard"
	request.Period = period
	request.MetricName = metricDimension

	response, err = client.DescribeMetricData(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	return response
}

/*
	获取阿里云 rds 的各项监控数据
*/
func GetRdsMonitor(instanceId string, timeDimension string, metricDimension string,
	region string, accessKey string, accessSecret string) (response *cms.DescribeMetricDataResponse) {
	client, err := cms.NewClientWithAccessKey(region, accessKey, accessSecret)
	if err != nil {
		fmt.Println("aliyun sdk client init 失败: ", err.Error())
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

/*
	获取阿里云 kv-store 的各项监控数据
*/
func GetKvMonitor(instanceId string, timeDimension string, metricDimension string,
	region string, accessKey string, accessSecret string) (response *cms.DescribeMetricDataResponse) {
	client, err := cms.NewClientWithAccessKey(region, accessKey, accessSecret)
	if err != nil {
		fmt.Println("aliyun sdk client init 失败: ", err.Error())
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
		fmt.Print(err.Error())
	}
	return response
}

/*
	获取阿里云 slb 的各项监控数据
*/
func GetSlbMonitor(instanceId string, timeDimension string, metricDimension string,
	region string, accessKey string, accessSecret string) (response *cms.DescribeMetricDataResponse) {
	client, err := cms.NewClientWithAccessKey(region, accessKey, accessSecret)
	if err != nil {
		fmt.Println("aliyun sdk client init 失败: ", err.Error())
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
		fmt.Print(err.Error())
	}
	return response
}
