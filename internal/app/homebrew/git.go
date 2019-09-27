package homebrew

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func (h *Homebrew) saveToGit(path, content string) error {
	r := regexp.MustCompile("^.*github.com[/:](?P<org>[^/]+)/(?P<name>[^/\\.]+)")
	if !r.MatchString(h.Tap.URL) {
		return errors.New("Only github supported")
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("Make sure to set GITHUB_TOKEN")
	}

	matches := r.FindStringSubmatch(h.Tap.URL)
	owner := matches[1]
	repo := matches[2]

	message := fmt.Sprintf("Upgrade %s to %s", h.Name, h.Version)

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))

	client := github.NewClient(tc)

	var sha string
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)
	if err == nil {
		sha = fileContent.GetSHA()
	}

	_, _, err = client.Repositories.CreateFile(
		ctx,
		owner,
		repo,
		path,
		&github.RepositoryContentFileOptions{
			SHA:     &sha,
			Message: &message,
			Content: []byte(content),
		},
	)

	if err != nil {
		return err
	}

	return nil
}
