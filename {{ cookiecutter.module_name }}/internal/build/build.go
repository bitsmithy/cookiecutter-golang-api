package build

import (
	"strings"
)

var (
	Version string
	Commit  string
	Time    string
	User    string
)

func isReleaseBuild() bool {
	return len(strings.TrimSpace(Version)) != 0
}

func Info() (ret map[string]string) {
	if isReleaseBuild() {
		ret = map[string]string{
			"version":  Version,
			"commit":   Commit,
			"built_at": Time,
			"built_by": User,
		}
	} else {
		ret = map[string]string{
			"version": "development",
		}
	}

	return
}
