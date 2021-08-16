package service

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	aliEcs "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/opslog"
	"github.com/chujieyang/commonops/ops/utils"
	ecs2 "github.com/chujieyang/commonops/ops/value_objects/ecs"
	"go.uber.org/zap"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/11/21 5:16 PM
 * @Desc:
 */

type ecs struct{}

var ecsInstance = &ecs{}

func GetEcsService() *ecs {
	return ecsInstance
}

func (e *ecs) GetEcsDataByPage(params ecs2.CloudServerQueryForm) (total int, data []models.EcsDetail) {
	offset := (params.Page - 1) * params.Size
	total = models.GetServersCount(params.UserId, params.QueryExpiredTime, params.QueryKeyword,
		params.QueryCloudAccount)
	data = models.GetServers(params.UserId, offset, params.Size,
		params.QueryExpiredTime, params.QueryKeyword, params.QueryCloudAccount)
	return
}

func (e *ecs) AddCloudServer(params ecs2.ServerInfoForm) error {
	return models.AddCloudServer(params)
}

func (e *ecs) UpdateCloudServer(params ecs2.ExtraInfoForm) error {
	return models.UpdateCloudServer(params)
}

func (e *ecs) DeleteCloudServer(id uint) error {
	return models.DeleteCloudServer(id)
}

func (e *ecs) GetServerDetail(id uint) models.EcsDetail {
	return models.GetServerDetail(id)
}

func (e *ecs) EcsDiffCacheClean() error {
	return models.CleanEcsDiffCaches()
}

func (e *ecs) GetServersTreeData() ([]models.Ecs, error) {
	return models.GetEcsTreeData()
}

type ecsClient struct {
	ecs *aliEcs.Client
}

func NewEcsClient(region string, accessKey string, keySecret string) (client *ecsClient, err error) {
	if c, err := aliEcs.NewClientWithAccessKey(region, accessKey, keySecret); err == nil {
		client = &ecsClient{
			ecs: c,
		}
	}
	return
}

func (e *ecs) GetEcsList(accessKey string, keySecret string) (data []aliEcs.Instance) {
	for _, region := range conf.RegionList {
		opslog.Info().Println("开始获取阿里云ecs数据", zap.String("accesskey", accessKey),
			zap.String("region", region))
		client, err := NewEcsClient(region, accessKey, keySecret)
		if err != nil {
			opslog.Error().Println(err)
			continue
		}
		request := aliEcs.CreateDescribeInstancesRequest()
		request.PageNumber = requests.NewInteger(1)
		request.PageSize = requests.NewInteger(50)
		resp, err := client.ecs.DescribeInstances(request)
		if err != nil {
			opslog.Error().Println(err)
			continue
		}
		if resp.TotalCount == 0 {
			continue
		} else {
			data = append(data, []aliEcs.Instance(resp.Instances.Instance)...)
			totalPage := int(math.Ceil(float64(resp.TotalCount) / float64(resp.PageSize)))
			if totalPage > 1 { // 获取第一页之后的数据
				for i := 1; i < totalPage; i++ {
					currentPage := i + 1
					request.PageNumber = requests.NewInteger(currentPage)
					resp, err := client.ecs.DescribeInstances(request)
					if err != nil {
						opslog.Error().Println(err.Error())
						break
					}
					data = append(data, []aliEcs.Instance(resp.Instances.Instance)...)
				}
			}
		}
	}
	return data
}

/*
	获取阿里云 ecs 的各项监控数据
*/
func (e *ecs) GetEcsMonitor(instanceId string, timeDimension string, metricDimension string,
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
	request.Namespace = "acs_ecs_dashboard"
	request.Period = period
	request.MetricName = metricDimension

	response, err = client.DescribeMetricData(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	return response
}

func (e *ecs) BatchSshExec(ids []uint, cmd string) (result []map[string]string, err error) {
	var wg sync.WaitGroup
	wg.Add(len(ids))
	resultChannel := make(chan map[string]string)
	for _, id := range ids {
		go execSshCommand(id, cmd, &resultChannel, &wg)
	}
	for i := 0; i < len(ids); i++ {
		result = append(result, <-resultChannel)
	}
	wg.Wait()
	return
}

func execSshCommand(id uint, cmd string, resultChannel *chan map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	ecsInfo := models.GetServerDetail(id)
	conn, err1 := Connect(fmt.Sprintf("%s:%d", ecsInfo.PublicIpAddress, ecsInfo.SshPort), ecsInfo.SshUser, utils.DesDecode(ecsInfo.SshPwd))
	if err1 != nil {
		*resultChannel <- map[string]string{
			fmt.Sprintf("%s: %s", ecsInfo.InstanceName, ecsInfo.PublicIpAddress): err1.Error(),
		}
		return
	}
	output, err1 := conn.SendCommands(cmd)
	if err1 != nil {
		*resultChannel <- map[string]string{
			fmt.Sprintf("%s: %s", ecsInfo.InstanceName, ecsInfo.PublicIpAddress): err1.Error(),
		}
		return
	}
	*resultChannel <- map[string]string{
		fmt.Sprintf("%s: %s", ecsInfo.InstanceName, ecsInfo.PublicIpAddress): string(output),
	}
	return
}
