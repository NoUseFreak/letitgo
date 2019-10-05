package config

type BaseConfig struct {
	LetItGo LetItGoConfig
}

func (bc *BaseConfig) Version() string { return bc.LetItGo.Version() }
