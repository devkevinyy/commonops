package aliyun_service

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/untils"
	"go.uber.org/zap"
	"math"
)

func EcsDiffCacheClean() (err error) {
	err = database.MysqlClient.Exec("delete from diff_caches where type = 'ecs' ").Error
	return
}

func GetEcsList(accessKey string, keySecret string)(data []ecs.Instance){
	for _, region := range conf.RegionList {
		untils.Log.Info("开始获取阿里云ecs数据", zap.String("accesskey", accessKey),
			zap.String("region", region))
		client, err := ecs.NewClientWithAccessKey(region, accessKey, keySecret)
		if err != nil {
			fmt.Println(err)
			continue
		}
		request := ecs.CreateDescribeInstancesRequest()
		request.PageNumber = requests.NewInteger(1)
		request.PageSize = requests.NewInteger(50)
		resp, err := client.DescribeInstances(request)
		if err != nil {
			untils.Log.Error(err.Error())
			continue
		}
		if resp.TotalCount == 0 {
			untils.Log.Info("该账号在该地区无ecs数据.")
			continue
		} else {
			untils.Log.Info("数据总数", zap.Int("msg", resp.TotalCount))
			data = append(data, []ecs.Instance(resp.Instances.Instance)...)
			totalPage := int(math.Ceil(float64(resp.TotalCount) / float64(resp.PageSize)))
			untils.Log.Info("数据总页数", zap.Int("msg", totalPage))
			if totalPage > 1 { // 获取第一页之后的数据
				for i := 1; i<totalPage; i++ {
					currentPage := i + 1
					untils.Log.Info("开始获取数据页数", zap.Int("msg", currentPage))
					request.PageNumber = requests.NewInteger(currentPage)
					resp, err := client.DescribeInstances(request)
					if err != nil {
						untils.Log.Error(err.Error())
						break
					}
					data = append(data, []ecs.Instance(resp.Instances.Instance)...)
				}
			}
		}
	}
	return data
}
