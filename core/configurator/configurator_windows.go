// +build windows

package configurator

import (
	"os"
	"path/filepath"
)

func insteadDir() string {
	homeInsteadDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "instead")
	os.MkdirAll(homeInsteadDir, os.ModePerm)

	return homeInsteadDir
}
