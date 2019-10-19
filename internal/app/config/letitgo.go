package config

// LetItGoConfig is the base config.
type LetItGoConfig struct {
	Name        string
	Description string

	version string
}

// NewConfig creates a config with the version set.
func NewConfig(v string) LetItGoConfig {
	return LetItGoConfig{
		version: v,
	}
}

// Version returns the version.
func (lig *LetItGoConfig) Version() string {
	return lig.version
}
