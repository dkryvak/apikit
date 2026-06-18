package meta

import "runtime"

func AppDirName() string {
	if runtime.GOOS == "windows" {
		return AppName
	}
	return "." + AppName
}
