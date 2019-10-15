package changelog

import "github.com/NoUseFreak/letitgo/internal/app/config"

type Config struct {
	config.BaseConfig `yaml:"-"`

	File    string
	Message string
}
