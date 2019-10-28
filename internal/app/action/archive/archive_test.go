package archive

import (
	"os"
	"path"
	"testing"

	"github.com/NoUseFreak/letitgo/internal/app/config"
)

func Test_archive_Execute(t *testing.T) {
	rootDir := "../../../../"
	cfg := config.LetItGoConfig{
		Name:        "test",
		Description: "test",
	}

	type fields struct {
		Source string
		Target string
		Extras []string
		Method string
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantArchive string
	}{
		{
			name: "success",
			fields: fields{
				Source: path.Join(rootDir, "./test/assets"),
				Target: path.Join(rootDir, "./build/test123"),
				Method: "zip",
			},
			wantErr:     false,
			wantArchive: path.Join(rootDir, "./build/test123/assets.zip"),
		},
		{
			name: "doesnotexist",
			fields: fields{
				Source: path.Join(rootDir, "./test/doesnotexist"),
				Method: "zip",
			},
			wantErr: true,
		},
		{
			name: "unsupported method",
			fields: fields{
				Method: "non-existing",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.wantArchive != "" {
			_ = os.RemoveAll(tt.wantArchive)
		}
		t.Run(tt.name, func(t *testing.T) {
			c := &archive{
				Source: tt.fields.Source,
				Target: tt.fields.Target,
				Extras: tt.fields.Extras,
				Method: tt.fields.Method,
			}
			if err := c.Execute(cfg); (err != nil) != tt.wantErr {
				t.Errorf("archive.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			info, err := os.Stat(tt.wantArchive)
			if tt.wantArchive != "" && err != nil {
				t.Errorf("archive.Execute() archive was not created")
			}
			if tt.wantArchive != "" && info.Size() < 1 {
				t.Errorf("Archive size to small")
			}
		})
	}
}
