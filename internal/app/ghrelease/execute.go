package ghrelease

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (r *GhRelease) Execute() error {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("Make sure to set GITHUB_TOKEN")
	}

	ctx := context.Background()
	ghClient := r.ghClient(ctx, token)
	rID := r.createRelease(ctx, ghClient)

	for _, a := range r.Assets {
		for _, f := range a.GetFiles() {
			logrus.Debugf("Uploading %s", f)
			r.uploadAsset(ctx, ghClient, rID, f)
		}
	}

	return nil
}

func (r *GhRelease) ghClient(ctx context.Context, token string) *github.Client {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	return github.NewClient(tc)
}

func (r *GhRelease) createRelease(ctx context.Context, client *github.Client) int64 {
	release, _, err := client.Repositories.GetReleaseByTag(ctx, r.Owner, r.Repo, r.Version)
	if release != nil && err == nil {
		logrus.Debug("Release exists")
		return release.GetID()
	}

	logrus.Debug("Creating release")
	release, _, err = client.Repositories.CreateRelease(ctx, r.Owner, r.Repo, &github.RepositoryRelease{
		TagName: &r.Version,
		Name:    &r.Title,
		Body:    &r.Description,
	})

	if err != nil {
		panic(err)
	}

	return release.GetID()
}

func (r *GhRelease) uploadAsset(ctx context.Context, client *github.Client, rID int64, path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, _, err = client.Repositories.UploadReleaseAsset(ctx, r.Owner, r.Repo, rID, &github.UploadOptions{
		Name: filepath.Base(path),
	}, file)

	if err != nil {
		logrus.Warn(err)
	}
}

func (r *GhRelease) templateInput() {
	r.interpolate(&r.Title)
	r.interpolate(&r.Description)
}

func (r *GhRelease) interpolate(v *string) {
	if tpl, err := template.New("value").Funcs(sprig.TxtFuncMap()).Parse(*v); err == nil {
		var out bytes.Buffer
		if err := tpl.Execute(&out, r); err == nil {
			*v = out.String()
		}
	}
}
