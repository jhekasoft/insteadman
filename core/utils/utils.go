package utils

import (
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
