package nacos_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/chujieyang/commonops/ops/opslog"
	"github.com/chujieyang/commonops/ops/value_objects/nacos"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 7/20/21 10:53 AM
 * @Desc:
 */

type nacosClient struct {
	endPoint    string
	username    string
	password    string
	accessToken string
}

func (r *nacosClient) requestPost(url string, body string, headers map[string]string) (respData string, statusCode int, err error) {
	url = fmt.Sprintf("http://%s/%s?accessToken=%s", r.endPoint, url, r.accessToken)
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			opslog.Error().Println(err.Error())
		}
	}()
	statusCode = resp.StatusCode
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	respData = string(respBytes)
	return
}

func (r *nacosClient) requestDelete(url string, params string, body string, headers map[string]string) (respData string, statusCode int, err error) {
	url = fmt.Sprintf("http://%s/%s?accessToken=%s&%s", r.endPoint, url, r.accessToken, params)
	req, err := http.NewRequest("DELETE", url, strings.NewReader(body))
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			opslog.Error().Println(err.Error())
		}
	}()
	statusCode = resp.StatusCode
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	respData = string(respBytes)
	return
}

type configDetail struct {
	Id         string
	DataId     string
	Group      string
	Content    string
	Tenant     string
	AppName    string
	Type       string
	CreateTime int
	ModifyTime int
	CreateUser string
	CreateIp   string
	Desc       string
	ConfigTags string
}

func (r *nacosClient) requestGet(url string, params string, headers map[string]string) (respData string, statusCode int, err error) {
	url = fmt.Sprintf("http://%s/%s?accessToken=%s&%s", r.endPoint, url, r.accessToken, params)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			opslog.Error().Println(err.Error())
		}
	}()
	statusCode = resp.StatusCode
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	respData = string(respBytes)
	return
}

type LoginResp struct {
	AccessToken string `json:"accessToken"`
}

func NewNacosClient(endpoint, username, password string) (nacos *nacosClient, err error) {
	nacos = &nacosClient{
		endPoint:    endpoint,
		username:    username,
		password:    password,
		accessToken: "",
	}
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	respData, statusCode, err := nacos.requestPost("nacos/v1/auth/login", fmt.Sprintf("username=%s&password=%s", username, password), headers)
	if err != nil {
		return
	}
	if statusCode != 200 {
		err = errors.New("操作失败")
		return
	}
	var loginResult LoginResp
	if err = json.Unmarshal([]byte(respData), &loginResult); err != nil {
		return
	}
	nacos.accessToken = loginResult.AccessToken
	return
}

type namespaceItem struct {
	Namespace         string
	NamespaceShowName string
	Quota             int
	ConfigCount       int
	Type              int
}

type namespaceResp struct {
	Data []namespaceItem `json:"data"`
}

func (r *nacosClient) GetNamespace() (namespaceList namespaceResp, err error) {
	data, statusCode, err := r.requestGet("nacos/v1/console/namespaces", "", nil)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	if err = json.Unmarshal([]byte(data), &namespaceList); err != nil {
		return
	}
	return
}

func (r *nacosClient) GetConfig(namespace, dataId, group string) (data configDetail, err error) {
	respData, statusCode, err := r.requestGet("nacos/v1/cs/configs",
		fmt.Sprintf("show=all&tenant=%s&dataId=%s&group=%s", namespace, dataId, group), nil)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, respData))
		return
	}
	if err = json.Unmarshal([]byte(respData), &data); err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	return
}

func (r *nacosClient) PublishConfig(namespace, dataId, group, content, configType, config_tags string) (err error) {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	body := fmt.Sprintf("tenant=%s&dataId=%s&group=%s&content=%s&type=%s&config_tags=%s", namespace, dataId, group, content, configType, config_tags)
	data, statusCode, err := r.requestPost("nacos/v1/cs/configs", body, headers)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	return
}

func (r *nacosClient) CopyConfig(srcNamespace, srcDataId, srcGroup, dstNamespace, dstDataId, dstGroup string) (err error) {
	srcConfig, err := r.GetConfig(srcNamespace, srcDataId, srcGroup)
	if err != nil {
		return
	}
	if err = r.PublishConfig(dstNamespace, dstDataId, dstGroup, srcConfig.Content, srcConfig.Type, srcConfig.ConfigTags); err != nil {
		return
	}
	return
}

func (r *nacosClient) AppendStaticConfigToSelectAllConfigs(srcNamespace, srcDataId, srcGroup string, destList []nacos.SyncDstConfig) (err error) {
	srcConfig, err := r.GetConfig(srcNamespace, srcDataId, srcGroup)
	if err != nil {
		return
	}
	if srcConfig.ConfigTags != "static" {
		err = errors.New("不是静态配置，不允许该操作.")
		return
	}
	if srcConfig.Type != "yaml" && srcConfig.Type != "text" && srcConfig.Type != "properties" {
		err = errors.New("当前仅支持text、yaml、properties格式的静态类型的追加合并.")
		return
	}
	for _, destConfigMap := range destList {
		destConfig, err1 := r.GetConfig(destConfigMap.Namespace, destConfigMap.DataId, destConfigMap.Group)
		if err1 != nil {
			err = err1
			return
		}
		if destConfig.ConfigTags == "static" {
			continue
		}
		if destConfig.Type != "yaml" && destConfig.Type != "text" && destConfig.Type != "properties" {
			continue
		}
		newContent := fmt.Sprintf("%s\n%s", destConfig.Content, srcConfig.Content)
		if err = r.PublishConfig(destConfigMap.Namespace, destConfig.DataId, destConfig.Group, newContent, destConfig.Type, destConfig.ConfigTags); err != nil {
			return
		}
	}
	return
}

func (r *nacosClient) GetNsConfigs(namespace string, page, size int, config_tags string) (data string, err error) {
	data, statusCode, err := r.requestGet("nacos/v1/cs/configs",
		fmt.Sprintf("pageNo=%d&pageSize=%d&search=accurate&dataId=&group=&tenant=%s&config_tags=%s", page, size, namespace, config_tags), nil)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	return
}

func (r *nacosClient) DeleteConfig(namespace, dataId, group string) (err error) {
	params := fmt.Sprintf("tenant=%s&dataId=%s&group=%s", namespace, dataId, group)
	data, statusCode, err := r.requestDelete("nacos/v1/cs/configs", params, "", nil)
	if statusCode != 200 {
		err = errors.New(fmt.Sprintf("%s \n %s", err, data))
		return
	}
	return
}

type configItem struct {
	Id     string `json:"id"`
	DataId string `json:"dataId"`
	Group  string `json:"group"`
	Tenant string `json:"tenant"`
	Type   string `json:"type"`
}

type nsConfigs struct {
	Namespace string       `json:"namespace"`
	Configs   []configItem `json:"configs"`
}

type configsResp struct {
	TotalCount int
	PageItems  []configItem
}

func (r *nacosClient) GetAllConfigs() (data []nsConfigs, err error) {
	nsList, err := r.GetNamespace()
	if err != nil {
		return
	}
	for _, ns := range nsList.Data {
		nsConfigData := nsConfigs{
			Namespace: ns.Namespace,
		}
		page := 1
		pageSize := 10
		leftCount := 1
		for leftCount > 0 {
			var configsData configsResp
			configs, err1 := r.GetNsConfigs(ns.Namespace, page, pageSize, "")
			if err1 != nil {
				err = err1
				return
			}
			if err = json.Unmarshal([]byte(configs), &configsData); err != nil {
				return
			}
			leftCount = configsData.TotalCount - page*pageSize
			page += 1
			nsConfigData.Configs = append(nsConfigData.Configs, configsData.PageItems...)
		}
		data = append(data, nsConfigData)
	}
	return
}
