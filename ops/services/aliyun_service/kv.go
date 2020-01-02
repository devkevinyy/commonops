package aliyun_service

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/untils"
	"go.uber.org/zap"
	"math"
)

func KvDiffCacheClean() (err error) {
	err = database.MysqlClient.Exec("delete from diff_caches where type = 'kv' ").Error
	return
}

func GetKvList(accessKey string, keySecret string)(data []r_kvstore.KVStoreInstance){
	for _, region := range conf.RegionList {
		untils.Log.Info("开始获取阿里云kv数据", zap.String("accessKey", accessKey),
			zap.String("region", region))
		client, err := r_kvstore.NewClientWithAccessKey(region, accessKey, keySecret)
		if err != nil {
			fmt.Println(err)
			continue
		}
		request := r_kvstore.CreateDescribeInstancesRequest()
		request.PageNumber = requests.NewInteger(1)
		request.PageSize = requests.NewInteger(50)
		resp, err := client.DescribeInstances(request)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		if resp.TotalCount == 0 {
			untils.Log.Info("该账号在该地区无kv数据.")
			continue
		} else {
			untils.Log.Info("数据总数", zap.Int("msg", resp.TotalCount))
			data = append(data, []r_kvstore.KVStoreInstance(resp.Instances.KVStoreInstance)...)
			totalPage := int(math.Ceil(float64(resp.TotalCount) / float64(resp.PageSize)))
			untils.Log.Info("数据总页数", zap.Int("msg", totalPage))
			if totalPage > 1 {
				for i := 1; i<totalPage; i++ {
					currentPage := i + 1
					untils.Log.Info("开始获取数据页数", zap.Int("msg", currentPage))
					request.PageNumber = requests.NewInteger(currentPage)
					resp, err := client.DescribeInstances(request)
					if err != nil {
						untils.Log.Error(err.Error())
						break
					}
					data = append(data, []r_kvstore.KVStoreInstance(resp.Instances.KVStoreInstance)...)
				}
			}
		}
	}
	return data
}
