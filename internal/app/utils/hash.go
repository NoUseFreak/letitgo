package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func BuildURLHash(alg, url string) (string, error) {
	logrus.Tracef("Downloading %s", url)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Not found")
	}

	var hasher hash.Hash

	switch alg {
	case "sha256":
		hasher = sha256.New()
	default:
		return "", errors.New("Unknown hashing algorithm")
	}

	io.Copy(hasher, resp.Body)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
