package k8s_structs

type NamespaceForm struct {
	Namespace string `json:"namespace" form:"namespace"`
	ResName   string `json:"resName" form:"resName"`
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

type PodContainerLogForm struct {
	Namespace     string `json:"namespace" form:"namespace"`
	PodName       string `json:"podName" form:"podName"`
	ContainerName string `json:"containerName" form:"containerName"`
}

type PodLogForm struct {
	Namespace     string `json:"namespace" form:"namespace"`
	PodName       string `json:"podName" form:"podName"`
	ContainerName string `json:"containerName" form:"containerName"`
}

type MetricsForm struct {
	ClusterId string `json:"clusterId" form:"clusterId"`
	MetricName string `json:"metricName" form:"metricName"`
	NodeName string `json:"nodeName" form:"nodeName"`
	Namespace string `json:"namespace" form:"namespace"`
	PodName string `json:"podName" form:"podName"`
}
