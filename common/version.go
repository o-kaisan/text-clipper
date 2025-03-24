package common

import (
	"log"
	"runtime/debug"
)

func GetVersionFromGit() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		log.Printf(bi.Main.Version)
		return bi.Main.Version
	}
	return "Unknown"
}
