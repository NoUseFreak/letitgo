package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"net"
	"net/http"
	"time"

	e "github.com/NoUseFreak/letitgo/internal/app/errors"
	"github.com/NoUseFreak/letitgo/internal/app/ui"
)

// BuildURLHash return the hash of a remote file.
func BuildURLHash(alg, url string) (string, error) {
	ui.Trace("Downloading %s", url)
	client := http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,

			ExpectContinueTimeout: 4 * time.Second,
			ResponseHeaderTimeout: 3 * time.Second,
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer e.DeferCheck(resp.Body.Close)

	if resp.StatusCode != 200 {
		return "", errors.New("not found")
	}

	return BuildHash(alg, resp.Body)
}

// BuildHash returns the hash of a local file.
func BuildHash(alg string, reader io.Reader) (string, error) {

	var hasher hash.Hash
	switch alg {
	case "sha256":
		hasher = sha256.New()
	case "sha1":
		hasher = sha1.New()
	default:
		return "", errors.New("unknown hashing algorithm")
	}

	_, err := io.Copy(hasher, reader)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
