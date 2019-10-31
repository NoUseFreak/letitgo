package git

import (
	"os"
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		repoURL string
		dryRun  bool
		env     map[string]string
	}
	tests := []struct {
		name     string
		args     args
		wantType Client
		wantErr  bool
	}{
		{
			name: "github",
			args: args{
				repoURL: "git@github.com:NoUseFreak/letitgo.git",
				env: map[string]string{
					"GITHUB_TOKEN": "123",
				},
			},
			wantType: &githubClient{},
			wantErr:  false,
		},
		{
			name: "github_notoken",
			args: args{
				repoURL: "git@github.com:NoUseFreak/letitgo.git",
				env: map[string]string{
					"GITHUB_TOKEN": "",
				},
			},
			wantType: nil,
			wantErr:  true,
		},
		{
			name: "unknown",
			args: args{
				repoURL: "git@somethingwrong.com:NoUseFreak/letitgo.git",
				env: map[string]string{
					"GITHUB_TOKEN": "",
				},
			},
			wantType: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.args.env {
				_ = os.Unsetenv(k)
				if v != "" {
					_ = os.Setenv(k, v)
				}
			}
			got, err := NewClient(tt.args.repoURL, tt.args.dryRun)
			for k := range tt.args.env {
				_ = os.Unsetenv(k)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(&got) != reflect.TypeOf(&tt.wantType) {
				t.Errorf("NewClient().(type) = %T, want %T", &got, &tt.wantType)
			}
		})
	}
}
