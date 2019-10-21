package githubrelease

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type assetConfig string

func (a assetConfig) GetFiles() []string {
	m, err := filepath.Glob(string(a))
	if err != nil {
		logrus.Error("Failed to read files")
		return []string{}
	}
	logrus.WithField("files", m).Trace("Found files")

	return m
}
