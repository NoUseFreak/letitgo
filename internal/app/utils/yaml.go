package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v1"
)

func ParseYamlFile(path string, data interface{}) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, data)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}
