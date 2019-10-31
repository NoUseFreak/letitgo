package utils

import (
	"reflect"
	"testing"
)

func TestParseURI(t *testing.T) {
	want := GitURI{
		Host:  "github.com",
		Owner: "NoUseFreak",
		Repo:  "letitgo",
	}
	tests := []struct {
		name    string
		uri     string
		want    *GitURI
		wantErr bool
	}{
		{
			name:    "gh",
			uri:     "git@github.com/NoUseFreak/letitgo",
			want:    &want,
			wantErr: false,
		},
		{
			name:    "gh_ext",
			uri:     "git@github.com/NoUseFreak/letitgo.git",
			want:    &want,
			wantErr: false,
		},
		{
			name:    "gh_user",
			uri:     "github.com/NoUseFreak/letitgo",
			want:    &want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
