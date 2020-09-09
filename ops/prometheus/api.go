package prometheus

import (
	"fmt"

	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/untils"
)

func PrometheusQuery(query string, start int, end int, step string) (resp string, err error) {
	config, err := models.GetPrometheusValue()
	if err != nil {
		return
	}
	targetUrl := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%d&end=%d&step=%s",
		config.Value, query, start, end, step)
	resp, err = untils.HttpGet(targetUrl, "", "")
	if err != nil {
		untils.Log.Error(err.Error())
		return
	}
	return
}
