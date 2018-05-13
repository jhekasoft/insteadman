package utils

import (
	"os"
	"path/filepath"
)

func BinAbsDir(executablePath string) (path string, e error) {
	exePath, e := filepath.Abs(filepath.Dir(executablePath))
	if e != nil {
		return
	}

	path, e = filepath.EvalSymlinks(exePath)

	return
}

func ExistsString(stack []string, element string) bool {
	for _, el := range stack {
		if el == element {
			return true
		}
	}
	return false
}

func PathExist(path string) bool {
	_, e := os.Stat(path)
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return true
	}

	return false
}
