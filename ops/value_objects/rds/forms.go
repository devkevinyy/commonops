package rds

import "github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:00 AM
 * @Desc:
 */

type RdsResp struct {
	Total uint                `json:"total"`
	Page  uint                `json:"page"`
	Rds   interface{} `json:"rds"`
}

type RdsBaseInfoForm struct {
	DbInstanceDescription string `json:"dbInstanceDescription"`
	DbInstanceStorage     string `json:"dbInstanceStorage"`
	DbMemory              string `json:"dbMemory"`
	DbExpiredTime         string `json:"dbExpiredTime"`
}

type RdsQueryForm struct {
	Page              uint    `form:"page" json:"page" binding:"required"`
	Size              uint    `form:"size" json:"size" binding:"required"`
	QueryExpiredTime  string `form:"queryExpiredTime" json:"queryExpiredTime"`
	QueryKeyword      string `form:"queryKeyword" json:"queryKeyword"`
	QueryCloudAccount uint    `form:"queryCloudAccount" json:"queryCloudAccount"`
}

type EcsBaseInfoForm struct {
	InstanceId       string `json:"instanceId"`
	InstanceName     string `json:"instanceName"`
	InnerIpAddress   string `json:"innerIpAddress"`
	PublicIpAddress  string `json:"publicIpAddress"`
	PrivateIpAddress string `json:"privateIpAddress"`
	Cpu              string `json:"cpu"`
	Memory           string `json:"memory"`
	ExpiredTime      string `json:"expiredTime"`
	SshPort          string `json:"sshPort"`
	SshUser          string `json:"sshUser"`
	SshPwd           string `json:"sshPwd"`
}

type KvBaseInfoForm struct {
	KvInstanceName string `json:"kvInstanceName"`
	KvBandwidth    string `json:"kvBandwidth"`
	KvCapacity     string `json:"kvCapacity"`
	KvExpiredTime  string `json:"kvExpiredTime"`
}

type ExtraInfoForm struct {
	Id string `json:"id"`
	EcsBaseInfoForm
	RdsBaseInfoForm
	KvBaseInfoForm
}

type RdsInfoForm struct {
	DBInstanceId          string `json:"dbInstanceId"`
	DBInstanceDescription string `json:"dbInstanceDescription"`
	Engine                string `json:"engine"`
	ConnectionString      string `json:"connectionString"`
	Port                  string `json:"port"`
	DBInstanceMemory      int    `json:"dbInstanceMemory"`
	DBInstanceStorage     int    `json:"dbInstanceStorage"`
	CreateTime            string `json:"createTime"`
	ExpireTime            string `json:"expireTime"`
	CloudAccountId        int    `json:"cloudAccountId"`
}

type ResDeleteForm struct {
	Id uint `json:"id"`
}

type DbInfo struct {
	DBInstance  rds.DBInstance
	DBAttribute rds.DBInstanceAttribute
}

