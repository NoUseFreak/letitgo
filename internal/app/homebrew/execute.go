package homebrew

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path"
)

func (h *Homebrew) Execute() error {
	h.templateInput()
	hash, err := h.buildURLHash(h.URL)
	if err != nil {
		return err
	}
	h.Hash = hash

	content, err := h.template()
	if err != nil {
		return err
	}
	filename := path.Join(h.Folder, fmt.Sprintf("%s.rb", h.Name))

	if err := h.saveToGit(filename, content); err != nil {
		return err
	}

	return nil
}

func (h *Homebrew) buildURLHash(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	hasher := sha256.New()
	io.Copy(hasher, resp.Body)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
