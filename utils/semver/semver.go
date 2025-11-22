package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var semverRegex = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z\-.]+))?(?:\+([0-9A-Za-z\-.]+))?$`)

// Version represents a semantic version.
type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Build      string
}

// Parse parses a semver string into a Version. Accepts an optional leading 'v'.
func Parse(s string) (Version, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Version{}, fmt.Errorf("empty version string")
	}
	m := semverRegex.FindStringSubmatch(s)
	if m == nil {
		return Version{}, fmt.Errorf("invalid semver: %q", s)
	}
	maj, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	pat, _ := strconv.Atoi(m[3])
	pr := m[4]
	bu := m[5]
	return Version{Major: maj, Minor: minor, Patch: pat, PreRelease: pr, Build: bu}, nil
}

// String returns the canonical string representation of the version (without a leading 'v').
func (v *Version) String() string {
	base := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.PreRelease != "" {
		base = base + "-" + v.PreRelease
	}
	if v.Build != "" {
		base = base + "+" + v.Build
	}
	return base
}

func (v *Version) IsDraft() bool {
	return v.PreRelease != ""
}

// IncrementPatch increments the patch version and resets pre-release and build metadata.
func (v *Version) IncrementPatch() {
	v.Patch++
	v.PreRelease = ""
	v.Build = ""
}

func (v *Version) IncrementMinor() {
	v.Minor++
	v.Patch = 0
	v.PreRelease = ""
	v.Build = ""
}

func (v *Version) IncrementMajor() {
	v.Major++
	v.Minor = 0
	v.Patch = 0
	v.PreRelease = ""
	v.Build = ""
}

func (v *Version) ReleaseDraft() {
	v.PreRelease = ""
}

// InitialVersion returns the default initial version for new definitions.
func InitialVersion() string {
	return "0.1.0"
}
