package models

import (
	"fmt"
	"github.com/chujieyang/commonops/ops/infrastructure/database"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"strconv"
	"strings"
	"time"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 4:09 PM
 * @Desc:
 */

type ResourceKind string

type DataPoints []DataPoint

type DataPoint struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type MetricPoint struct {
	Timestamp string `json:"timestamp"`
	Value     uint64 `json:"value"`
}

type ResourceSelector struct {
	Namespace    string
	ResourceType ResourceKind
	ResourceName string
	Selector     map[string]string
	UID          types.UID
}

type SidecarMetric struct {
	DataPoints   `json:"dataPoints"`
	MetricPoints []MetricPoint `json:"metricPoints"`
	MetricName   string        `json:"metricName"`
	UIDs         []types.UID   `json:"uids"`
}

func (metric *SidecarMetric) AddMetricPoint(item MetricPoint) []MetricPoint {
	metric.MetricPoints = append(metric.MetricPoints, item)
	return metric.MetricPoints
}

type SidecarMetricResultList struct {
	Items []SidecarMetric `json:"items"`
}

type NodeMetric struct {
	Uid     string `json:"uid"`
	Cid     string `json:"cid"`
	Name    string `json:"name"`
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
	Time    string `json:"time"`
}

type MetricItem struct {
	Uid       string `json:"uid"`
	Cid       string `json:"cid"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Container string `json:"container"`
	Cpu       string `json:"cpu"`
	Memory    string `json:"memory"`
	Storage   string `json:"storage"`
	Time      string `json:"time"`
}

func UpdateNodeMetrics(clusterId string, nodeMetrics *v1beta1.NodeMetricsList) (err error) {
	for _, v := range nodeMetrics.Items {
		insertNodeMetricsSql := fmt.Sprintf("insert into nodes(uid, cid, name, cpu, memory, storage, time) values(?, ?, ?, ?, ?, ?, now())")
		err = database.Mysql().Exec(insertNodeMetricsSql, v.UID, clusterId, v.Name, v.Usage.Cpu().MilliValue(), v.Usage.Memory().Value()/1024.0/1024.0, v.Usage.StorageEphemeral().Value()).Error
		if err != nil {
			return
		}
	}
	return
}

func UpdatePodMetrics(clusterId string, podMetrics *v1beta1.PodMetricsList) (err error) {
	for _, v := range podMetrics.Items {
		for _, u := range v.Containers {
			insertPodMetricsSql := fmt.Sprintf("insert into pods(uid, cid, name, namespace, container, cpu, memory, storage, time) values(?, ?, ?, ?, ?, ?, ?, ?, now())")
			err = database.Mysql().Exec(insertPodMetricsSql, v.UID, clusterId, v.Name, v.Namespace, u.Name, u.Usage.Cpu().MilliValue(), u.Usage.Memory().Value()/1024.0/1024.0, u.Usage.StorageEphemeral().Value()).Error
			if err != nil {
				return
			}
		}
	}
	return
}

/*
	CleanDatabase deletes rows from nodes and pods based on a time window.
*/
func CleanDatabase(window *time.Duration) (err error) {
	cleanNodeMetricsSql := "delete from nodes where time <= date_sub(now(), interval 300 second);"
	err = database.Mysql().Exec(cleanNodeMetricsSql).Error
	if err != nil {
		return
	}

	cleanPodMetricsSql := "delete from pods where time <= date_sub(now(), interval 300 second);"
	err = database.Mysql().Exec(cleanPodMetricsSql).Error
	if err != nil {
		return
	}
	return
}

func getRows(clusterId string, table string, metricName string, selector ResourceSelector) (rows []MetricItem, err error) {
	var query string
	var values []interface{}
	var args []string
	if metricName == "cpu" {
		query = "select cpu, name, uid, time from %s "
	} else {
		query = "select memory, name, uid, time from %s "
	}

	if table == "pods" {
		args = append(args, "namespace=?")
		if selector.Namespace != "" {
			values = append(values, selector.Namespace)
		} else {
			values = append(values, "default")
		}
	}

	if selector.ResourceName != "" {
		if strings.ContainsAny(selector.ResourceName, ",") {
			subargs := []string{}
			for _, v := range strings.Split(selector.ResourceName, ",") {
				subargs = append(subargs, "?")
				values = append(values, v)
			}
			args = append(args, " name in ("+strings.Join(subargs, ",")+")")
		} else {
			values = append(values, selector.ResourceName)
			args = append(args, " name = ?")
		}
	}
	if selector.UID != "" {
		args = append(args, " uid = ?")
		values = append(values, selector.UID)
	}
	args = append(args, " cid = ?")
	values = append(values, clusterId)

	query = fmt.Sprintf(query+" where "+strings.Join(args, " and ")+" order by time asc;", table)

	err = database.Mysql().Raw(query, values...).Scan(&rows).Error
	return
}

/*
	getPodMetrics: With a database connection and a resource selector
	Queries SQLite and returns a list of metrics.
*/
func GetPodMetrics(clusterId string, metricName string, selector ResourceSelector) (SidecarMetricResultList, error) {
	rows, err := getRows(clusterId, "pods", metricName, selector)
	if err != nil {
		fmt.Printf("Error getting pod metrics: %s \n", err)
		return SidecarMetricResultList{}, err
	}

	var result SidecarMetricResultList
	var metrics []MetricPoint
	for _, row := range rows {
		layout := "2006-01-02T15:04:05+08:00"
		local, _ := time.LoadLocation("Local")
		t, err := time.ParseInLocation(layout, row.Time, local)
		if err != nil {
			return SidecarMetricResultList{}, err
		}
		var metricValue uint64
		if metricName == "cpu" {
			metricValue, err = strconv.ParseUint(row.Cpu, 10, 64)
			if err != nil {
				return SidecarMetricResultList{}, err
			}
		} else {
			metricValue, err = strconv.ParseUint(row.Memory, 10, 64)
			if err != nil {
				return SidecarMetricResultList{}, err
			}
		}

		newMetric := MetricPoint{
			Timestamp: t.Format("15:04:05"),
			Value:     metricValue,
		}
		metrics = append(metrics, newMetric)
	}
	item := SidecarMetric{
		MetricName:   metricName,
		MetricPoints: metrics,
		DataPoints:   []DataPoint{},
		UIDs: []types.UID{
			types.UID(selector.UID),
		},
	}
	result.Items = append(result.Items, item)
	return result, nil
}

/*
	getNodeMetrics: With a database connection and a resource selector
	Queries SQLite and returns a list of metrics.
*/
func GetNodeMetrics(clusterId string, metricName string, selector ResourceSelector) (SidecarMetricResultList, error) {
	rows, err := getRows(clusterId, "nodes", metricName, selector)
	if err != nil {
		fmt.Printf("Error getting pod metrics: %s \n", err)
		return SidecarMetricResultList{}, err
	}
	var result SidecarMetricResultList
	var metrics []MetricPoint
	for _, row := range rows {
		layout := "2006-01-02T15:04:05+08:00"
		local, _ := time.LoadLocation("Local")
		t, err := time.ParseInLocation(layout, row.Time, local)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		var metricValue uint64
		if metricName == "cpu" {
			metricValue, err = strconv.ParseUint(row.Cpu, 10, 64)
			if err != nil {
				return SidecarMetricResultList{}, err
			}
		} else {
			metricValue, err = strconv.ParseUint(row.Memory, 10, 64)
			if err != nil {
				return SidecarMetricResultList{}, err
			}
		}

		newMetric := MetricPoint{
			Timestamp: t.Format("15:04:05"),
			Value:     metricValue,
		}
		metrics = append(metrics, newMetric)
	}
	item := SidecarMetric{
		MetricName:   metricName,
		MetricPoints: metrics,
		DataPoints:   []DataPoint{},
		UIDs: []types.UID{
			types.UID(selector.UID),
		},
	}
	result.Items = append(result.Items, item)
	return result, nil
}
