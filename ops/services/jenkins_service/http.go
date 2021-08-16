package jenkins_service

import (
	"github.com/chujieyang/commonops/ops/opslog"
	"io"
	"io/ioutil"
	"net/http"
)

func httpGet(url string, authUsername string, authToken string) (respData string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	req.SetBasicAuth(authUsername, authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	respData = string(respBytes)
	if resp.StatusCode == 404 {
		respData = "暂无资源"	
	}
	return
}

/*
	用于访问jenkins获取实时的构建日志
*/
func httpGetRespMap(url string, authUsername string, authToken string) (data map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	req.SetBasicAuth(authUsername, authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	data = map[string]interface{}{
		"content": string(respBytes),
		"start":   resp.Header.Get("X-Text-Size"),
		"hasMore": resp.Header.Get("X-More-Data"),
	}
	return
}

func httpPost(url string, body io.Reader, authUsername string, authToken string, headers map[string]string) (respData string, statusCode int, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.SetBasicAuth(authUsername, authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		opslog.Error().Println(err.Error())
		return
	}
	respData = string(respBytes)
	return
}
