package interpreterfinder

import (
	"../configurator"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// InterpreterFinder provides finding INSTEAD interpreter in the filesystem
type InterpreterFinder struct {
	// CurrentDir is a currently directory of the executable file
	CurrentDir string
}

// HaveBuiltIn checks is there is built-in INSTEAD with InsteadMan
func (f *InterpreterFinder) HaveBuiltIn() bool {
	_, e := os.Stat(filepath.Join(f.CurrentDir, builtinRelativeFilePath))
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return true
	}

	return false
}

// FindBuiltIn returns built-in INSTEAD interpreter path
func (f *InterpreterFinder) FindBuiltIn() (path string) {
	if f.HaveBuiltIn() {
		path = filepath.Join(f.CurrentDir, builtinRelativeFilePath)
	}
	return
}

// Find finds INSTEAD interpreter in the filesystem
func (f *InterpreterFinder) Find() *string {
	// External interpreter
	for _, path := range exactFilePaths() {
		_, e := os.Stat(path)
		exists := !os.IsNotExist(e)

		if exists && e == nil {
			return &path
		}
	}

	return nil
}

// Check checks the INSTEAD interpreter and returns version of INSTEAS
// If INSTEAD could not be found returns error
func (f *InterpreterFinder) Check(command string) (version string, e error) {
	out, e := exec.Command(configurator.ExpandInterpreterCommand(command), "-version").Output()
	if e != nil {
		return "", e
	}

	replacer := strings.NewReplacer("\n", "", "\r", "")
	version = replacer.Replace(string(out))

	return
}
