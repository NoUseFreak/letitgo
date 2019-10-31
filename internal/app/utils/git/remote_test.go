package git

import (
	"strings"
	"testing"
)

func TestGetRemote(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		want    string
		wantErr bool
	}{
		{
			name:    "current",
			dir:     "../../../../",
			want:    "NoUseFreak/letitgo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRemote(tt.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.Index(got, tt.want) < 1 {
				t.Errorf("GetRemote() = %v, needs to contain %v", got, tt.want)
			}
		})
	}
}
