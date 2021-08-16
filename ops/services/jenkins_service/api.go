package jenkins_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chujieyang/commonops/ops/opslog"
	"net/url"
	"strconv"
	"strings"
	"text/template"

	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
)

var JenkinsInstance jenkinsClient

type jenkinsClient struct {
	authHost  string
	authUser  string
	authToken string
}

func init() {
	JenkinsInstance = jenkinsClient{
		authHost:  conf.JenkinsHost,
		authUser:  conf.JenkinsAuthUser,
		authToken: conf.JenkinsAuthToken,
	}
}

func (c jenkinsClient) CreateSystemCredentials(credentialsId string, username string, passwd string, description string) (jobInfo JobInfo, err error) {
	crumb, err := c.getSystemCrumb()
	if err != nil {
		return
	}
	targetUrl := fmt.Sprintf("%s/credentials/store/system/domain/_/createCredentials", c.authHost)
	formData := url.Values{}
	credentialsData := map[string]string{
		"scope":         "GLOBAL",
		"id":            credentialsId,
		"username":      username,
		"password":      passwd,
		"description":   description,
		"stapler-class": "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl",
	}
	formDataMap := map[string]interface{}{
		"":            "0",
		"credentials": credentialsData,
	}
	credentialsDataBytes, err := json.Marshal(formDataMap)
	if err != nil {
		return
	}
	formData.Add("json", string(credentialsDataBytes))
	headers := map[string]string{
		"Content-Type":   "application/x-www-form-urlencoded",
		"Jenkins-Crumb":  crumb.Crumb,
		"Content-Length": strconv.Itoa(len(formData.Encode())),
	}
	_, statusCode, err := httpPost(targetUrl, strings.NewReader(formData.Encode()), c.authUser, c.authToken, headers)
	if err != nil {
		return
	}
	if statusCode != 200 {
		err = errors.New("操作失败")
		return
	}
	return
}

func (c jenkinsClient) getSystemCrumb() (crumb CrumbItem, err error) {
	url := fmt.Sprintf("%s/crumbIssuer/api/json", c.authHost)
	respStr, err := httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(respStr), &crumb); err != nil {
		return
	}
	return
}

// 获取指定的系统凭证信息
func (c jenkinsClient) GetSystemCredentialsInfo(credentialsId string) (credentialInfo CredentailsInfo, err error) {
	url := fmt.Sprintf("%s/credentials/store/system/domain/_/credential/%s/api/json", c.authHost, credentialsId)
	respStr, err := httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(respStr), &credentialInfo); err != nil {
		return
	}
	return
}

// 获取系统凭证信息列表
func (c jenkinsClient) GetSystemCredentialList() (credentialsList CredentialsList, err error) {
	url := fmt.Sprintf("%s/credentials/store/system/domain/_/api/json?depth=1", c.authHost)
	respStr, err := httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(respStr), &credentialsList); err != nil {
		return
	}
	return
}

// 创建新的构建任务
func (c jenkinsClient) CreateJobItem(jobName string, configData models.CiConfigInfo) (jobInfo JobInfo, err error) {
	crumb, err := c.getSystemCrumb()
	if err != nil {
		return
	}
	targetUrl := fmt.Sprintf("%s/createItem?name=%s", c.authHost, jobName)
	configXml, err := generateTemplate("./services/jenkins_service/config.template.xml", configData)
	if err != nil {
		return
	}
	headers := map[string]string{
		"Content-Type":   "application/xml;charset=UTF-8",
		"Jenkins-Crumb":  crumb.Crumb,
		"Content-Length": strconv.Itoa(len(configXml)),
	}
	resp, statusCode, err := httpPost(targetUrl, strings.NewReader(configXml), c.authUser, c.authToken, headers)
	if err != nil {
		return
	}
	if statusCode != 200 {
		err = errors.New("操作失败: " + resp)
		return
	}
	return
}

// 更新构建任务的构建文件
func (c jenkinsClient) UpdateJobItem(jobName string, configScript models.CiScriptInfo) (jobInfo JobInfo, err error) {
	crumb, err := c.getSystemCrumb()
	if err != nil {
		return
	}
	targetUrl := fmt.Sprintf("%s/job/%s/config.xml", c.authHost, jobName)

	configXml, err := generateTemplate("/Users/yangchujie/GoProjects/src/github.com/chujieyang/commonops/ops/services/jenkins_service/script.template.xml", configScript)
	if err != nil {
		return
	}
	headers := map[string]string{
		"Content-Type":   "application/xml;charset=UTF-8",
		"Jenkins-Crumb":  crumb.Crumb,
		"Content-Length": strconv.Itoa(len(configXml)),
	}
	resp, statusCode, err := httpPost(targetUrl, strings.NewReader(configXml), c.authUser, c.authToken, headers)
	if err != nil {
		return
	}
	if statusCode != 200 {
		err = errors.New("操作失败: " + resp)
		return
	}
	return
}

// 删除构建任务
func (c jenkinsClient) DeleteJobItem(jobName string) (credentialsList CredentialsList, err error) {
	url := fmt.Sprintf("%s/job/%s/doDelete", c.authHost, jobName)
	resp, statusCode, err := httpPost(url, nil, c.authUser, c.authToken, nil)
	if err != nil {
		return
	}
	if statusCode != 200 {
		err = errors.New("操作失败: " + resp)
		return
	}
	return
}

// 获取项目列表
func (c jenkinsClient) GetJobList() (jenkinsJobInfo JenkinsJobInfo, err error) {
	targetUrl := fmt.Sprintf("%s/api/json", c.authHost)
	respStr, err := httpGet(targetUrl, c.authUser, c.authToken)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(respStr), &jenkinsJobInfo); err != nil {
		return
	}
	return
}

// 获取项目的构建文件信息
func (c jenkinsClient) GetJobItemConfig(jobName string) (configXml string, err error) {
	targetUrl := fmt.Sprintf("%s/job/%s/config.xml", c.authHost, jobName)
	configXml, err = httpGet(targetUrl, c.authUser, c.authToken)
	if err != nil {
		return
	}
	return
}

func (c jenkinsClient) GetJobBuildList(jobName string) (jobInfo JobInfo, err error) {
	url := fmt.Sprintf("%s/job/%s/api/json", c.authHost, jobName)
	respStr, err := httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(respStr), &jobInfo); err != nil {
		return
	}
	var builds []BuildInfo
	for _, buildItem := range jobInfo.Builds {
		buildInfo, err1 := c.GetBuildInfo(jobName, buildItem.Number)
		if err1 != nil {
			err = err1
			return
		}
		builds = append(builds, buildInfo)
	}
	jobInfo.BuildsDetails = builds
	return
}

func (c jenkinsClient) GetBuildInfo(jobName string, number int32) (buildInfo BuildInfo, err error) {
	url := fmt.Sprintf("%s/job/%s/%d/api/json", c.authHost, jobName, number)
	respStr, err := httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(respStr), &buildInfo); err != nil {
		return
	}
	return
}

func (c jenkinsClient) PostJobBuild(jobName string) (buildInfo BuildInfo, err error) {
	url := fmt.Sprintf("%s/job/%s/build/api/json", c.authHost, jobName)
	body := bytes.NewBuffer([]byte(""))
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	_, statusCode, err := httpPost(url, body, c.authUser, c.authToken, headers)
	if err != nil {
		return
	}
	if statusCode != 201 {
		err = errors.New("操作失败")
		return
	}
	return
}

func (c jenkinsClient) StopJobBuild(jobName string, buildNumber int32) (jobInfo JobInfo, err error) {
	url := fmt.Sprintf("%s/job/%s/%d/api/json", c.authHost, jobName, buildNumber)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	_, statusCode, err := httpPost(url, nil, c.authUser, c.authToken, headers)
	if err != nil {
		return
	}
	if statusCode != 201 {
		err = errors.New("操作失败")
		return
	}
	return
}

func (c jenkinsClient) DeleteJobBuild(jobName string, buildNumber int32) (buildInfo BuildInfo, err error) {
	url := fmt.Sprintf("%s/job/%s/%d/doDelete", c.authHost, jobName, buildNumber)
	body := bytes.NewBuffer([]byte(""))
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	_, statusCode, err := httpPost(url, body, c.authUser, c.authToken, headers)
	if err != nil {
		return
	}
	if statusCode != 200 {
		err = errors.New("操作失败")
		return
	}
	return
}

func (c jenkinsClient) GetBuildNumLog(jobName string, buildNumber int32) (buildLog string, err error) {
	url := fmt.Sprintf("%s/job/%s/%d/logText/progressiveText/api/json", c.authHost, jobName, buildNumber)
	buildLog, err = httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	return
}

func (c jenkinsClient) GetBuildNumStageList(jobName string, buildNumber int32) (buildLog string, err error) {
	url := fmt.Sprintf("%s/job/%s/%d/wfapi/", c.authHost, jobName, buildNumber)
	buildLog, err = httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	return
}

func (c jenkinsClient) GetBuildNumStageDetailLog(jobName string, buildNumber int32, nodeId int32) (buildLog BuildStageDescription, err error) {
	url := fmt.Sprintf("%s/job/%s/%d/execution/node/%d/wfapi/describe", c.authHost, jobName, buildNumber, nodeId)
	respStr, err := httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(respStr), &buildLog); err != nil {
		return
	}
	var newStateNodes []StageFlowNodes
	for _, stageNode := range buildLog.StageFlowNodes {
		nodeId, _ := strconv.Atoi(stageNode.Id)
		logInfo, err1 := c.GetBuildNumStageFlowNodeLog(jobName, buildNumber, int32(nodeId))
		if err1 != nil {
			err = err1
			return
		}
		stageNode.Log = logInfo.Text
		if stageNode.Log != "" {
			newStateNodes = append(newStateNodes, stageNode)
		}
	}
	buildLog.StageFlowNodes = newStateNodes
	return
}

func (c jenkinsClient) GetBuildNumStageFlowNodeLog(jobName string, buildNumber int32, nodeId int32) (log StageFlowNodeLog, err error) {
	url := fmt.Sprintf("%s/job/%s/%d/execution/node/%d/wfapi/log/", c.authHost, jobName, buildNumber, nodeId)
	respStr, err := httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(respStr), &log); err != nil {
		return
	}
	return
}

func (c jenkinsClient) GetBuildArchiveArtifactsInfo(jobName string, buildNumber int32) (result string, err error) {
	url := fmt.Sprintf("%s/job/%s/%d/artifact/jenkins_ci_image_result.info/*view*/", c.authHost, jobName, buildNumber)
	result, err = httpGet(url, c.authUser, c.authToken)
	if err != nil {
		return
	}
	return
}

func generateTemplate(templateFilePath string, data interface{}) (fileContent string, err error) {
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		opslog.Info().Println("create template failed, err:", err)
		return
	}
	var tmplBytes bytes.Buffer
	err = tmpl.Execute(&tmplBytes, data)
	if err != nil {
		return
	}
	fileContent = tmplBytes.String()
	return
}
