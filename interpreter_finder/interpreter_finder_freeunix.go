// +build !windows,!darwin

package interpreterFinder

import (
	"os/exec"
	"strings"
)

func exactFilePaths() []string {
	interpreterCommand := "instead"

	out, e := exec.Command("which", interpreterCommand).Output()

	if e != nil {
		return []string{}
	}

	path := strings.Replace(string(out), "\n", "", -1)

	return []string{path}
}
