package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetCommand(t *testing.T) {
	commands := map[string]string{
		"list":             "list",
		"search cat":       "search",
		"list --installed": "list",
	}

	for argsStr, mustBeCommand := range commands {
		command := GetCommand(strings.Split(argsStr, " "))
		assert.Equal(t, command, mustBeCommand)
	}
}

func TestGetCommandArg(t *testing.T) {
	commands := map[string]string{
		"search cat":                       "cat",
		"search dog --installed --lang=en": "dog",
		"run лифтер":                       "лифтер",
	}

	for argsStr, mustBeArg := range commands {
		arg := GetCommandArg(strings.Split(argsStr, " "))
		assert.Equal(t, *arg, mustBeArg)
	}

	nilArg := GetCommandArg(strings.Split("list", " "))
	assert.Nil(t, nilArg)
}

func TestFindBoolArg(t *testing.T) {
	commands := map[string]bool{
		"search cat":                       false,
		"search dog --installed --lang=en": true,
		"list --lang=en --installed":       true,
	}

	for argsStr, mustBeVal := range commands {
		argVal := FindBoolArg("--installed", strings.Split(argsStr, " "))
		assert.Equal(t, argVal, mustBeVal)
	}
}

func TestFindStringlArg(t *testing.T) {
	commands := map[string]string{
		"search cat --lang=en":                                   "en",
		"search dog --installed --lang=ru --repository=official": "ru",
		"list --lang=en --installed":                             "en",
	}

	for argsStr, mustBeVal := range commands {
		argVal := FindStringArg("--lang", strings.Split(argsStr, " "))
		assert.Equal(t, *argVal, mustBeVal)
	}

	nilVal := FindStringArg("--lang", strings.Split("list", " "))
	assert.Nil(t, nilVal)

	nilVal2 := FindStringArg("--lang", strings.Split("list --repository=official", " "))
	assert.Nil(t, nilVal2)
}
