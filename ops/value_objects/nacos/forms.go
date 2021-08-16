package nacos

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 7/20/21 3:09 PM
 * @Desc:
 */

type NacosServer struct {
	Alias    string `json:"alias" form:"alias" required:"true"`
	EndPoint string `json:"endpoint" form:"endpoint" required:"true"`
	Username string `json:"username" form:"username" required:"true"`
	Password string `json:"password" form:"password" required:"true"`
}

type NacosNsListReq struct {
	ClusterId string `json:"clusterId" form:"clusterId" required:"true"`
}

type NacosConfigsReq struct {
	ClusterId  string `json:"clusterId" form:"clusterId" required:"true"`
	Namespace  string `json:"namespace" form:"namespace" required:"true"`
	Page       int    `json:"page" form:"page" required:"true"`
	Size       int    `json:"size" form:"size" required:"true"`
	ConfigTags string `json:"configTags" form:"configTags"`
}

type CreateNacosConfigReq struct {
	ClusterId  string `json:"clusterId" form:"clusterId" required:"true"`
	TemplateId int    `json:"templateId" form:"templateId" required:"true"`
	Namespace  string `json:"namespace" form:"namespace" required:"true"`
	DataId     string `json:"dataId" form:"dataId" required:"true"`
	Group      string `json:"group" form:"group" required:"true"`
	// Content    string `json:"content" form:"content" required:"true"`
	// ConfigType string `json:"configType" form:"configType" required:"true"`
	// ConfigTags string `json:"configTags" form:"configTags" required:"true"`
}

type GetNacosConfigDetailReq struct {
	ClusterId string `json:"clusterId" form:"clusterId" required:"true"`
	Namespace string `json:"namespace" form:"namespace" required:"true"`
	DataId    string `json:"dataId" form:"dataId" required:"true"`
	Group     string `json:"group" form:"group" required:"true"`
}

type UpdateNacosConfigReq struct {
	ClusterId  string `json:"clusterId" form:"clusterId" required:"true"`
	Id         int    `json:"id" form:"id" required:"true"`
	ConfigId   string `json:"configId" form:"configId" required:"true"`
	Namespace  string `json:"namespace" form:"namespace" required:"true"`
	DataId     string `json:"dataId" form:"dataId" required:"true"`
	Group      string `json:"group" form:"group" required:"true"`
	Content    string `json:"content" form:"content" required:"true"`
	ConfigType string `json:"configType" form:"configType" required:"true"`
	ConfigTags string `json:"configTags" form:"configTags" required:"true"`
}

type DeleteNacosConfigReq struct {
	ClusterId string `json:"clusterId" form:"idclusterId" required:"true"`
	Id        int `json:"id" form:"id" required:"true"`
	Namespace string `json:"namespace" form:"namespace" required:"true"`
	DataId    string `json:"dataId" form:"dataId" required:"true"`
	Group     string `json:"group" form:"group" required:"true"`
}

type CreateNacosConfigCopyReq struct {
	ClusterId    string `json:"clusterId" form:"clusterId" required:"true"`
	SrcNamespace string `json:"srcNamespace" form:"srcNamespace" required:"true"`
	SrcDataId    string `json:"srcDataId" form:"srcDataId" required:"true"`
	SrcGroup     string `json:"srcGroup" form:"srcGroup" required:"true"`
	DstNamespace string `json:"dstNamespace" form:"dstNamespace" required:"true"`
	DstDataId    string `json:"dstDataId" form:"dstDataId" required:"true"`
	DstGroup     string `json:"dstGroup" form:"dstGroup" required:"true"`
}

type SyncDstConfig struct {
	Namespace string `json:"namespace" form:"namespace" required:"true"`
	DataId    string `json:"dataId" form:"dataId" required:"true"`
	Group     string `json:"group" form:"group" required:"true"`
}

type CreateNacosConfigSyncReq struct {
	ClusterId    string          `json:"clusterId" form:"clusterId" required:"true"`
	SrcNamespace string          `json:"srcNamespace" form:"srcNamespace" required:"true"`
	SrcDataId    string          `json:"srcDataId" form:"srcDataId" required:"true"`
	SrcGroup     string          `json:"srcGroup" form:"srcGroup" required:"true"`
	DstConfigs   []SyncDstConfig `json:"dstConfigs" form:"dstConfigs" required:"true"`
}

type ConfigTemplatesReq struct {
	Name string `json:"name" form:"name"`
	Page int    `form:"page" json:"page" binding:"required"`
	Size int    `form:"size" json:"size" binding:"required"`
}

type ConfigTemplateAddReq struct {
	Name          string   `json:"name" form:"name" binding:"required"`
	ConfigContent string   `form:"configContent" json:"configContent" binding:"required"`
	FillField     []string `form:"fillField" json:"fillField" binding:"required"`
}

type ConfigTemplateUpdateReq struct {
	Id            int      `json:"id" form:"id" binding:"required"`
	Name          string   `json:"name" form:"name" binding:"required"`
	ConfigContent string   `form:"configContent" json:"configContent" binding:"required"`
	FillField     []string `form:"fillField" json:"fillField" binding:"required"`
}

type ConfigTemplateDeleteReq struct {
	Id int `json:"id" form:"id" binding:"required"`
}
