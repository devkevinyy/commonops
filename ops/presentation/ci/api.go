package ci

import (
	"github.com/chujieyang/commonops/ops/value_objects/ci"
	"net/http"

	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/services/jenkins_service"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
)

// 获取所有Job
func IGetJenkinsJobList(c *gin.Context) {
	jobsInfo, err := jenkins_service.JenkinsInstance.GetJobList()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: jobsInfo})
}

// 创建构建job
func IPostJenkinsJob(c *gin.Context) {
	var req ci.JenkinsCreateJobReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	ciParams := models.CiConfigInfo{
		GitRepo:                 req.GitRepo,
		GitBranch:               req.GitBranch,
		GitCredentialId:         req.GitCredentials,
		DockerImageRepo:         req.DockerImageRepo,
		DockerImageCredentialId: req.DockerImageCredentials,
		DockerImageName:         req.DockerImageName,
	}
	buildInfo, err := jenkins_service.JenkinsInstance.CreateJobItem(req.JobName, ciParams)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 更新构建job
func IPutJenkinsJob(c *gin.Context) {
	var req ci.JenkinsUpdateJobReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	ciParams := models.CiScriptInfo{
		PipelineScript: req.PipelineScript,
	}
	buildInfo, err := jenkins_service.JenkinsInstance.UpdateJobItem(req.JobName, ciParams)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 获取构建Job的配置文件
func IGetJenkinsJobConfig(c *gin.Context) {
	var req ci.JenkinsBuildListReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	configInfo, err := jenkins_service.JenkinsInstance.GetJobItemConfig(req.JobName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: configInfo})
}

// 删除构建job
func IDeleteJenkinsJob(c *gin.Context) {
	var req ci.JenkinsBuildListReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	jobInfo, err := jenkins_service.JenkinsInstance.DeleteJobItem(req.JobName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: jobInfo})
}

// 获取构建列表
func IGetJenkinsBuildList(c *gin.Context) {
	var req ci.JenkinsBuildListReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	jobInfo, err := jenkins_service.JenkinsInstance.GetJobBuildList(req.JobName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: jobInfo})
}

// 获取构建详情
func IGetJenkinsBuildInfo(c *gin.Context) {
	var req ci.JenkinsBuildInfoReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	buildInfo, err := jenkins_service.JenkinsInstance.GetBuildInfo(req.JobName, req.Number)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 获取构建阶段列表
func IGetJenkinsBuildStageList(c *gin.Context) {
	var req ci.JenkinsBuildInfoReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	buildInfo, err := jenkins_service.JenkinsInstance.GetBuildNumStageList(req.JobName, req.Number)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 获取构建子阶段详情日志
func IGetJenkinsBuildStageDetailLog(c *gin.Context) {
	var req ci.JenkinsBuildStageLogReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	buildInfo, err := jenkins_service.JenkinsInstance.GetBuildNumStageDetailLog(req.JobName, req.BuildNum, req.NodeId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 触发Job进行构建
func IPostJenkinsBuild(c *gin.Context) {
	var req ci.JenkinsBuildListReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	buildInfo, err := jenkins_service.JenkinsInstance.PostJobBuild(req.JobName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 取消Job构建
func IPostCancelJenkinsBuild(c *gin.Context) {
	var req ci.JenkinsCancelBuildReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	buildInfo, err := jenkins_service.JenkinsInstance.StopJobBuild(req.JobName, req.BuildNum)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 删除Job构建记录
func IDeleteJenkinsBuildLog(c *gin.Context) {
	var req ci.JenkinsCancelBuildReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	buildInfo, err := jenkins_service.JenkinsInstance.DeleteJobBuild(req.JobName, req.BuildNum)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 获取构建日志
func IGetJenkinsBuildLog(c *gin.Context) {
	var req ci.JenkinsBuildInfoReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	buildInfo, err := jenkins_service.JenkinsInstance.GetBuildNumLog(req.JobName, req.Number)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}

// 创建系统凭证
func IPostJenkinsCredentials(c *gin.Context) {
	var req ci.JenkinsCreateCredentialsReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	info, err := jenkins_service.JenkinsInstance.CreateSystemCredentials(req.CredentialId,
		req.Username, req.Password, req.Description)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: info})
}

// 获取系统凭证信息
func IGetJenkinsCredentialsInfo(c *gin.Context) {
	var req ci.JenkinsGetCredentialsReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	info, err := jenkins_service.JenkinsInstance.GetSystemCredentialsInfo(req.CredentialId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: info})
}

// 获取系统凭证列表
func IGetJenkinsCredentialsList(c *gin.Context) {
	info, err := jenkins_service.JenkinsInstance.GetSystemCredentialList()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: info})
}

// 获取构建制品信息
func IGetJenkinsBuildArchiveArtifactsInfo(c *gin.Context) {
	var req ci.JenkinsBuildInfoReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	buildInfo, err := jenkins_service.JenkinsInstance.GetBuildArchiveArtifactsInfo(req.JobName, req.Number)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: buildInfo})
}