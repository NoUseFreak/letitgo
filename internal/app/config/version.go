package config

import (
	"strings"

	sv "github.com/coreos/go-semver/semver"
)

// Version is an struct representing the version of the project you are
// releasing.
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

// Major returns the major part of a semver version.
func (v Version) Major() int64 {
	return v.semver.Major
}

// Minor returns the minor part of a semver version.
func (v Version) Minor() int64 {
	return v.semver.Minor
}

// Patch returns the patch part of a semver version.
func (v Version) Patch() int64 {
	return v.semver.Patch
}

// PreRelease returns the prerelease part of a semver version.
func (v Version) PreRelease() string {
	return strings.Join(v.semver.PreRelease.Slice(), "-")
}
