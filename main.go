package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/interpreterfinder"
	"github.com/jhekasoft/insteadman3/core/manager"
	"github.com/jhekasoft/insteadman3/core/utils"
	"github.com/wailsapp/wails"
)

func basic() string {
	return "World!"
}

func games() string {
	res, err := json.Marshal(gamesTitles)
	log.Println(string(res))
	log.Println(err)
	return string(res)
}

//go:embed frontend/build/static/js/main.js
var js string

//go:embed frontend/build/static/css/main.css
var css string

var gamesTitles []string

func main() {
	m, _ := initManagerAndConfigurator()
	gameItems, _ := m.GetSortedGamesByDateDesc()
	for _, v := range gameItems {
		gamesTitles = append(gamesTitles, v.Title)
	}

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "insteadman-wails",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(basic)
	app.Bind(games)
	app.Run()
}

func initManagerAndConfigurator() (*manager.Manager, *configurator.Configurator) {
	executablePath, e := os.Executable()
	ExitIfError(e)

	currentDir, e := utils.BinAbsDir(executablePath)
	ExitIfError(e)

	c := configurator.Configurator{FilePath: "", CurrentDir: currentDir, Version: manager.Version}
	config, e := c.GetConfig()
	ExitIfError(e)

	finder := &interpreterfinder.InterpreterFinder{CurrentDir: currentDir}

	m := manager.Manager{Config: config, InterpreterFinder: finder}

	return &m, &c
}

func ExitIfError(e error) {
	if e == nil {
		return
	}

	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}
