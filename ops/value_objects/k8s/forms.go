package k8s

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 3:03 PM
 * @Desc:
 */

type K8sForm struct {
	Id          int    `form:"id" json:"id"`
	ClusterId   string `form:"clusterId" json:"clusterId"`
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
	Token       string `form:"token" json:"token"`
	ApiServer   string `form:"apiServer" json:"apiServer"`
}

type NamespaceForm struct {
	Namespace string `json:"namespace" form:"namespace"`
	ResName   string `json:"resName" form:"resName"`
}

type PodContainerLogForm struct {
	Namespace     string `json:"namespace" form:"namespace"`
	PodName       string `json:"podName" form:"podName"`
	ContainerName string `json:"containerName" form:"containerName"`
}

type ResourceForm struct {
	Namespace    string `json:"namespace" form:"namespace"`
	ResType      string `json:"resType" form:"resType"`
	ResName      string `json:"resName" form:"resName"`
	ReplicaCount int32  `json:"replicaCount" form:"replicaCount"`
}

type YamlResource struct {
	NamespaceForm
	YamlContent string `json:"yamlContent" form:"yamlContent"`
}

type MetricsForm struct {
	ClusterId  string `json:"clusterId" form:"clusterId"`
	MetricName string `json:"metricName" form:"metricName"`
	NodeName   string `json:"nodeName" form:"nodeName"`
	Namespace  string `json:"namespace" form:"namespace"`
	PodName    string `json:"podName" form:"podName"`
}