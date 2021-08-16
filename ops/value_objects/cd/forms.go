package cd

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 3:11 PM
 * @Desc:
 */

type CdProcessTemplateReq struct {
	JobName       string `json:"jobName" form:"jobName" required:"true"`
	TemplateName  string `json:"templateName" form:"templateName" required:"true"`
	ClusterId     string `json:"clusterId" form:"clusterId" required:"true"`
	Namespace     string `json:"namespace" form:"namespace" required:"true"`
	ImageName     string `json:"imageName" form:"imageName" required:"true"`
	DeployYaml    string `json:"deployYaml" form:"deployYaml" required:"true"`
	ServiceYaml   string `json:"serviceYaml" form:"serviceYaml" required:"true"`
	ConfigmapYaml string `json:"configmapYaml" form:"configmapYaml" required:"true"`
	IngressYaml   string `json:"ingressYaml" form:"ingressYaml" required:"true"`
}

type CdProcessTemplateLog struct {
	TemplateId int    `json:"templateId" form:"templateId" required:"true"`
	ImageName  string `json:"imageName" form:"imageName" required:"true"`
}

type GetCdProcessLogReq struct {
	Page             int `json:"page" form:"page" required:"true"`
	Size             int `json:"size" form:"size" required:"true"`
	QueryProjectName int `json:"queryProjectName" form:"queryProjectName"`
}

type CdProcessLogResp struct {
	Total int                         `json:"total"`
	Page  int                         `json:"page"`
	Logs  interface{} `json:"logs"`
}