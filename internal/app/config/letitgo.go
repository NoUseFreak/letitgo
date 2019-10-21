package config

// LetItGoConfig is the base config.
type LetItGoConfig struct {
	Name        string
	Description string

	version Version
}

// NewConfig creates a config with the version set.
func NewConfig(v string) LetItGoConfig {
	return LetItGoConfig{
		version: newVersion(v),
	}
}

// Version returns the version.
func (lig *LetItGoConfig) Version() Version {
	return lig.version
}
