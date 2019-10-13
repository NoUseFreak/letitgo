package utils

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	Owner  string
	Repo   string
	ctx    context.Context
	client *github.Client
}

func (c *GithubClient) init() error {
	if DryRun {
		return &e.SkipError{"dryrun"}
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("Make sure to set GITHUB_TOKEN")
	}

	c.ctx = context.Background()
	c.client = c.ghClient(token)

	return nil
}

func (c *GithubClient) CreateRelease(version, title, description string) (int64, error) {
	if err := c.init(); err != nil {
		return -1, err
	}

	release, _, err := c.client.Repositories.GetReleaseByTag(
		c.ctx,
		c.Owner,
		c.Repo,
		version,
	)
	if release != nil && err == nil {
		ui.Debug("Release exists")
		return release.GetID(), nil
	}

	ui.Debug("Creating release")
	release, _, err = c.client.Repositories.CreateRelease(c.ctx, c.Owner, c.Repo, &github.RepositoryRelease{
		TagName: &version,
		Name:    &title,
		Body:    &description,
	})

	if err != nil {
		return -1, err
	}

	return release.GetID(), err
}

func (c *GithubClient) UploadAssets(releaseID int64, assets []string) error {
	if err := c.init(); err != nil {
		return err
	}

	for _, f := range assets {
		ui.Debug("Uploading %s", f)
		if err := c.uploadAsset(releaseID, f); err != nil {
			return err
		}
	}

	return nil
}

func (c *GithubClient) ghClient(token string) *github.Client {
	tc := oauth2.NewClient(c.ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	return github.NewClient(tc)
}

func (c *GithubClient) uploadAsset(rID int64, path string) error {
	ui.Trace("Uploading %s", path)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, _, err = c.client.Repositories.UploadReleaseAsset(c.ctx, c.Owner, c.Repo, rID, &github.UploadOptions{
		Name: filepath.Base(path),
	}, file)

	return err
}
