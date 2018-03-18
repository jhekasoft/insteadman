// +build darwin

package interpreterFinder

import (
	"os"
)

func Find() *string {
	// builtInPath := "../../MacOS/sdl-instead"
	exactFilePaths := []string{"/Applications/Instead.app/Contents/MacOS/sdl-instead"}
	for _, path := range exactFilePaths {
		_, e := os.Stat(path)
		exists := !os.IsNotExist(e)

		if exists && e == nil {
			return &path
		}
	}

	return nil
}
