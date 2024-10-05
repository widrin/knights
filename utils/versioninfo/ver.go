package versioninfo

import (
	"fmt"
	"os"
	"strings"
)

var ver = ""

func init() {
	if len(os.Args) < 2 {
		return
	}
	if os.Args[1] == "-version" || os.Args[1] == "-v" {
		fmt.Println(GetVer())
		os.Exit(0)
		return
	}
	if os.Args[1] == "-info" {
		fmt.Println("Ver:", GetVer())
		fmt.Println("Version:", Version)
		fmt.Println("Revision:", Revision)
		fmt.Println("DirtyBuild:", DirtyBuild)
		fmt.Println("LastCommit:", LastCommit)
		os.Exit(0)
		return
	}
}
func GetVer() string {
	if ver == "" {
		if Version != "unknown" && Version != "(devel)" {
			//v0.0.0-20230619063042-e3565d42175d
			str := strings.Split(Version, "-")
			if len(str) == 3 {
				ver = str[1]
				return ver
			}
		}
		if Revision != "unknown" && Revision != "" {
			ver = LastCommit.Format("20060102150405")
			return ver
		}
		ver = "unknown"
		return ver
	}
	return ver
}
func GetStartTime() int64 {
	return startTime
}
