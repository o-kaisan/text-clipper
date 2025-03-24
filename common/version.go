package common

import (
	"runtime/debug"
)

func GetVersionFromGit() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		return bi.Main.Version
	}
	return "Unknown"
}
