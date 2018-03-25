package main

import (
	"os"
	"fmt"
	"strings"
)

func GetCommand(argsWithoutProg []string) string {
	if len(argsWithoutProg) > 0 {
		return argsWithoutProg[0]
	}

	return ""
}

func GetCommandArg(argsWithoutProg []string) *string {
	if len(argsWithoutProg) > 1 {
		return &argsWithoutProg[1]
	}

	return nil
}

func FindBoolArg(name string, args []string) bool {
	for _, arg := range args {
		if arg == name {
			return true
		}
	}

	return false
}

func FindStringArg(name string, args []string) *string {
	for _, arg := range args {
		searchArgPrefix := name + "="
		if strings.HasPrefix(arg, searchArgPrefix) {
			value := strings.TrimPrefix(arg, searchArgPrefix)
			if value == "" {
				return nil
			}

			return &value
		}
	}

	return nil
}

func ExitIfError(e error) {
	if e == nil {
		return
	}

	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}