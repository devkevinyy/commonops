package prometheus

import (
	"fmt"

	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/untils"
)

func PrometheusQuery(query string) (resp string, err error) {
	config, err := models.GetPrometheusValue()
	if err != nil {
		return
	}
	targetUrl := fmt.Sprintf("%s/api/v1/query?query=%s", config.Value, query)
	resp, err = untils.HttpGet(targetUrl, "", "")
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}
