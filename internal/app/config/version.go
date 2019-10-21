package config

import (
	"strings"

	sv "github.com/coreos/go-semver/semver"
)

type Version struct {
	semver *sv.Version
}

func newVersion(v string) Version {
	return Version{
		semver: sv.New(v),
	}
}

func (v Version) String() string {
	return v.semver.String()
}

func (v Version) Major() int64 {
	return v.semver.Major
}

func (v Version) Minor() int64 {
	return v.semver.Minor
}

func (v Version) Patch() int64 {
	return v.semver.Patch
}

func (v Version) PreRelease() string {
	return strings.Join(v.semver.PreRelease.Slice(), "-")
}
