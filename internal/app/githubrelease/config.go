package githubrelease

import (
	"path/filepath"

	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/sirupsen/logrus"
)

type Config struct {
	config.BaseConfig

	Owner string
	Repo  string

	Title       string
	Description string

	Assets []GhReleaseAssetConfig
}

type GhReleaseAssetConfig string

func (a GhReleaseAssetConfig) GetFiles() []string {
	m, err := filepath.Glob(string(a))
	if err != nil {
		panic(err)
	}
	logrus.WithField("files", m).Trace("Found files")

	return m
}
