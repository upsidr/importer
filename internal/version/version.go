package version

import (
	"embed"
	"fmt"
	"strings"
)

var (
	//go:embed *.txt
	versionData embed.FS

	unknownVersion = &Version{
		Version: "(version unknown)",
	}
)

type Version struct {
	Version          string
	Revision         string
	ReleaseCandidate string
}

func GetVersion() *Version {
	base, err := versionData.ReadFile("VERSION.txt")
	if err != nil {
		return unknownVersion
	}
	v := &Version{
		Version: string(base),
	}

	if rev, err := versionData.ReadFile("REVISION.txt"); err == nil {
		v.Revision = strings.TrimSpace(string(rev))
	}
	if rc, err := versionData.ReadFile("RC.txt"); err == nil {
		v.ReleaseCandidate = strings.TrimSpace(string(rc))
	}

	return v
}

func (v *Version) VersionInfo() string {
	result := v.Version
	if v.ReleaseCandidate != "" {
		result = fmt.Sprintf("%s-%s", v.Version, v.ReleaseCandidate)
	}

	if v.Revision != "" {
		result = result + " ( " + v.Revision + " )"
	}

	return result
}
