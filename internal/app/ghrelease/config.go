package ghrelease

import (
	"path/filepath"
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
	return m
}
