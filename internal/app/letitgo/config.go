package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/githubrelease"
	"github.com/NoUseFreak/letitgo/internal/app/homebrew"
)

type Config struct {
	Name        string
	Description string
	Version     string

	Homebrew      []homebrew.Config
	GithubRelease []githubrelease.Config
}
