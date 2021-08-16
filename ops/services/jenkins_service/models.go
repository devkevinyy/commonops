package jenkins_service

type JobInfo struct {
	FullName              string      `json:"fullName"`
	DisplayName           string      `json:"displayName"`
	Description           string      `json:"description"`
	Color                 string      `json:"color"`
	Buildable             bool        `json:"buildable"`
	FirstBuild            BuildItem   `json:"firstBuild"`
	LastBuild             BuildItem   `json:"lastBuild"`
	LastSuccessfulBuild   BuildItem   `json:"lastSuccessfulBuild"`
	LastUnsuccessfulBuild BuildItem   `json:"lastUnsuccessfulBuild"`
	Builds                []BuildItem `json:"builds"`
	BuildsDetails         []BuildInfo `json:"buildsDetails"`
}
type BuildItem struct {
	Number int32  `json:"number"`
	Url    string `json:"url"`
}

type BranchItem struct {
	SHA1 string `json:"SHA1"`
	Name string `json:"name"`
}
type LastBuiltRevision struct {
	SHA1   string       `json:"SHA1"`
	Branch []BranchItem `json:"branch"`
}

type ActionItem struct {
	Class             string            `json:"_class"`
	LastBuiltRevision LastBuiltRevision `json:"lastBuiltRevision"`
	RemoteUrls        []string          `json:"remoteUrls"`
}

type BuildInfo struct {
	DisplayName       string       `json:"displayName"`
	Description       string       `json:"description"`
	Duration          int32        `json:"duration"`
	EstimatedDuration int32        `json:"estimatedDuration"`
	FullDisplayName   string       `json:"fullDisplayName"`
	Id                string       `json:"id"`
	Building          bool         `json:"building"`
	KeepLog           bool         `json:"keepLog"`
	Result            string       `json:"result"`
	Timestamp         int64        `json:"timestamp"`
	Actions           []ActionItem `json:"actions"`
}

type CrumbItem struct {
	Crumb string `json:"crumb"`
}

type CredentailsInfo struct {
	Class       string `json:"_class"`
	Description string `json:"description"`
	DisplayName string `json:"displayName"`
	FullName    string `json:"fullName"`
	Id          string `json:"id"`
	TypeName    string `json:"typeName"`
}

type CredentialsList struct {
	Credentials []CredentailsInfo `json:"credentials"`
}

type JobItem struct {
	Class string `json:"_class"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Color string `json:"color"`
}

type JenkinsJobInfo struct {
	Jobs []JobItem `json:"jobs"`
}

type StageFlowNodes struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	Status               string `json:"status"`
	ParameterDescription string `json:"parameterDescription"`
	DurationMillis       int64  `json:"durationMillis"`
	Log                  string `json:"log"`
}

type BuildStageDescription struct {
	Name           string           `json:"name"`
	Status         string           `json:"status"`
	DurationMillis int64            `json:"durationMillis"`
	StageFlowNodes []StageFlowNodes `json:"stageFlowNodes"`
}

type StageFlowNodeLog struct {
	NodeId     string `json:"nodeId"`
	NodeStatus string `json:"nodeStatus"`
	HasMore    bool   `json:"hasMore"`
	Text       string `json:"text"`
}
