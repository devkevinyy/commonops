package aliyun_service

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/untils"
	"go.uber.org/zap"
	"math"
)

func SlbDiffCacheClean() (err error) {
	err = database.MysqlClient.Exec("delete from diff_caches where type = 'slb' ").Error
	return
}

func GetSlbList(accessKey string, keySecret string)(data []slb.LoadBalancer){
	for _, region := range conf.RegionList {
		untils.Log.Info("开始获取阿里云slb数据", zap.String("accessKey", accessKey),
			zap.String("region", region))
		client, err := slb.NewClientWithAccessKey(region, accessKey, keySecret)
		if err != nil {
			fmt.Println(err)
			continue
		}
		request := slb.CreateDescribeLoadBalancersRequest()
		request.PageNumber = requests.NewInteger(1)
		request.PageSize = requests.NewInteger(50)
		resp, err := client.DescribeLoadBalancers(request)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		if resp.TotalCount == 0 {
			untils.Log.Info("该账号在该地区无slb数据.")
			continue
		} else {
			untils.Log.Info("数据总数", zap.Int("msg", resp.TotalCount))
			data = append(data, []slb.LoadBalancer(resp.LoadBalancers.LoadBalancer)...)
			totalPage := int(math.Ceil(float64(resp.TotalCount) / float64(resp.PageSize)))
			untils.Log.Info("数据总页数", zap.Int("msg", totalPage))
			if totalPage > 1 {
				for i := 1; i<totalPage; i++ {
					currentPage := i + 1
					untils.Log.Info("开始获取数据页数", zap.Int("msg", currentPage))
					request.PageNumber = requests.NewInteger(currentPage)
					resp, err := client.DescribeLoadBalancers(request)
					if err != nil {
						untils.Log.Error(err.Error())
						break
					}
					data = append(data, []slb.LoadBalancer(resp.LoadBalancers.LoadBalancer)...)
				}
			}
		}
	}
	return data
}
