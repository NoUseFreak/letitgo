package utils

import "io/ioutil"

func WriteFile(path, content string) error {
	if DryRun.IsEnabled() {
		return nil
	}

	return ioutil.WriteFile(path, []byte(content), 0644)
}
