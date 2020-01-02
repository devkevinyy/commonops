package cron

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/database"
	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/services/aliyun_service"
	"github.com/chujieyang/commonops/ops/untils"
	"github.com/robfig/cron"
	"strings"
)

func init(){
	untils.Log.Info("cron init exec ...")
	c := cron.New()
	err := c.AddFunc("0 0 23 * *", SyncAliYunEcsData)
	err = c.AddFunc("0 10 23 * *", SyncAliYunRdsData)
	err = c.AddFunc("0 20 23 * *", SyncAliYunKvData)
	err = c.AddFunc("0 30 23 * *", SyncAliYunSlbData)
	if err != nil {
		untils.Log.Error(err.Error())
	}
	c.Start()
}


/**
	同步阿里云的 ecs 服务器数据
 */
func SyncAliYunEcsData(){
	untils.Log.Info("定时任务 [同步ecs数据]")
	err := aliyun_service.EcsDiffCacheClean()
	if err != nil {
		untils.Log.Error(err.Error())
	}
	accountList := models.GetAllCloudAccounts()
	for _, account := range accountList {
		accountId := account.Id
		accessKey := account.Key
		accessSecret := account.Secret
		if accessKey == "" || accessKey == "未知" {
			continue
		}
		untils.Log.Info(fmt.Sprintf("同步云账号信息：%s - %s", account.Name, account.Key))
		data := aliyun_service.GetEcsList(accessKey, accessSecret)
		for _, v := range data {
			ecs := models.Ecs{
				ImageId: v.ImageId,
				InstanceType: v.InstanceType,
				InstanceNetworkType: v.InstanceNetworkType,
				LocalStorageAmount: fmt.Sprintf("%d", v.LocalStorageAmount),
				InstanceChargeType: v.InstanceChargeType,
				ClusterId: v.ClusterId,
				InstanceName: v.InstanceName,
				StartTime: untils.ConvertUtcTimeToLocal(v.StartTime, "2006-01-02T15:04Z"),
				ZoneId: v.ZoneId,
				InternetChargeType: v.InternetChargeType,
				InternetMaxBandwidthIn: fmt.Sprintf("%d", v.InternetMaxBandwidthIn),
				HostName: v.HostName,
				Cpu: fmt.Sprintf("%d", v.Cpu),
				Status: v.Status,
				OSName: v.OSName,
				OSNameEn: v.OSNameEn,
				SerialNumber: v.SerialNumber,
				RegionId: v.RegionId,
				InternetMaxBandwidthOut: fmt.Sprintf("%d", v.InternetMaxBandwidthOut),
				ResourceGroupId: v.ResourceGroupId,
				InstanceTypeFamily: v.InstanceTypeFamily,
				InstanceId: v.InstanceId,
				DeploymentSetId: v.DeploymentSetId,
				Description: v.Description,
				ExpiredTime: untils.ConvertUtcTimeToLocal(v.ExpiredTime, "2006-01-02T15:04Z"),
				OSType: v.OSType,
				Memory: fmt.Sprintf("%d", v.Memory),
				CreationTime: untils.ConvertUtcTimeToLocal(v.CreationTime, "2006-01-02T15:04Z"),
				LocalStorageCapacity: fmt.Sprintf("%d", v.LocalStorageCapacity),
				InnerIpAddress: strings.Join(v.InnerIpAddress.IpAddress, ","),
				PublicIpAddress: strings.Join(v.PublicIpAddress.IpAddress, ","),
				PrivateIpAddress: strings.Join(v.VpcAttributes.PrivateIpAddress.IpAddress, ","),
				DataStatus: 1,
			}
			models.SaveOrUpdateEcs(v.InstanceId, accountId, ecs)
			untils.Log.Info(fmt.Sprintf("[cron jobs - SyncAliYunEcsData]: %s, %s, %s", v.HostName,
				v.InstanceId, v.InstanceName))
		}
	}
}

/**
	同步阿里云的 Rds 服务器数据
*/
func SyncAliYunRdsData(){
	untils.Log.Info("定时任务 [同步rds数据]")
	err := aliyun_service.RdsDiffCacheClean()
	if err != nil {
		untils.Log.Error(err.Error())
	}
	accountList := models.GetAllCloudAccounts()
	for _, account := range accountList {
		accountId := account.Id
		accessKey := account.Key
		accessSecret := account.Secret
		if accessKey == "" || accessKey == "未知" {
			continue
		}
		untils.Log.Info(fmt.Sprintf("[SyncAliYunRdsData] 同步云账号信息：%s, %s", account.Name, account.Key))
		data := aliyun_service.GetRdsList(accessKey, accessSecret)
		for _, v := range data {
			untils.Log.Info(fmt.Sprintf("[SyncAliYunRdsData] rds: %s", v.DBInstance.DBInstanceId))
			attribute := v.DBAttribute
			rds := models.Rds{
				DataStatus: 1,
				InsId: v.DBInstance.InsId,
				DBInstanceId: v.DBInstance.DBInstanceId,
				DBInstanceDescription: v.DBInstance.DBInstanceDescription,
				PayType: v.DBInstance.PayType,
				DBInstanceType: v.DBInstance.DBInstanceType,
				RegionId: v.DBInstance.RegionId,
				ExpireTime: untils.ConvertUtcTimeToLocal(v.DBInstance.ExpireTime, "2006-01-02T15:04:05Z"),
				DestroyTime: untils.ConvertUtcTimeToLocal(v.DBInstance.DestroyTime, "2006-01-02T15:04:05Z"),
				DBInstanceStatus: v.DBInstance.DBInstanceStatus,
				Engine: v.DBInstance.Engine,
				DBInstanceNetType: v.DBInstance.DBInstanceNetType,
				ConnectionMode: v.DBInstance.ConnectionMode,
				LockMode: v.DBInstance.LockMode,
				Category: v.DBInstance.Category,
				DBInstanceStorageType: v.DBInstance.DBInstanceStorageType,
				DBInstanceClass: v.DBInstance.DBInstanceClass,
				InstanceNetworkType: v.DBInstance.InstanceNetworkType,
				VpcCloudInstanceId: v.DBInstance.VpcCloudInstanceId,
				LockReason: v.DBInstance.LockReason,
				ZoneId: v.DBInstance.ZoneId,
				MutriORsignle: v.DBInstance.MutriORsignle,
				CreateTime: untils.ConvertUtcTimeToLocal(v.DBInstance.CreateTime, "2006-01-02T15:04:05Z"),
				EngineVersion: v.DBInstance.EngineVersion,
				GuardDBInstanceId: v.DBInstance.GuardDBInstanceId,
				TempDBInstanceId: v.DBInstance.TempDBInstanceId,
				MasterInstanceId: v.DBInstance.MasterInstanceId,
				VpcId: v.DBInstance.VpcId,
				VSwitchId: v.DBInstance.VSwitchId,
				ReplicateId: v.DBInstance.ReplicateId,
				ResourceGroupId: v.DBInstance.ResourceGroupId,
				ConnectionString: attribute.ConnectionString,
				Port: attribute.Port,
				DBInstanceMemory: attribute.DBInstanceMemory,
				DBInstanceStorage: attribute.DBInstanceStorage,
			}
			isSuccess := models.SaveOrUpdateRds(v.DBInstance.DBInstanceId, accountId, rds)
			untils.Log.Info(fmt.Sprintf("%s - %s - %s - %s - %s ", v.DBInstance.DBInstanceId,
				v.DBInstance.DBInstanceDescription, v.DBInstance.DBInstanceStatus,
				attribute.ConnectionString, isSuccess))
		}
	}
}

/**
	同步阿里云的 KvStore 服务器数据
*/
func SyncAliYunKvData() {
	untils.Log.Info("[SyncAliYunKvData] 定时任务 [同步kv数据]")
	err := aliyun_service.KvDiffCacheClean()
	if err != nil {
		untils.Log.Error(err.Error())
	}
	accountList := models.GetAllCloudAccounts()
	for _, account := range accountList {
		accountId := account.Id
		accessKey := account.Key
		accessSecret := account.Secret
		if accessKey == "" || accessKey == "未知" {
			continue
		}
		untils.Log.Info(fmt.Sprintf("[SyncAliYunKvData] 同步云账号信息：%s, %s", account.Name, account.Key))
		data := aliyun_service.GetKvList(accessKey, accessSecret)
		for _, v := range data {
			kv := models.Kv{
				DataStatus: 1,
				InstanceClass: v.InstanceClass,
				PackageType: v.PackageType,
				ChargeType: v.ChargeType,
				ConnectionDomain: v.ConnectionDomain,
				SearchKey: v.SearchKey,
				CreateTime: untils.ConvertUtcTimeToLocal(v.CreateTime, "2006-01-02T15:04:05Z"),
				HasRenewChangeOrder: v.HasRenewChangeOrder,
				InstanceType: v.InstanceType,
				DestroyTime: untils.ConvertUtcTimeToLocal(v.DestroyTime, "2006-01-02T15:04:05Z"),
				RegionId: v.RegionId,
				PrivateIp: v.PrivateIp,
				InstanceId: v.InstanceId,
				InstanceStatus: v.InstanceStatus,
				Bandwidth: v.Bandwidth,
				NetworkType: v.NetworkType,
				VpcId: v.VpcId,
				NodeType: v.NodeType,
				Connections: v.Connections,
				ArchitectureType: v.ArchitectureType,
				ReplacateId: v.ReplacateId,
				EngineVersion: v.EngineVersion,
				Capacity: v.Capacity,
				VSwitchId: v.VSwitchId,
				InstanceName: v.InstanceName,
				Port: v.Port,
				ZoneId: v.ZoneId,
				EndTime: untils.ConvertUtcTimeToLocal(v.EndTime, "2006-01-02T15:04:05Z"),
				QPS: v.QPS,
				UserName: v.UserName,
				Config: v.Config,
				IsRds: v.IsRds,
				ConnectionMode: v.ConnectionMode,
			}
			isSuccess := models.SaveOrUpdateKv(v.InstanceId, accountId, kv)
			untils.Log.Info(fmt.Sprintf("%s - %s - %s - %s", v.InstanceId, v.InstanceName, v.InstanceClass, isSuccess))
		}
	}
}


/**
	同步阿里云的 SLB 数据
*/
func SyncAliYunSlbData() {
	untils.Log.Info("[SyncAliYunSlbData] 定时任务 [同步slb数据]")
	err := aliyun_service.SlbDiffCacheClean()
	if err != nil {
		untils.Log.Error(err.Error())
	}
	accountList := models.GetAllCloudAccounts()
	for _, account := range accountList {
		accountId := account.Id
		accessKey := account.Key
		accessSecret := account.Secret
		if accessKey == "" || accessKey == "未知" {
			continue
		}
		untils.Log.Info(fmt.Sprintf("[SyncAliYunSlbData] 同步云账号信息：%s, %s", account.Name, account.Key))
		data := aliyun_service.GetSlbList(accessKey, accessSecret)
		for _, v := range data {
			slb := models.Slb{
				DataStatus: 1,
				Count: v.Count,
				SlaveZoneId: v.SlaveZoneId,
				LoadBalancerStatus: v.LoadBalancerStatus,
				VSwitchId: v.VSwitchId,
				MasterZoneId: v.MasterZoneId,
				PayType: v.PayType,
				RegionIdAlias: v.RegionIdAlias,
				CreateTime: untils.ConvertUtcTimeToLocal(v.CreateTime, "2006-01-02T15:04Z"),
				Address: v.Address,
				LoadBalancerId: v.LoadBalancerId,
				AddressIPVersion: v.AddressIPVersion,
				RegionId: v.RegionId,
				ResourceGroupId: v.ResourceGroupId,
				LoadBalancerName: v.LoadBalancerName,
				InternetChargeType: v.InternetChargeType,
				AddressType: v.AddressType,
				VpcId: v.VpcId,
				NetworkType: v.NetworkType,
			}
			isSuccess := models.SaveOrUpdateSlb(v.LoadBalancerId, accountId, slb)
			untils.Log.Info(fmt.Sprintf("%s - %s - %s - %s", v.LoadBalancerId, v.LoadBalancerName,
				v.LoadBalancerStatus, isSuccess))
		}
	}
}

func DiffDataCheck() {
	type Item struct {
		InstanceId string
	}
	var resList []Item

	// ecs data check
	database.MysqlClient.Raw("select instance_id from ecs inner join cloud_account as ca " +
		"on ecs.cloud_account_id = ca.id where ca.cloud_type = '阿里云' and ecs.data_status = 1").Scan(&resList)
	for _, item := range resList {
		var count = 0
		database.MysqlClient.Raw("select count(*) from diff_caches where " +
			"type = 'ecs' and instance_id = ? limit 1 ", item.InstanceId).Count(&count)
		if count == 0 {
			err := database.MysqlClient.Exec("update ecs set data_status = 0 where " +
				"instance_id = ?", item.InstanceId).Error
			if err != nil {
				untils.Log.Error(err.Error())
			}
		}
	}

	// rds data check
	database.MysqlClient.Raw("select db_instance_id as instance_id from rds inner join cloud_account as ca " +
		"on rds.cloud_account_id = ca.id where ca.cloud_type = '阿里云' and rds.data_status = 1").Scan(&resList)
	for _, item := range resList {
		var count = 0
		database.MysqlClient.Raw("select count(*) from diff_caches where " +
			"type = 'rds' and instance_id = ? limit 1 ", item.InstanceId).Count(&count)
		if count == 0 {
			err := database.MysqlClient.Exec("update rds set data_status = 0 where " +
				"db_instance_id = ?", item.InstanceId).Error
			if err != nil {
				untils.Log.Error(err.Error())
			}
		}
	}

	// redis data check
	database.MysqlClient.Raw("select instance_id from kv inner join cloud_account as ca " +
		"on kv.cloud_account_id = ca.id where ca.cloud_type = '阿里云' and kv.data_status = 1").Scan(&resList)
	for _, item := range resList {
		var count = 0
		database.MysqlClient.Raw("select count(*) from diff_caches where " +
			"type = 'kv' and instance_id = ? limit 1 ", item.InstanceId).Count(&count)
		if count == 0 {
			err := database.MysqlClient.Exec("update kv set data_status = 0 where " +
				"instance_id = ?", item.InstanceId).Error
			if err != nil {
				untils.Log.Error(err.Error())
			}
		}
	}

	// slb data check
	database.MysqlClient.Raw("select load_balancer_id as instance_id from slb inner join cloud_account as ca " +
		"on slb.cloud_account_id = ca.id where ca.cloud_type = '阿里云' and slb.data_status = 1").Scan(&resList)
	for _, item := range resList {
		var count = 0
		database.MysqlClient.Raw("select count(*) from diff_caches where " +
			"type = 'slb' and instance_id = ? limit 1 ", item.InstanceId).Count(&count)
		if count == 0 {
			err := database.MysqlClient.Exec("update slb set data_status = 0 where " +
				"load_balancer_id = ?", item.InstanceId).Error
			if err != nil {
				untils.Log.Error(err.Error())
			}
		}
	}

}
