package utils

import (
	"testing"
)

func TestExtractUriPath(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		//{name: "URI路径提取-1", args:args{uri: "/cloud/servers?page=1"}, want: "/cloud/servers"},
		//{name: "URI路径提取-2", args:args{uri: "/user/tokenRefresh"}, want: "/user/tokenRefresh"},
		{name: "URI路径提取-3", args: args{uri: "/daily_job?page=1"}, want: "/daily_job"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractUriPath(tt.args.uri); got != tt.want {
				t.Errorf("ExtractUriPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
