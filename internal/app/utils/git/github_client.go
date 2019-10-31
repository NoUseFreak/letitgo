package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/NoUseFreak/letitgo/internal/app/ui"
	"github.com/google/go-github/github"
	"github.com/matryer/try"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type githubClient struct {
	owner string
	repo  string

	client *github.Client
	ctx    context.Context
}

func NewGithubClient(owner, repo string) (*githubClient, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, errors.New("Make sure to set GITHUB_TOKEN")
	}

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	return &githubClient{
		owner: owner,
		repo:  repo,

		ctx:    ctx,
		client: github.NewClient(tc),
	}, nil
}

func (c *githubClient) PublishFile(path, content, message string, branch *string) error {
	opts := github.RepositoryContentGetOptions{}
	if branch != nil {
		opts.Ref = *branch
		if err := c.CreateBranch(*branch); err != nil {
			return err
		}
	}

	var sha string
	fileContent, _, _, err := c.client.Repositories.GetContents(
		c.ctx,
		c.owner,
		c.repo,
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
		c.owner,
		c.repo,
		path,
		&github.RepositoryContentFileOptions{
			SHA:     &sha,
			Message: &message,
			Content: []byte(content),
			Branch:  branch,
		},
	)

	return errors.Wrapf(err, "failed creating/update file %s", path)
}

func (c *githubClient) CreateRelease(version, title, description string) (string, error) {
	release, _, err := c.client.Repositories.GetReleaseByTag(
		c.ctx,
		c.owner,
		c.repo,
		version,
	)
	if release != nil && err == nil {
		return strconv.FormatInt(release.GetID(), 10), nil
	}

	release, _, err = c.client.Repositories.CreateRelease(
		c.ctx,
		c.owner,
		c.repo,
		&github.RepositoryRelease{
			TagName: &version,
			Name:    &title,
			Body:    &description,
		},
	)

	if err != nil {
		return "", errors.Wrap(err, "Failed creating release.")
	}

	return strconv.FormatInt(release.GetID(), 10), err
}

func (c *githubClient) UploadReleaseAssets(releaseID string, assets []string) error {
	rID, err := strconv.ParseInt(releaseID, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "failed to parse releaseID %s", releaseID)
	}

	for _, a := range assets {
		if err := c.uploadReleaseAsset(rID, a); err != nil {
			return err
		}
	}

	return nil
}

func (c *githubClient) uploadReleaseAsset(rID int64, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "Failed to open file %v", path)
	}
	defer func() { _ = file.Close() }()

	_, _, err = c.client.Repositories.UploadReleaseAsset(
		c.ctx,
		c.owner,
		c.repo,
		rID,
		&github.UploadOptions{
			Name: filepath.Base(path),
		},
		file,
	)

	return err
}

func (c *githubClient) CreateForkFrom(owner, repo string) error {
	_, _, err := c.client.Repositories.CreateFork(c.ctx, owner, repo, &github.RepositoryCreateForkOptions{})

	if _, ok := err.(*github.AcceptedError); !ok {
		if err != nil {
			return errors.Wrapf(err, "Failed to fork %s/%s to %s/%s", owner, repo, c.owner, repo)
		}
	}

	err = try.Do(func(attempt int) (bool, error) {
		time.Sleep(5 * time.Second)
		var err error
		if !c.Exists() {
			err = errors.New("does not exist")
		}

		return attempt < 6, err // try 5 times
	})

	return errors.Wrap(err, "Failed to verify fork got created.")
}
func (c *githubClient) Exists() bool {
	_, _, err := c.client.Repositories.Get(
		c.ctx,
		c.owner,
		c.repo,
	)

	return err == nil
}

func (c *githubClient) CreateBranch(branch string) error {
	ref := fmt.Sprintf("heads/%s", branch)

	r, _, _ := c.client.Git.GetRef(c.ctx, c.owner, c.repo, ref)
	if r == nil {
		r, _, err := c.client.Git.GetRef(c.ctx, c.owner, c.repo, "heads/master")
		if err != nil {
			return errors.Wrap(err, "Failed to find master branch sha")
		}

		_, _, err = c.client.Git.CreateRef(c.ctx, c.owner, c.repo, &github.Reference{
			Ref:    &ref,
			Object: r.GetObject(),
		})

		if err != nil {
			return errors.Wrapf(err, "Failed to create branch %s on %s/%s", branch, c.owner, c.repo)
		}
	}

	return nil
}
