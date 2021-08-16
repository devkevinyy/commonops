package kv

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 11:23 AM
 * @Desc:
 */

type KvResp struct {
	Total uint               `json:"total"`
	Page  uint               `json:"page"`
	Kv    interface{} `json:"kv"`
}

type KvBaseInfoForm struct {
	KvInstanceName string `json:"kvInstanceName"`
	KvBandwidth    string `json:"kvBandwidth"`
	KvCapacity     string `json:"kvCapacity"`
	KvExpiredTime  string `json:"kvExpiredTime"`
}

type KvQueryForm struct {
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

type RdsBaseInfoForm struct {
	DbInstanceDescription string `json:"dbInstanceDescription"`
	DbInstanceStorage     string `json:"dbInstanceStorage"`
	DbMemory              string `json:"dbMemory"`
	DbExpiredTime         string `json:"dbExpiredTime"`
}

type ExtraInfoForm struct {
	Id string `json:"id"`
	EcsBaseInfoForm
	RdsBaseInfoForm
	KvBaseInfoForm
}

type KvInfoForm struct {
	InstanceId       string `json:"instanceId"`
	InstanceName     string `json:"instanceName"`
	InstanceType     string `json:"instanceType"`
	ConnectionString string `json:"connectionString"`
	Port             string `json:"port"`
	Capacity         int    `json:"capacity"`
	Bandwidth        int    `json:"bandwidth"`
	CreateTime       string `json:"createTime"`
	EndTime          string `json:"endTime"`
	CloudAccountId   uint    `json:"cloudAccountId"`
}

type ResDeleteForm struct {
	Id uint `json:"id"`
}