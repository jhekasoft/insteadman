package interpreterFinder

import (
	// "fmt"
	"../configurator"
	"os/exec"
)

func CheckInterpreter(config *configurator.InsteadmanConfig) (string, error) {
	out, e := exec.Command(config.InterpreterCommand, "-version").Output()
	if e != nil {
		return "", e
	}

	return string(out), nil
}
