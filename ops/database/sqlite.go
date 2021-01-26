package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

var metricsDB *sql.DB

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

// ResourceSelector is a structure used to quickly and uniquely identify given resource.
// This struct can be later used for heapster data download etc.
type ResourceSelector struct {
	// Namespace of this resource.
	Namespace string
	// Type of this resource
	ResourceType ResourceKind
	// Name of this resource.
	ResourceName string
	// Selector used to identify this resource (should be used only for Deployments!).
	Selector map[string]string
	// UID is resource unique identifier.
	UID types.UID
}

// SidecarMetric is a format of data used by our sidecar. This is also the format of data that is being sent by backend API.
type SidecarMetric struct {
	// DataPoints is a list of X, Y int64 data points, sorted by X.
	DataPoints `json:"dataPoints"`
	// MetricPoints is a list of value, timestamp metrics used for sparklines on a pod list page.
	MetricPoints []MetricPoint `json:"metricPoints"`
	// MetricName is the name of metric stored in this struct.
	MetricName string `json:"metricName"`
	// Label stores information about identity of resources (UIDS) described by this metric.
	UIDs []types.UID `json:"uids"`
}

func (metric *SidecarMetric) AddMetricPoint(item MetricPoint) []MetricPoint {
	metric.MetricPoints = append(metric.MetricPoints, item)
	return metric.MetricPoints
}

type SidecarMetricResultList struct {
	Items []SidecarMetric `json:"items"`
}

func init() {
	// os.Remove("./k8s_metrics.db")
	metricsDB, err := sql.Open("sqlite3", "./k8s_metrics.db")
	if err != nil {
		fmt.Println("init metrics sqlite db error: ", err)
	}
	CreateDatabase(metricsDB)
	defer metricsDB.Close()
}

/*
	CreateDatabase creates tables for node and pod metrics
*/
func CreateDatabase(db *sql.DB) error {
	sqlStmt := `
	create table if not exists nodes (uid text, cid text, name text, cpu text, memory text, storage text, time datetime);
	create table if not exists pods (uid text, cid text, name text, namespace text, container text, cpu text, memory text, storage text, time datetime);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

/*
	UpdateDatabase updates nodeMetrics and podMetrics with scraped data
*/
func UpdateDatabase(clusterId string, nodeMetrics *v1beta1.NodeMetricsList, podMetrics *v1beta1.PodMetricsList) error {
	metricsDB, err := sql.Open("sqlite3", "./k8s_metrics.db")
	if err != nil {
		fmt.Println("init metrics sqlite db error: ", err)
	}
	defer metricsDB.Close()
	tx, err := metricsDB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into nodes(uid, cid, name, cpu, memory, storage, time) values(?, ?, ?, ?, ?, ?, datetime('now'))")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range nodeMetrics.Items {
		_, err = stmt.Exec(v.UID, clusterId, v.Name, v.Usage.Cpu().MilliValue(), v.Usage.Memory().Value()/1024.0/1024.0, v.Usage.StorageEphemeral().Value())
		if err != nil {
			return err
		}
	}

	stmt, err = tx.Prepare("insert into pods(uid, cid, name, namespace, container, cpu, memory, storage, time) values(?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range podMetrics.Items {
		for _, u := range v.Containers {
			_, err = stmt.Exec(v.UID, clusterId, v.Name, v.Namespace, u.Name, u.Usage.Cpu().MilliValue(), u.Usage.Memory().Value()/1024.0/1024.0, u.Usage.StorageEphemeral().Value())
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()

	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return rberr
		}
		return err
	}

	return nil
}

/*
	CullDatabase deletes rows from nodes and pods based on a time window.
*/
func CullDatabase(window *time.Duration) error {
	metricsDB, err := sql.Open("sqlite3", "./k8s_metrics.db")
	if err != nil {
		fmt.Println("init metrics sqlite db error: ", err)
	}
	defer metricsDB.Close()
	tx, err := metricsDB.Begin()
	if err != nil {
		return err
	}

	windowStr := fmt.Sprintf("-%.0f seconds", window.Seconds())

	nodestmt, err := tx.Prepare("delete from nodes where time <= datetime('now', ?);")
	if err != nil {
		return err
	}

	defer nodestmt.Close()
	res, err := nodestmt.Exec(windowStr)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	fmt.Printf("Cleaning up nodes: %d rows removed \n", affected)

	podstmt, err := tx.Prepare("delete from pods where time <= datetime('now', ?);")
	if err != nil {
		return err
	}

	defer podstmt.Close()
	res, err = podstmt.Exec(windowStr)
	if err != nil {
		return err
	}

	affected, _ = res.RowsAffected()
	fmt.Printf("Cleaning up pods: %d rows removed \n", affected)
	err = tx.Commit()

	if err != nil {
		rberr := tx.Rollback()
		if rberr != nil {
			return rberr
		}
		return err
	}

	return nil
}

func getRows(clusterId string, table string, metricName string, selector ResourceSelector) (*sql.Rows, error) {
	metricsDB, err := sql.Open("sqlite3", "./k8s_metrics.db")
	if err != nil {
		fmt.Println("init metrics sqlite db error: ", err)
	}
	defer metricsDB.Close()
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

	return metricsDB.Query(query, values...)
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

	defer rows.Close()

	resultList := make(map[string]SidecarMetric)

	for rows.Next() {
		var metricValue string
		var pod string
		var metricTime string
		var uid string
		var newMetric MetricPoint
		err = rows.Scan(&metricValue, &pod, &uid, &metricTime)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		layout := "2006-01-02T15:04:05Z"
		local, _ := time.LoadLocation("Local")
		t, err := time.ParseInLocation(layout, metricTime, local)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		v, err := strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		newMetric = MetricPoint{
			Timestamp: t.Add(time.Hour * time.Duration(8)).Format("15:04:05"),
			Value:     v,
		}

		if _, ok := resultList[pod]; ok {
			metricThing := resultList[pod]
			metricThing.AddMetricPoint(newMetric)
			resultList[pod] = metricThing
		} else {
			resultList[pod] = SidecarMetric{
				MetricName:   metricName,
				MetricPoints: []MetricPoint{newMetric},
				DataPoints:   []DataPoint{},
				UIDs: []types.UID{
					types.UID(pod),
				},
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return SidecarMetricResultList{}, err
	}

	result := SidecarMetricResultList{}
	for _, v := range resultList {
		result.Items = append(result.Items, v)
	}

	return result, nil
}

/*
	getNodeMetrics: With a database connection and a resource selector
	Queries SQLite and returns a list of metrics.
*/
func GetNodeMetrics(clusterId string, metricName string, selector ResourceSelector) (SidecarMetricResultList, error) {
	resultList := make(map[string]SidecarMetric)
	rows, err := getRows(clusterId, "nodes", metricName, selector)

	if err != nil {
		fmt.Printf("Error getting node metrics: %v \n", err)
		return SidecarMetricResultList{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var metricValue string
		var node string
		var metricTime string
		var uid string
		var newMetric MetricPoint
		err = rows.Scan(&metricValue, &node, &uid, &metricTime)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		layout := "2006-01-02T15:04:05Z"
		local, _ := time.LoadLocation("Local")
		t, err := time.ParseInLocation(layout, metricTime, local)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		v, err := strconv.ParseUint(metricValue, 10, 64)
		if err != nil {
			return SidecarMetricResultList{}, err
		}

		newMetric = MetricPoint{
			Timestamp: t.Add(time.Hour * time.Duration(8)).Format("15:04:05"),
			Value:     v,
		}

		if _, ok := resultList[node]; ok {
			metricThing := resultList[node]
			metricThing.AddMetricPoint(newMetric)
			resultList[node] = metricThing
		} else {
			resultList[node] = SidecarMetric{
				MetricName:   metricName,
				MetricPoints: []MetricPoint{newMetric},
				DataPoints:   []DataPoint{},
				UIDs: []types.UID{
					types.UID(node),
				},
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return SidecarMetricResultList{}, err
	}

	result := SidecarMetricResultList{}
	for _, v := range resultList {
		result.Items = append(result.Items, v)
	}

	return result, nil
}
