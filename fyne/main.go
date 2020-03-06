package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/utils"
)

var version = "3"

func main() {
	executablePath, e := os.Executable()
	ExitIfError(e)

	currentDir, e := utils.BinAbsDir(executablePath)
	ExitIfError(e)

	c := configurator.Configurator{FilePath: "", CurrentDir: currentDir, Version: version}
	config, e := c.GetConfig()
	ExitIfError(e)

	fmt.Printf("Config: %v\n", config)

	app := app.New()

	w := app.NewWindow("InsteadMan")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	w.ShowAndRun()
}

func ExitIfError(e error) {
	if e == nil {
		return
	}

	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}
