package githubrelease

import (
	"reflect"
	"testing"
)

func TestAssetConfig_GetFiles(t *testing.T) {
	root := "../../../../"
	tests := []struct {
		name string
		a    assetConfig
		want []string
	}{
		{
			name: "Find single",
			a:    assetConfig(root + "test/assets/1.txt"),
			want: []string{
				root + "test/assets/1.txt",
			},
		},
		{
			name: "Find glob",
			a:    assetConfig(root + "test/assets/folder/*"),
			want: []string{
				root + "test/assets/folder/2.txt",
				root + "test/assets/folder/3.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.GetFiles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GhReleaseAsset.GetFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
