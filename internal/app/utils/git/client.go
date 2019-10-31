package git

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	giturls "github.com/whilp/git-urls"
)

type Client interface {
	PublishFile(path, content, message string, branch *string) error
	CreateRelease(version, title, description string) (string, error)
	UploadReleaseAssets(releaseID string, assets []string) error
	CreateForkFrom(owner, repo string) error
	CreateBranch(branch string) error
	Exists() bool
}

func NewClient(repoURL string, dryRun bool) (Client, error) {
	if dryRun {
		return NewDryRunClient()
	}
	url, err := giturls.Parse(repoURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed parsing giturl")
	}

	switch url.Hostname() {
	case "github.com":
		owner := strings.ReplaceAll(filepath.Dir(url.Path), "/", "")
		basename := filepath.Base(url.Path)
		repo := strings.TrimSuffix(basename, filepath.Ext(basename))

		return NewGithubClient(owner, repo)
	}

	return nil, fmt.Errorf("could not find client for %s", repoURL)
}
