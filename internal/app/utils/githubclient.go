package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
)

// GithubClient handles github api interactions.
type GithubClient struct {
	Owner  string
	Repo   string
	ctx    context.Context
	client *github.Client
}

func (c *GithubClient) init() error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("Make sure to set GITHUB_TOKEN")
	}

	c.ctx = context.Background()
	c.client = c.ghClient(token)

	return nil
}

// CreateRelease creates a release for a given version if it does not exist.
func (c *GithubClient) CreateRelease(version, title, description string) (int64, error) {
	if err := c.init(); err != nil {
		return -1, err
	}

	ui.Step("Validating release")
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

	ui.Step("Creating release")
	if DryRun {
		return -1, nil
	}
	release, _, err = c.client.Repositories.CreateRelease(c.ctx, c.Owner, c.Repo, &github.RepositoryRelease{
		TagName: &version,
		Name:    &title,
		Body:    &description,
	})

	if err != nil {
		return -1, errors.Wrap(err, "Failed creating release.")
	}

	return release.GetID(), err
}

// UploadAssets uploads multiple assets to a given release.
func (c *GithubClient) UploadAssets(releaseID int64, assets []string) error {
	if err := c.init(); err != nil {
		return errors.Wrap(err, "Failed to init client")
	}
	if DryRun {
		return &e.SkipError{
			Reason: "dryrun",
			Part:   fmt.Sprintf("Upload assets %s", strings.Join(assets, ", ")),
		}
	}

	for _, f := range assets {
		if err := c.UploadAsset(releaseID, f); err != nil {
			return errors.Wrap(err, "Failed to upload asset")
		}
	}

	return nil
}

// UploadAsset uploads a single asset to a given release.
func (c *GithubClient) UploadAsset(releaseID int64, asset string) error {
	ui.Step("Uploading %s", asset)
	if err := c.init(); err != nil {
		return errors.Wrap(err, "Failed to init client")
	}
	ui.Debug("Uploading %s", asset)
	if DryRun {
		return &e.SkipError{
			Reason: "dryrun",
			Part:   fmt.Sprintf("Upload asset %s", asset),
		}
	}

	return c.uploadAsset(releaseID, asset)
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
		return errors.Wrapf(err, "Failed to open file %v", path)
	}
	defer file.Close()

	_, _, err = c.client.Repositories.UploadReleaseAsset(c.ctx, c.Owner, c.Repo, rID, &github.UploadOptions{
		Name: filepath.Base(path),
	}, file)

	return err
}
