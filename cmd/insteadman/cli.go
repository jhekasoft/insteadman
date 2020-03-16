package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
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

func FmtTitle(name string) string {
	return name
}

func FmtName(name string) string {
	return color.New(color.Bold).Sprint(name)
}

func FmtVersion(name string) string {
	return color.RedString(name)
}

func FmtSize(val string) string {
	return color.MagentaString(val)
}

func FmtRepo(name string) string {
	return color.YellowString(name)
}

func FmtInstalled(name string) string {
	return color.GreenString(name)
}

func FmtLang(name string) string {
	return color.BlueString(name)
}

func FmtURL(name string) string {
	return color.CyanString(name)
}

func ExitIfError(e error) {
	if e == nil {
		return
	}

	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}
