package snapcraft

import "github.com/NoUseFreak/letitgo/internal/app/config"

type Config struct {
	config.BaseConfig

	Name        string
	Description string

	Assets       []string
	Architecture string
}
