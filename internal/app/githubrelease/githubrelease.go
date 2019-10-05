package githubrelease

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/NoUseFreak/letitgo/internal/app/utils"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func Execute(c Config) error {
	templateProps(&c)

	if c.Owner == "" || c.Repo == "" {
		resolveOwnerRepo(&c)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("Make sure to set GITHUB_TOKEN")
	}

	ctx := context.Background()
	ghClient := ghClient(ctx, token)
	rID, err := createRelease(ctx, ghClient, c)
	if err != nil {
		return err
	}

	for _, a := range c.Assets {
		for _, f := range a.GetFiles() {
			logrus.Debugf("Uploading %s", f)

			if err = uploadAsset(ctx, ghClient, rID, f, c); err != nil {
				return err
			}
		}
	}

	return nil
}

func ghClient(ctx context.Context, token string) *github.Client {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	return github.NewClient(tc)
}

func createRelease(ctx context.Context, client *github.Client, c Config) (int64, error) {
	version := c.BaseConfig.LetItGo.Version()
	release, _, err := client.Repositories.GetReleaseByTag(ctx, c.Owner, c.Repo, version)
	if release != nil && err == nil {
		logrus.Debug("Release exists")
		return release.GetID(), nil
	}

	logrus.Debug("Creating release")
	release, _, err = client.Repositories.CreateRelease(ctx, c.Owner, c.Repo, &github.RepositoryRelease{
		TagName: &version,
		Name:    &c.Title,
		Body:    &c.Description,
	})

	if err != nil {
		return -1, err
	}

	return release.GetID(), err
}

func uploadAsset(ctx context.Context, client *github.Client, rID int64, path string, c Config) error {
	logrus.Tracef("Uploading %s", path)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, _, err = client.Repositories.UploadReleaseAsset(ctx, c.Owner, c.Repo, rID, &github.UploadOptions{
		Name: filepath.Base(path),
	}, file)

	return err
}

func templateProps(c *Config) {
	utils.TemplateProperty(&c.Title, c)
	utils.TemplateProperty(&c.Description, c)
}

func resolveOwnerRepo(c *Config) error {
	repo, err := utils.GetRemote(".")
	if err != nil {
		return err
	}

	uri, err := utils.ParseURI(repo)
	if err != nil {
		return err
	}

	c.Owner = uri.Owner
	c.Repo = uri.Repo

	return nil
}
