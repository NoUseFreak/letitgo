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

	owner := strings.ReplaceAll(filepath.Dir(url.Path), "/", "")
	basename := filepath.Base(url.Path)
	repo := strings.TrimSuffix(basename, filepath.Ext(basename))

	switch url.Hostname() {
	case "github.com":
		return NewGithubClient(owner, repo)
	case "gitlab.com":
		return NewGitlabClient(owner, repo)
	}

	return nil, fmt.Errorf("could not find client for %s", repoURL)
}
