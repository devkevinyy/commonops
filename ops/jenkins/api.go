package jenkins

import (
	"bytes"
	"fmt"
	"github.com/chujieyang/commonops/ops/conf"
	"github.com/chujieyang/commonops/ops/untils"
	"strings"
)

type Jenkins struct {
	Host string
	BasicUsername string
	BasicToken string
}

var JenkinsClient Jenkins

func init() {
	JenkinsClient = Jenkins{
		Host: conf.JenkinsHost,
		BasicUsername: conf.JenkinsUsername,
		BasicToken: conf.JenkinsToken,
	}
}

/*
	获取所有任务的信息
*/
func (this *Jenkins) GetAllJobList() (resp string, err error) {
	targetUrl := fmt.Sprintf("%s/api/json", this.Host)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	获取指定任务的信息
 */
func (this *Jenkins) GetJobInfo(jobName string) (resp string, err error) {
	targetUrl := fmt.Sprintf("%s/job/%s/api/json?pretty=true", this.Host, jobName)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	获取任务的描述信息
*/
func (this *Jenkins) GetJobDescription(jobName string) (resp string, err error) {
	targetUrl := fmt.Sprintf("%s/job/%s/description", this.Host, jobName)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	更新任务的描述信息
*/
func (this *Jenkins) UpdateJobDescription(jobName string, description string) (resp string, err error) {
	targetUrl := fmt.Sprintf("%s/job/%s/description", this.Host, jobName)
	body := strings.NewReader(fmt.Sprintf("description=%s", description))
	resp, err = untils.HttpPost(targetUrl, body, this.BasicUsername, this.BasicToken, "application/x-www-form-urlencoded")
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	获取任务某次的构建信息
 */
func (this *Jenkins) GetJobBuildInfo(jobName string, count int) (resp string, err error) {
	targetUrl := fmt.Sprintf("%s/job/%s/%d/api/json?pretty=true", this.Host, jobName, count)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	获取任务的配置文件
*/
func (this *Jenkins) GetJobConfigInfo(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/config.xml", this.Host, jobName)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}


/*
	删除任务
*/
func (this *Jenkins) DeleteJob(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/doDelete", this.Host, jobName)
	resp, err = untils.HttpPost(targetUrl, nil, this.BasicUsername, this.BasicToken, "application/json")
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	允许任务
*/
func (this *Jenkins) EnableJob(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/enable", this.Host, jobName)
	resp, err = untils.HttpPost(targetUrl, nil, this.BasicUsername, this.BasicToken, "application/json")
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	禁止任务
*/
func (this *Jenkins) DisableJob(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/disable", this.Host, jobName)
	resp, err = untils.HttpPost(targetUrl, nil, this.BasicUsername, this.BasicToken, "application/json")
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	构建任务
*/
func (this *Jenkins) BuildJob(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/build", this.Host, jobName)
	resp, err = untils.HttpPost(targetUrl, nil, this.BasicUsername, this.BasicToken, "application/json")
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	获取上次构建的序号
*/
func (this *Jenkins) GetJobLastBuildNumber(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/lastBuild/buildNumber", this.Host, jobName)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	获取上次构建的时间戳
*/
func (this *Jenkins) GetJobLastBuildTimestamp(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/lastBuild/buildTimestamp", this.Host, jobName)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	获取最近一次构建控制台输出
*/
func (this *Jenkins) GetJobLastBuildConsoleLog(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/lastBuild/logText/progressiveText", this.Host, jobName)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}


/*
	获取指定构建控制台输出
*/
func (this *Jenkins) GetJobBuildConsoleLog(jobName string, count int, start int) (data map[string]interface{}, err error)  {
	targetUrl := fmt.Sprintf("%s/job/%s/%d/logText/progressiveText?start=%d", this.Host, jobName, count, start)
	data, err = untils.HttpGetRespMap(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	创建任务
*/
func (this *Jenkins) CreateJob(jobName string, configXml string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/createItem?name=%s", this.Host, jobName)
	bodyContent := bytes.NewBuffer([]byte(configXml))
	resp, err = untils.HttpPost(targetUrl, bodyContent, this.BasicUsername, this.BasicToken, "text/xml")
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}

/*
	获取构建历史记录
*/
func (this *Jenkins) GetJobBuildList(jobName string) (resp string, err error)  {
	targetUrl := fmt.Sprintf("%s/queue/api/json", this.Host)
	resp, err = untils.HttpGet(targetUrl, this.BasicUsername, this.BasicToken)
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}