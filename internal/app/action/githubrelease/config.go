package githubrelease

import (
	"path/filepath"

	"github.com/NoUseFreak/letitgo/internal/app/ui"
)

type assetConfig string

func (a assetConfig) GetFiles() []string {
	m, err := filepath.Glob(string(a))
	if err != nil {
		ui.Error("Failed to read files")
		return []string{}
	}

	return m
}
