package helm

import "github.com/NoUseFreak/letitgo/internal/app/config"

type Config struct {
	config.BaseConfig

	Charts     []string
	BuildDir   string
	Repository string
}
