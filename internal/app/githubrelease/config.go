package githubrelease

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type assetConfig string

func (a assetConfig) GetFiles() []string {
	m, err := filepath.Glob(string(a))
	if err != nil {
		panic(err)
	}
	logrus.WithField("files", m).Trace("Found files")

	return m
}
