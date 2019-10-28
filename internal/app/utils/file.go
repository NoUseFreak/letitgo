package utils

import "io/ioutil"

// WriteFile is a safe way to write to a file, adhering to dryRun.
func WriteFile(path, content string) error {
	if DryRun.IsEnabled() {
		return nil
	}

	return ioutil.WriteFile(path, []byte(content), 0644)
}
