package ecs

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/11/21 5:14 PM
 * @Desc:
 */
type CloudServerResp struct {
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Servers interface{} `json:"servers"`
}

type CloudServerQueryForm struct {
	UserId            uint
	Page              int    `form:"page" json:"page" binding:"required"`
	Size              int    `form:"size" json:"size" binding:"required"`
	QueryExpiredTime  string `form:"queryExpiredTime" json:"queryExpiredTime"`
	QueryKeyword      string `form:"queryKeyword" json:"queryKeyword"`
	QueryCloudAccount int    `form:"queryCloudAccount" json:"queryCloudAccount"`
}

type ServerInfoForm struct {
	HostName        string   `json:"hostName"`
	InstanceId      string   `json:"instanceId"`
	InstanceName    string   `json:"instanceName"`
	Description     string   `json:"description"`
	InnerIpAddress  string   `json:"innerIpAddress"`
	PublicIpAddress string   `json:"publicIpAddress"`
	Memory          int      `json:"memory"`
	Disk            int      `json:"disk"`
	CreateTime      string   `json:"createTime"`
	ExpiredTime     string   `json:"expiredTime"`
	CloudAccountId  uint     `json:"cloudAccountId"`
	OsType          string   `json:"osType"`
	Cpu             int      `json:"cpu"`
	SshPort         int      `json:"sshPort"`
	SshUser         string   `json:"sshUser"`
	SshPwd          string   `json:"sshPwd"`
	Tags            []string `json:"tags"`
}

type EcsBaseInfoForm struct {
	InstanceId       string   `json:"instanceId"`
	InstanceName     string   `json:"instanceName"`
	InnerIpAddress   string   `json:"innerIpAddress"`
	PublicIpAddress  string   `json:"publicIpAddress"`
	PrivateIpAddress string   `json:"privateIpAddress"`
	Cpu              string   `json:"cpu"`
	Memory           string   `json:"memory"`
	ExpiredTime      string   `json:"expiredTime"`
	SshPort          string   `json:"sshPort"`
	SshUser          string   `json:"sshUser"`
	SshPwd           string   `json:"sshPwd"`
	Tags             []string `json:"tags"`
}

type RdsBaseInfoForm struct {
	DbInstanceDescription string `json:"dbInstanceDescription"`
	DbInstanceStorage     string `json:"dbInstanceStorage"`
	DbMemory              string `json:"dbMemory"`
	DbExpiredTime         string `json:"dbExpiredTime"`
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

type ResDeleteForm struct {
	Id uint `json:"id"`
}

type SshForm struct {
	Token    string `form:"token" json:"token"`
	ServerId string `form:"serverId" json:"serverId"`
}

type BatchSshFrom struct {
	Ids     []uint `json:"ids"`
	Command string `json:"command"`
}
