package jenkins

import (
	"encoding/json"
	"github.com/chujieyang/commonops/ops/forms"
	"github.com/chujieyang/commonops/ops/jenkins"
	"github.com/chujieyang/commonops/ops/untils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IGetJenkinsAllJobs(c *gin.Context){
	resp, err := jenkins.JenkinsClient.GetAllJobList()
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}
	var mapData map[string]interface{}
	err = json.Unmarshal([]byte(resp), &mapData)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg: "success", Data: mapData["jobs"]})
}

func IGetJenkinsJobBuildList(c *gin.Context){
	var req forms.JobBuildForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	resp, err := jenkins.JenkinsClient.GetJobInfo(req.Name)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}

	var mapData map[string]interface{}
	err = json.Unmarshal([]byte(resp), &mapData)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}

	var buildLog []interface{}
	count := 1
	for i:=int(mapData["nextBuildNumber"].(float64))-1; i>0; i-- {
		if count > 5 {
			break
		}
		buildInfo, _ := jenkins.JenkinsClient.GetJobBuildInfo(req.Name, i)
		buildLog = append(buildLog, buildInfo)
		count++
	}

	respData := map[string]interface{}{
		"displayName": mapData["displayName"],
		"description": mapData["description"],
		"buildable": mapData["buildable"],
		"inQueue": mapData["inQueue"],
		"concurrentBuild": mapData["concurrentBuild"],
		"nextBuildNumber": mapData["nextBuildNumber"],
		"buildLog": buildLog,
	}
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg: "success", Data: respData})
}

func IGetJenkinsJobBuildLog(c *gin.Context){
	var req forms.JobBuildForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	resp, err := jenkins.JenkinsClient.GetJobBuildConsoleLog(req.Name, req.BuildId, req.Start)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg: "success", Data: resp})
}

func IGetJenkinsJobLastBuildLog(c *gin.Context){
	var req forms.JobBuildForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	resp, err := jenkins.JenkinsClient.GetJobLastBuildConsoleLog(req.Name)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg: "success", Data: resp})
}

func IPostJenkinsStartJob(c *gin.Context){
	var req forms.JobBuildForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	resp, err := jenkins.JenkinsClient.BuildJob(req.Name)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg: "success", Data: resp})
}

func IPostJenkinsEnableJob(c *gin.Context){
	var req forms.JobBuildForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	resp, err := jenkins.JenkinsClient.EnableJob(req.Name)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg: "success", Data: resp})
}

func IPostJenkinsDisableJob(c *gin.Context){
	var req forms.JobBuildForm
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code: -1, Msg: string(err.Error()), Data: nil})
		return
	}
	resp, err := jenkins.JenkinsClient.DisableJob(req.Name)
	if err != nil {
		c.JSON(http.StatusOK, untils.RespData{Code:-1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, untils.RespData{Code:0, Msg: "success", Data: resp})
}