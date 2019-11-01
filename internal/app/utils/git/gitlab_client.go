package git

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type gitlabClient struct {
	client *gitlab.Client
	pID    string
}

func NewGitlabClient(owner, repo string) (*gitlabClient, error) {
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		return nil, errors.New("Make sure to set GITLAB_TOKEN")
	}

	return &gitlabClient{
		client: gitlab.NewClient(nil, token),
		pID:    fmt.Sprintf("%s/%s", owner, repo),
	}, nil
}

func (c *gitlabClient) PublishFile(path, content, message string, branch *string) error {
	_, _, err := c.client.RepositoryFiles.CreateFile(
		c.pID,
		path,
		&gitlab.CreateFileOptions{
			Content: gitlab.String(content),
			Branch:  branch,
		},
	)

	return err
}

func (c *gitlabClient) CreateForkFrom(owner, repo string) error {
	return errors.New("Not implemented")
}

func (c *gitlabClient) Exists() bool {
	return false
}

func (c *gitlabClient) CreateBranch(branch string) error {
	return errors.New("Not implemented")
}
