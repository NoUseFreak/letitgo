package letitgo

import (
	"github.com/NoUseFreak/letitgo/internal/app/changelog"
	"github.com/NoUseFreak/letitgo/internal/app/config"
	"github.com/NoUseFreak/letitgo/internal/app/githubrelease"
	"github.com/NoUseFreak/letitgo/internal/app/homebrew"
	"github.com/NoUseFreak/letitgo/internal/app/snapcraft"
)

type Config struct {
	LetItGo config.LetItGoConfig

	Changelog     []changelog.Config
	GithubRelease []githubrelease.Config
	Homebrew      []homebrew.Config
	Snapcraft     []snapcraft.Config
}
