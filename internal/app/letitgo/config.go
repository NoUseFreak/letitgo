package letitgo

import (
	"io/ioutil"
	"log"

	"github.com/NoUseFreak/letitgo/internal/app/homebrew"
	"gopkg.in/yaml.v2"
)

func NewConfig(file string) Config {
	cfg := Config{file: file}
	cfg.parse()

	return cfg
}

type Config struct {
	file     string
	Homebrew []homebrew.Homebrew
}

func (c *Config) parse() {
	yamlFile, err := ioutil.ReadFile(c.file)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
