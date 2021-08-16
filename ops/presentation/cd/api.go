package cd

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/domain/service"
	"github.com/chujieyang/commonops/ops/value_objects/cd"
	"net/http"
	"strings"

	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/presentation/k8s"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/gin-gonic/gin"
)

// 获取发布模板
func IGetCdProcessTemplateList(c *gin.Context) {
	empId, _ := c.Get("empId")
	templateList, err := models.GetAllCdProcessTemplateName(empId.(string))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: templateList})
}

// 创建发布模板
func IPostCdProcessTemplate(c *gin.Context) {
	var req cd.CdProcessTemplateReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	empId, _ := c.Get("empId")
	templateId, err := models.AddNewCdProcessTemplate(empId.(string), req.JobName, req.TemplateName, req.ClusterId, req.Namespace,
		req.DeployYaml, req.ServiceYaml, req.ConfigmapYaml, req.IngressYaml)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	cdProcessId, err := models.AddCdProcessLog(empId.(string), templateId, req.ImageName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	// 执行项目发布过程
	apiServer, token := models.GetK8sInfoByClusterId(req.ClusterId)
	k8sClient, err := service.NewKubernetesService(apiServer, token)
	if err != nil {
		return
	}
	var errLog []string
	_, err = k8sClient.ApplyYaml(req.Namespace, req.ConfigmapYaml)
	if err != nil {
		errLog = append(errLog, err.Error())
		models.UpdateCdProcessLogResult(cdProcessId, -1, strings.Join(errLog, "; "))
		return
	}
	_, err = k8sClient.ApplyYaml(req.Namespace, req.DeployYaml)
	if err != nil {
		errLog = append(errLog, err.Error())
		models.UpdateCdProcessLogResult(cdProcessId, -1, strings.Join(errLog, "; "))
		return
	}
	_, err = k8sClient.ApplyYaml(req.Namespace, req.ServiceYaml)
	if err != nil {
		errLog = append(errLog, err.Error())
		models.UpdateCdProcessLogResult(cdProcessId, -1, strings.Join(errLog, "; "))
		return
	}
	_, err = k8sClient.ApplyYaml(req.Namespace, req.IngressYaml)
	if err != nil {
		errLog = append(errLog, err.Error())
		models.UpdateCdProcessLogResult(cdProcessId, -1, strings.Join(errLog, "; "))
		return
	}
	models.UpdateCdProcessLogResult(cdProcessId, 1, "发布成功")
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: cdProcessId})
}

// 根据模板发布
func IPostCdProcessLog(c *gin.Context) {
	var req cd.CdProcessTemplateLog
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	empId, _ := c.Get("empId")
	cdProcessId, err := models.AddCdProcessLog(empId.(string), req.TemplateId, req.ImageName)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	templateInfo, err := models.GetTemplateInfoById(req.TemplateId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	err = k8s.UpdateDeploymentWithImage(templateInfo.ClusterId, templateInfo.Namespace, fmt.Sprintf("%s-deploy", templateInfo.JobName), req.ImageName)
	if err != nil {
		models.UpdateCdProcessLogResult(cdProcessId, -1, err.Error())
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	} else {
		models.UpdateCdProcessLogResult(cdProcessId, 1, "发布成功")
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetCdProcessLog(c *gin.Context) {
	var req cd.GetCdProcessLogReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	empId, _ := c.Get("empId")
	offset := (req.Page - 1) * req.Size
	totalCount, err := models.GetCdProcessLogCount(empId.(string))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	data, err := models.GetCdProcessLogByPage(empId.(string), offset, req.Size)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	resp := cd.CdProcessLogResp{
		Total: totalCount,
		Page:  req.Page,
		Logs:  data,
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: resp})
}
