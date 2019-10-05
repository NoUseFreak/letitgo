package config

type LetItGoConfig struct {
	Name        string
	Description string

	version string
}

func NewConfig(v string) LetItGoConfig {
	return LetItGoConfig{
		version: v,
	}
}

func (lig *LetItGoConfig) Version() string {
	return lig.version
}
