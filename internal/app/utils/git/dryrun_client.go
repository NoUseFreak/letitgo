package git

import (
	"github.com/NoUseFreak/letitgo/internal/app/ui"
)

type dryRunClient struct {
}

func NewDryRunClient() (*dryRunClient, error) {
	return &dryRunClient{}, nil
}

func (c *dryRunClient) PublishFile(path, content, message string, branch *string) error {
	ui.Warn("DRYRUN: Skipping publish %s", path)
	return nil
}

func (c *dryRunClient) CreateForkFrom(owner, repo string) error {
	ui.Warn("DRYRUN: Skipping create fork %s/%s", owner, repo)
	return nil
}
func (c *dryRunClient) Exists() bool {
	return true
}

func (c *dryRunClient) CreateBranch(branch string) error {
	ui.Warn("DRYRUN: Skipping create branch %s", branch)
	return nil
}
