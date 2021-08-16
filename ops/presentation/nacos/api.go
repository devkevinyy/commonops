package nacos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"github.com/chujieyang/commonops/ops/infrastructure/database/models"
	"github.com/chujieyang/commonops/ops/services/nacos_service"
	"github.com/chujieyang/commonops/ops/utils"
	"github.com/chujieyang/commonops/ops/value_objects/nacos"
	"github.com/gin-gonic/gin"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 7/20/21 3:07 PM
 * @Desc:
 */

func IPostNacos(c *gin.Context) {
	var req nacos.NacosServer
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	_, err = models.AddNewNacosServer(req.Alias, req.EndPoint, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPostNacosConfig(c *gin.Context) {
	var req map[string]interface{}
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	templateId := 0
	configContent := ""
	fillData := ""
	if req["content"] != nil {
		configContent = req["content"].(string)
	}
	if req["isUseTemplate"].(string) == "useTemplate" {
		templateId = int(req["templateId"].(float64))
		templateDetail, err := models.GetConfigTemplateDetail(int(req["templateId"].(float64)))
		if err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		var fieldsList []string
		if err = json.Unmarshal([]byte(templateDetail.FillField), &fieldsList); err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		fillMap := map[string]interface{}{}
		for _, field := range fieldsList {
			fillMap[field] = req[field]
		}
		fillDataBytes, err := json.Marshal(fillMap)
		if err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		fillData = string(fillDataBytes)

		tmpl, err := template.New("parseConfigTemplate").Parse(templateDetail.ConfigContent)
		if err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		configContentBuf := new(bytes.Buffer)
		if err := tmpl.Execute(configContentBuf, fillMap); err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		configContent = configContentBuf.String()
	}

	nacosInfo, err := models.GetNacosInfoById(req["clusterId"].(string))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}

	tx := database.Mysql().Begin()
	models.AddNacosConfig(tx, req["clusterId"].(string), templateId,
		req["namespace"].(string), req["dataId"].(string), req["group"].(string), fillData)
	if err = nacosClient.PublishConfig(req["namespace"].(string), req["dataId"].(string), req["group"].(string),
		configContent, "yaml", ""); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetNacosConfigDetail(c *gin.Context) {
	var req nacos.GetNacosConfigDetailReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.ClusterId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	originData, err := models.GetNacosConfig(req.ClusterId, req.Namespace, req.DataId, req.Group)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	configDetail, err := nacosClient.GetConfig(req.Namespace, req.DataId, req.Group)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"originData": originData,
		"nacosData":  configDetail,
	}})
}

func IPostNacosConfigCopy(c *gin.Context) {
	var req nacos.CreateNacosConfigCopyReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.ClusterId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = nacosClient.CopyConfig(req.SrcNamespace, req.SrcDataId, req.SrcGroup, req.DstNamespace, req.DstDataId, req.DstGroup); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPostNacosConfigSync(c *gin.Context) {
	var req nacos.CreateNacosConfigSyncReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.ClusterId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = nacosClient.AppendStaticConfigToSelectAllConfigs(req.SrcNamespace, req.SrcDataId, req.SrcGroup, req.DstConfigs); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPutNacosConfig(c *gin.Context) {
	var req map[string]interface{}
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	templateId := 0
	configContent := ""
	fillData := ""
	if req["content"] != nil {
		configContent = req["content"].(string)
	}
	if req["isUseTemplate"].(string) == "useTemplate" {
		templateId = int(req["templateId"].(float64))
		templateDetail, err := models.GetConfigTemplateDetail(int(req["templateId"].(float64)))
		if err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		var fieldsList []string
		if err = json.Unmarshal([]byte(templateDetail.FillField), &fieldsList); err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		fillMap := map[string]interface{}{}
		for _, field := range fieldsList {
			fillMap[field] = req[field]
		}
		fillDataBytes, err := json.Marshal(fillMap)
		if err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		fillData = string(fillDataBytes)

		tmpl, err := template.New("parseConfigTemplate").Parse(templateDetail.ConfigContent)
		if err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		configContentBuf := new(bytes.Buffer)
		if err := tmpl.Execute(configContentBuf, fillMap); err != nil {
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		configContent = configContentBuf.String()
	}

	nacosInfo, err := models.GetNacosInfoById(req["clusterId"].(string))
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	tx := database.Mysql().Begin()
	models.UpdateNacosConfig(tx, int(req["id"].(float64)), templateId, req["dataId"].(string), req["group"].(string), fillData)
	if err = nacosClient.PublishConfig(req["namespace"].(string), req["dataId"].(string), req["group"].(string),
		configContent, "yaml", ""); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IDeleteNacosConfig(c *gin.Context) {
	var req nacos.DeleteNacosConfigReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.ClusterId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	tx := database.Mysql().Begin()
	err = models.DeleteNacosConfig(tx, req.Id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = nacosClient.DeleteConfig(req.Namespace, req.DataId, req.Group); err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IGetNacosList(c *gin.Context) {
	nacosList, err := models.GetNacosList()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nacosList})
}

func IGetNacosNamespaceList(c *gin.Context) {
	var req nacos.NacosNsListReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.ClusterId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nsList, err := nacosClient.GetNamespace()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nsList})
}

func IGetNacosConfigList(c *gin.Context) {
	var req nacos.NacosConfigsReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	// nacosInfo, err := models.GetNacosInfoById(req.ClusterId)
	// if err != nil {
	// 	c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
	// 	return
	// }
	// nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	// if err != nil {
	// 	c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
	// 	return
	// }
	// configsList, err := nacosClient.GetNsConfigs(req.Namespace, req.Page, req.Size, req.ConfigTags)
	// if err != nil {
	// 	c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
	// 	return
	// }
	offset := (req.Page - 1) * req.Size
	count, err := models.GetOpsNacosConfigCount(req.ClusterId, req.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	data, err := models.GetOpsNacosConfigByPage(req.ClusterId, req.Namespace, offset, req.Size)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"count":       count,
		"currentPage": req.Page,
		"data":        data,
	}})
}

func IGetNacosAllConfigs(c *gin.Context) {
	var req nacos.NacosNsListReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosInfo, err := models.GetNacosInfoById(req.ClusterId)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	allConfigs, err := nacosClient.GetAllConfigs()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: allConfigs})
}

func IGetConfigTemplates(c *gin.Context) {
	var req nacos.ConfigTemplatesReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	count, err := models.GetConfigTemplateCount(req.Name)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	offset := (req.Page - 1) * req.Size
	data, err := models.GetConfigTemplateByPage(req.Name, offset, req.Size)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: map[string]interface{}{
		"data":  data,
		"count": count,
	}})
}

func IGetConfigTemplatesAll(c *gin.Context) {
	data, err := models.GetConfigTemplatesAll()
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: data})
}

func IPostConfigTemplate(c *gin.Context) {
	var req nacos.ConfigTemplateAddReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = models.AddConfigTemplate(req.Name, req.ConfigContent, req.FillField); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}

func IPutConfigTemplate(c *gin.Context) {
	var req nacos.ConfigTemplateUpdateReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}

	fillFieldBytes, err := json.Marshal(req.FillField)
	if err != nil {
		return
	}
	templateInfo, err := models.GetConfigTemplateById(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}

	tx := database.Mysql().Begin()
	warnMsg := ""
	if templateInfo.FillField != string(fillFieldBytes) {
		warnMsg = "由于你修改了动态填充项，当前保存只保存了模板数据，nacos具体配置的生效请单独对配置进行修改!"
		if err = models.UpdateConfigTemplate(tx, req.Id, req.Name, req.ConfigContent, string(fillFieldBytes)); err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
	} else {
		if err = models.UpdateConfigTemplate(tx, req.Id, req.Name, req.ConfigContent, string(fillFieldBytes)); err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		configsList, err := models.GetAllNacosConfigRelatedToTheTemplate(tx, req.Id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
			return
		}
		for _, config := range configsList {
			clusterId, namespace, dataId, group, templateContent, fillData := config.ClusterId, config.Namespace, config.DataId, config.ConfigGroup, config.ConfigContent, config.FillData
			fmt.Println(clusterId, namespace, dataId, group)
			nacosInfo, err := models.GetNacosInfoById(clusterId)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
				return
			}
			nacosClient, err := nacos_service.NewNacosClient(nacosInfo.EndPoint, nacosInfo.Username, nacosInfo.Password)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
				return
			}
			tmpl, err := template.New("parseConfigTemplate").Parse(templateContent)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
				return
			}
			var fillMap map[string]interface{}
			if err = json.Unmarshal([]byte(fillData), &fillMap); err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
				return
			}
			configContentBuf := new(bytes.Buffer)
			if err := tmpl.Execute(configContentBuf, fillMap); err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
				return
			}
			if err = nacosClient.PublishConfig(namespace, dataId, group, configContentBuf.String(), "yaml", ""); err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
				return
			}
		}
		warnMsg = "由于你未修改动态填充项，已对nacos配置进行同步更新!"
	}
	tx.Commit()
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: fmt.Sprintf("修改成功: %s", warnMsg), Data: nil})
}

func IDeleteConfigTemplate(c *gin.Context) {
	var req nacos.ConfigTemplateDeleteReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	if err = models.DeleteConfigTemplate(req.Id); err != nil {
		c.JSON(http.StatusOK, utils.RespData{Code: -1, Msg: err.Error(), Data: nil})
		return
	}
	c.JSON(http.StatusOK, utils.RespData{Code: 0, Msg: "success", Data: nil})
}
