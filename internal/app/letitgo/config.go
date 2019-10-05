package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/changelog"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/githubrelease"
	"github.com/NoUseFreak/letitgo/internal/app/homebrew"
)

type Config struct {
	LetItGo config.LetItGoConfig

	Homebrew      []homebrew.Config
	GithubRelease []githubrelease.Config
	Changelog     []changelog.Config
}
