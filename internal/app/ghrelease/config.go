package ghrelease

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type GhRelease struct {
	Owner string
	Repo  string

	Title       string
	Description string
	Version     string

	Assets []GhReleaseAsset
}

type GhReleaseAsset string

func (a GhReleaseAsset) GetFiles() []string {
	m, err := filepath.Glob(string(a))
	if err != nil {
		panic(err)
	}
	logrus.WithField("files", m).Trace("Found files")

	return m
}
