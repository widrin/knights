package versioninfo

import "runtime/debug"

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	Version = info.Main.Version
}
