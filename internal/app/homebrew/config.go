package homebrew

import "github.com/NoUseFreak/letitgo/internal/app/config"

type Config struct {
	config.BaseConfig

	Name        string
	Description string
	Homepage    string

	URL string

	Dependencies []string
	Conflicts    []string

	Tap     HomebrewTapConfig
	Folder  string
	Install string
	Test    string

	// internal
	hash string
}

func (c *Config) Hash() string { return c.hash }

type HomebrewTapConfig struct {
	URL         string
	AuthorName  string
	AuthorEmail string
}
