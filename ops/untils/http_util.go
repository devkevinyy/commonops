package untils

import (
	"io"
	"io/ioutil"
	"net/http"
)


func HttpGet(url string, authUsername string, authToken string) (respData string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	req.SetBasicAuth(authUsername, authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	respData = string(respBytes)
	return
}

/*
	用于访问jenkins获取实时的构建日志
 */
func HttpGetRespMap(url string, authUsername string, authToken string) (data map[string]interface{}, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	req.SetBasicAuth(authUsername, authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	data = map[string]interface{}{
		"content": string(respBytes),
		"start": resp.Header.Get("X-Text-Size"),
		"hasMore": resp.Header.Get("X-More-Data"),
	}
	return
}


func HttpPost(url string, body io.Reader, authUsername string, authToken string, contentType string) (respData string, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.SetBasicAuth(authUsername, authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Error(err.Error())
		return
	}
	respData = string(respBytes)
	return
}