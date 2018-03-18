package interpreterFinder

import (
	"../configurator"
	"os"
	"os/exec"
	"strings"
)

type InterpreterFinder struct {
	Config *configurator.InsteadmanConfig
}

func (f *InterpreterFinder) Find() *string {
	for _, path := range exactFilePaths {
		_, e := os.Stat(path)
		exists := !os.IsNotExist(e)

		if exists && e == nil {
			return &path
		}
	}

	return nil
}

func (f *InterpreterFinder) CheckInterpreter() (string, error) {
	out, e := exec.Command(f.Config.InterpreterCommand, "-version").Output()
	if e != nil {
		return "", e
	}

	version := strings.Replace(string(out), "\n", "", -1)

	return version, nil
}
