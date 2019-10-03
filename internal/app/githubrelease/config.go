package githubrelease

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Owner string
	Repo  string

	Title       string
	Description string
	Version     string

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
