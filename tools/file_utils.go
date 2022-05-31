package tools

import "os"

func IsPathExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	} else {
		return true
	}
}
