package homebrew

import (
	"errors"
	"fmt"
	"path"
)

func (h *Homebrew) Execute() error {
	h.templateInput()
	hash, err := h.buildURLHash(h.URL)
	if err != nil {
		return errors.New("Failed to build hash from url - " + err.Error())
	}
	h.Hash = hash

	content, err := h.template()
	if err != nil {
		return errors.New("Failed to template homebrew spec - " + err.Error())
	}
	filename := path.Join(h.Folder, fmt.Sprintf("%s.rb", h.Name))

	if err := h.saveToGit(filename, content); err != nil {
		return errors.New("Failed to save homebrew spec to git - " + err.Error())
	}

	return nil
}
