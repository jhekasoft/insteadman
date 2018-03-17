// +build linux freebsd netbsd openbsd darwin

package configurator

import "os"

func userHomeDir() string {
	return os.Getenv("HOME")
}