// +build linux freebsd netbsd openbsd

package interpreterFinder

import (
	"os/exec"
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
