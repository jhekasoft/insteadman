// +build darwin

package interpreterFinder

import (
	"os/exec"
	"strings"
)

//const builtInPath = "../../MacOS/sdl-instead"

func exactFilePaths() []string {
	// Add /Application path
	// Can be installed by: "brew install caskroom/cask/instead"
	paths := []string{"/Applications/Instead.app/Contents/MacOS/sdl-instead"}

	// Add command line path
	// Can be installed by: "brew search instead"
	interpreterCommand := "instead"
	out, e := exec.Command("which", interpreterCommand).Output()

	if e != nil {
		return paths
	}
	commandPath := strings.Replace(string(out), "\n", "", -1)
	paths = append(paths, commandPath)


	return paths
}
