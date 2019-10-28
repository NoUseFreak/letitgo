package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
	try "gopkg.in/matryer/try.v1"
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
	if DryRun.IsEnabled() {
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
	if DryRun.IsEnabled() {
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
	if DryRun.IsEnabled() {
		return &e.SkipError{
			Reason: "dryrun",
			Part:   fmt.Sprintf("Upload asset %s", asset),
		}
	}

	return c.uploadAsset(releaseID, asset)
}

// RepoExists checks if the repo exists. It will return an error if it does not.
func (c *GithubClient) RepoExists() error {
	if err := c.init(); err != nil {
		return errors.Wrap(err, "Failed to init client")
	}
	_, _, err := c.client.Repositories.Get(c.ctx, c.Owner, c.Repo)
	if err != nil {
		return err
	}

	return nil
}

// CreateForkFrom will create a for of owner/repo to your github account.
func (c *GithubClient) CreateForkFrom(owner, repo string) error {
	if err := c.init(); err != nil {
		return errors.Wrap(err, "Failed to init client")
	}

	_, _, err := c.client.Repositories.CreateFork(c.ctx, owner, repo, &github.RepositoryCreateForkOptions{})

	if _, ok := err.(*github.AcceptedError); !ok {
		if err != nil {
			return errors.Wrapf(err, "Failed to fork %s/%s to %s/%s", owner, repo, c.Owner, repo)
		}
	}

	err = try.Do(func(attempt int) (bool, error) {
		time.Sleep(5 * time.Second)
		return attempt < 6, c.RepoExists() // try 5 times
	})

	return errors.Wrap(err, "Failed to verify fork got created.")
}

// PublishFile will commit changes to a single file.
func (c *GithubClient) PublishFile(path, content, message, branch string) error {
	if err := c.init(); err != nil {
		return errors.Wrap(err, "Failed to init client")
	}

	if branch != "" {
		// create branch
		ref := fmt.Sprintf("heads/%s", branch)

		r, _, _ := c.client.Git.GetRef(c.ctx, c.Owner, c.Repo, ref)
		if r == nil {
			r, _, err := c.client.Git.GetRef(c.ctx, c.Owner, c.Repo, "heads/master")
			if err != nil {
				return errors.Wrap(err, "Failed to find master branch sha")
			}

			_, _, err = c.client.Git.CreateRef(c.ctx, c.Owner, c.Repo, &github.Reference{
				Ref:    &ref,
				Object: r.GetObject(),
			})

			if err != nil {
				return errors.Wrapf(err, "Failed to create branch %s on %s/%s", branch, c.Owner, c.Repo)
			}
		}
	}

	var sha string
	opts := github.RepositoryContentGetOptions{}
	if branch != "" {
		opts.Ref = branch
	}
	fileContent, _, _, err := c.client.Repositories.GetContents(
		c.ctx,
		c.Owner,
		c.Repo,
		path,
		&opts,
	)
	if err == nil {
		if oldContent, err := fileContent.GetContent(); err == nil {
			if oldContent == content {
				ui.Step("No changes to %s", path)
				return nil
			}
		}
		sha = fileContent.GetSHA()
	}

	_, _, err = c.client.Repositories.CreateFile(
		c.ctx,
		c.Owner,
		c.Repo,
		path,
		&github.RepositoryContentFileOptions{
			SHA:     &sha,
			Message: &message,
			Content: []byte(content),
			Branch:  &branch,
		},
	)

	return errors.Wrapf(err, "Failed to publish file %s", path)
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
	defer DeferCheck(file.Close)

	_, _, err = c.client.Repositories.UploadReleaseAsset(c.ctx, c.Owner, c.Repo, rID, &github.UploadOptions{
		Name: filepath.Base(path),
	}, file)

	return err
}
