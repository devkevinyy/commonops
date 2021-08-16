package ci

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 3:11 PM
 * @Desc:
 */

type JenkinsCreateJobReq struct {
	JobName                string `json:"jobName" form:"jobName" required:"true"`
	GitRepo                string `json:"gitRepo" form:"gitRepo" required:"true"`
	GitBranch              string `json:"gitBranch" form:"gitBranch" required:"true"`
	GitCredentials         string `json:"gitCredentials" form:"gitCredentials" required:"true"`
	DockerImageRepo        string `json:"dockerImageRepo" form:"dockerImageRepo" required:"true"`
	DockerImageCredentials string `json:"dockerImageCredentials" form:"dockerImageCredentials" required:"true"`
	DockerImageName        string `json:"dockerImageName" form:"dockerImageName" required:"true"`
}

type JenkinsUpdateJobReq struct {
	JobName        string `json:"jobName" form:"jobName" required:"true"`
	PipelineScript string `json:"pipelineScript" form:"pipelineScript" required:"true"`
}

type JenkinsBuildListReq struct {
	JobName string `json:"jobName" form:"jobName" required:"true"`
}

type JenkinsCancelBuildReq struct {
	JobName  string `json:"jobName" form:"jobName" required:"true"`
	BuildNum int32  `json:"buildNum" form:"buildNum" required:"true"`
}

type JenkinsBuildInfoReq struct {
	JobName string `json:"jobName" form:"jobName" required:"true"`
	Number  int32  `json:"number" form:"number" required:"true"`
}

type JenkinsBuildStageLogReq struct {
	JobName  string `json:"jobName" form:"jobName" required:"true"`
	BuildNum int32  `json:"buildNum" form:"buildNum" required:"true"`
	NodeId   int32  `json:"nodeId" form:"nodeId" required:"true"`
}

type JenkinsCreateCredentialsReq struct {
	CredentialId string `json:"credentialId" form:"credentialId" required:"true"`
	Username     string `json:"username" form:"username" required:"true"`
	Password     string `json:"password" form:"password" required:"true"`
	Description  string `json:"description" form:"description" required:"true"`
}

type JenkinsGetCredentialsReq struct {
	CredentialId string `json:"credentialId" form:"credentialId" required:"true"`
}
