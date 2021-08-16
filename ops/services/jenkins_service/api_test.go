package jenkins_service

import (
	"reflect"
	"testing"
)

// func Test_jenkinsClient_CreateSystemCredentials(t *testing.T) {
// 	type fields struct {
// 		authHost  string
// 		authUser  string
// 		authToken string
// 	}
// 	type args struct {
// 		credentialsId string
// 		username      string
// 		passwd        string
// 		description   string
// 	}
// 	tests := []struct {
// 		name        string
// 		fields      fields
// 		args        args
// 		wantJobInfo JobInfo
// 		wantErr     bool
// 	}{
// 		{
// 			name: "create credential",
// 			fields: fields{
// 				authHost:  "http://10.0.0.31:9090",
// 				authUser:  "admin",
// 				authToken: "11b2183f37905ab28164bd0cea60cf7d90",
// 			},
// 			args: args{
// 				credentialsId: "测试凭证",
// 				username:      "123",
// 				passwd:        "211",
// 				description:   "凭证描述",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := jenkinsClient{
// 				authHost:  tt.fields.authHost,
// 				authUser:  tt.fields.authUser,
// 				authToken: tt.fields.authToken,
// 			}
// 			gotJobInfo, err := c.CreateSystemCredentials(tt.args.credentialsId, tt.args.username, tt.args.passwd, tt.args.description)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("jenkinsClient.CreateSystemCredentials() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotJobInfo, tt.wantJobInfo) {
// 				t.Errorf("jenkinsClient.CreateSystemCredentials() = %v, want %v", gotJobInfo, tt.wantJobInfo)
// 			}
// 		})
// 	}
// }

// func Test_jenkinsClient_GetSystemCredentialsInfo(t *testing.T) {
// 	type fields struct {
// 		authHost  string
// 		authUser  string
// 		authToken string
// 	}
// 	type args struct {
// 		credentialsId string
// 	}
// 	tests := []struct {
// 		name               string
// 		fields             fields
// 		args               args
// 		wantCredentialInfo CredentailsIndo
// 		wantErr            bool
// 	}{
// 		{
// 			name: "get credential",
// 			fields: fields{
// 				authHost:  "http://10.0.0.31:9090",
// 				authUser:  "admin",
// 				authToken: "11b2183f37905ab28164bd0cea60cf7d90",
// 			},
// 			args: args{
// 				credentialsId: "测试凭证",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := jenkinsClient{
// 				authHost:  tt.fields.authHost,
// 				authUser:  tt.fields.authUser,
// 				authToken: tt.fields.authToken,
// 			}
// 			gotCredentialInfo, err := c.GetSystemCredentialsInfo(tt.args.credentialsId)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("jenkinsClient.GetSystemCredentialsInfo() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotCredentialInfo, tt.wantCredentialInfo) {
// 				t.Errorf("jenkinsClient.GetSystemCredentialsInfo() = %v, want %v", gotCredentialInfo, tt.wantCredentialInfo)
// 			}
// 		})
// 	}
// }

// func Test_jenkinsClient_CreateJobItem(t *testing.T) {
// 	type fields struct {
// 		authHost  string
// 		authUser  string
// 		authToken string
// 	}
// 	type args struct {
// 		jobName    string
// 		configData models.CiConfigInfo
// 	}
// 	tests := []struct {
// 		name        string
// 		fields      fields
// 		args        args
// 		wantJobInfo JobInfo
// 		wantErr     bool
// 	}{
// 		{
// 			name: "get credential",
// 			fields: fields{
// 				authHost:  "http://10.0.0.31:9090",
// 				authUser:  "admin",
// 				authToken: "11b2183f37905ab28164bd0cea60cf7d90",
// 			},
// 			args: args{
// 				jobName: "test_job",
// 				configData: models.CiConfigInfo{
// 					GitRepo:                 "http://10.0.0.43:8090/test.git",
// 					GitBranch:               "dev",
// 					GitCredentialId:         "123",
// 					DockerImageRepo:         "aliyun",
// 					DockerImageCredentialId: "234",
// 					DockerImageName:         "test_image",
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := jenkinsClient{
// 				authHost:  tt.fields.authHost,
// 				authUser:  tt.fields.authUser,
// 				authToken: tt.fields.authToken,
// 			}
// 			gotJobInfo, err := c.CreateJobItem(tt.args.jobName, tt.args.configData)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("jenkinsClient.CreateJobItem() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotJobInfo, tt.wantJobInfo) {
// 				t.Errorf("jenkinsClient.CreateJobItem() = %v, want %v", gotJobInfo, tt.wantJobInfo)
// 			}
// 		})
// 	}
// }

func Test_jenkinsClient_GetJobItemConfig(t *testing.T) {
	type fields struct {
		authHost  string
		authUser  string
		authToken string
	}
	type args struct {
		jobName string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantJobInfo JobInfo
		wantErr     bool
	}{
		{
			name: "get credential",
			fields: fields{
				authHost:  "http://10.0.0.31:9090",
				authUser:  "admin",
				authToken: "",
			},
			args: args{
				jobName: "java_test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := jenkinsClient{
				authHost:  tt.fields.authHost,
				authUser:  tt.fields.authUser,
				authToken: tt.fields.authToken,
			}
			gotJobInfo, err := c.GetJobItemConfig(tt.args.jobName)
			if (err != nil) != tt.wantErr {
				t.Errorf("jenkinsClient.GetJobItemConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJobInfo, tt.wantJobInfo) {
				t.Errorf("jenkinsClient.GetJobItemConfig() = %v, want %v", gotJobInfo, tt.wantJobInfo)
			}
		})
	}
}

// func Test_jenkinsClient_UpdateJobItem(t *testing.T) {
// 	type fields struct {
// 		authHost  string
// 		authUser  string
// 		authToken string
// 	}
// 	type args struct {
// 		jobName    string
// 		configData models.CiConfigInfo
// 	}
// 	tests := []struct {
// 		name        string
// 		fields      fields
// 		args        args
// 		wantJobInfo JobInfo
// 		wantErr     bool
// 	}{
// 		{
// 			name: "get credential",
// 			fields: fields{
// 				authHost:  "http://10.0.0.31:9090",
// 				authUser:  "admin",
// 				authToken: "11b2183f37905ab28164bd0cea60cf7d90",
// 			},
// 			args: args{
// 				jobName: "test_job",
// 				configData: models.CiConfigInfo{
// 					GitRepo:                 "http://10.0.0.43:8090/test.git",
// 					GitBranch:               "dev",
// 					GitCredentialId:         "12222",
// 					DockerImageRepo:         "aliyun",
// 					DockerImageCredentialId: "23333",
// 					DockerImageName:         "test_image_xxx",
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := jenkinsClient{
// 				authHost:  tt.fields.authHost,
// 				authUser:  tt.fields.authUser,
// 				authToken: tt.fields.authToken,
// 			}
// 			gotJobInfo, err := c.UpdateJobItem(tt.args.jobName, tt.args.configData)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("jenkinsClient.UpdateJobItem() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotJobInfo, tt.wantJobInfo) {
// 				t.Errorf("jenkinsClient.UpdateJobItem() = %v, want %v", gotJobInfo, tt.wantJobInfo)
// 			}
// 		})
// 	}
// }

// func Test_jenkinsClient_GetJobBuildList(t *testing.T) {
// 	type fields struct {
// 		authHost  string
// 		authUser  string
// 		authToken string
// 	}
// 	type args struct {
// 		jobName string
// 	}
// 	tests := []struct {
// 		name        string
// 		fields      fields
// 		args        args
// 		wantJobInfo JobInfo
// 		wantErr     bool
// 	}{
// 		{
// 			name: "get credential",
// 			fields: fields{
// 				authHost:  "http://10.0.0.31:9090",
// 				authUser:  "admin",
// 				authToken: "11b2183f37905ab28164bd0cea60cf7d90",
// 			},
// 			args: args{
// 				jobName: "java_test",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := jenkinsClient{
// 				authHost:  tt.fields.authHost,
// 				authUser:  tt.fields.authUser,
// 				authToken: tt.fields.authToken,
// 			}
// 			gotJobInfo, err := c.GetJobBuildList(tt.args.jobName)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("jenkinsClient.GetJobBuildList() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotJobInfo, tt.wantJobInfo) {
// 				t.Errorf("jenkinsClient.GetJobBuildList() = %v, want %v", gotJobInfo, tt.wantJobInfo)
// 			}
// 		})
// 	}
// }

// func Test_jenkinsClient_GetBuildInfo(t *testing.T) {
// 	type fields struct {
// 		authHost  string
// 		authUser  string
// 		authToken string
// 	}
// 	type args struct {
// 		jobName string
// 		number  int32
// 	}
// 	tests := []struct {
// 		name          string
// 		fields        fields
// 		args          args
// 		wantBuildInfo BuildInfo
// 		wantErr       bool
// 	}{
// 		{
// 			name: "get credential",
// 			fields: fields{
// 				authHost:  "http://10.0.0.31:9090",
// 				authUser:  "admin",
// 				authToken: "11b2183f37905ab28164bd0cea60cf7d90",
// 			},
// 			args: args{
// 				jobName: "java_test",
// 				number:  9,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := jenkinsClient{
// 				authHost:  tt.fields.authHost,
// 				authUser:  tt.fields.authUser,
// 				authToken: tt.fields.authToken,
// 			}
// 			gotBuildInfo, err := c.GetBuildInfo(tt.args.jobName, tt.args.number)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("jenkinsClient.GetBuildInfo() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotBuildInfo, tt.wantBuildInfo) {
// 				t.Errorf("jenkinsClient.GetBuildInfo() = %v, want %v", gotBuildInfo, tt.wantBuildInfo)
// 			}
// 		})
// 	}
// }

// func Test_jenkinsClient_GetBuildNumLog(t *testing.T) {
// 	type fields struct {
// 		authHost  string
// 		authUser  string
// 		authToken string
// 	}
// 	type args struct {
// 		jobName     string
// 		buildNumber int32
// 	}
// 	tests := []struct {
// 		name          string
// 		fields        fields
// 		args          args
// 		wantBuildInfo BuildInfo
// 		wantErr       bool
// 	}{
// 		{
// 			name: "get credential",
// 			fields: fields{
// 				authHost:  "http://10.0.0.31:9090",
// 				authUser:  "admin",
// 				authToken: "11b2183f37905ab28164bd0cea60cf7d90",
// 			},
// 			args: args{
// 				jobName: "java_test",
// 				buildNumber:  9,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := jenkinsClient{
// 				authHost:  tt.fields.authHost,
// 				authUser:  tt.fields.authUser,
// 				authToken: tt.fields.authToken,
// 			}
// 			gotBuildInfo, err := c.GetBuildNumLog(tt.args.jobName, tt.args.buildNumber)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("jenkinsClient.GetBuildNumLog() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotBuildInfo, tt.wantBuildInfo) {
// 				t.Errorf("jenkinsClient.GetBuildNumLog() = %v, want %v", gotBuildInfo, tt.wantBuildInfo)
// 			}
// 		})
// 	}
// }

// func Test_jenkinsClient_GetJobList(t *testing.T) {
// 	type fields struct {
// 		authHost  string
// 		authUser  string
// 		authToken string
// 	}
// 	tests := []struct {
// 		name               string
// 		fields             fields
// 		wantJenkinsJobInfo JenkinsJobInfo
// 		wantErr            bool
// 	}{
// 		{
// 			name: "get credential",
// 			fields: fields{
// 				authHost:  "http://10.0.0.31:9090",
// 				authUser:  "admin",
// 				authToken: "11b2183f37905ab28164bd0cea60cf7d90",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			c := jenkinsClient{
// 				authHost:  tt.fields.authHost,
// 				authUser:  tt.fields.authUser,
// 				authToken: tt.fields.authToken,
// 			}
// 			gotJenkinsJobInfo, err := c.GetJobList()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("jenkinsClient.GetJobList() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotJenkinsJobInfo, tt.wantJenkinsJobInfo) {
// 				t.Errorf("jenkinsClient.GetJobList() = %v, want %v", gotJenkinsJobInfo, tt.wantJenkinsJobInfo)
// 			}
// 		})
// 	}
// }
