package models

type CiConfigInfo struct {
	GitRepo                 string `json:"GitRepo"`
	GitBranch               string `json:"GitBranch"`
	GitCredentialId         string `json:"GitCredentialId"`
	DockerImageRepo         string `json:"DockerImageRepo"`
	DockerImageCredentialId string `json:"DockerImageCredentialId"`
	DockerImageName         string `json:"DockerImageName"`
}

type CiScriptInfo struct {
	PipelineScript string `json:"PipelineScript"`
}
