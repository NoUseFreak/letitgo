package utils

import "io/ioutil"

func WriteFile(path, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0644)
}
