package homebrew

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Homebrew) buildURLHash(url string) (string, error) {
	logrus.Tracef("Downloading %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Not found")
	}

	hasher := sha256.New()
	io.Copy(hasher, resp.Body)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
