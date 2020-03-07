package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne/theme"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/interpreterfinder"
	"github.com/jhekasoft/insteadman3/core/manager"
	"github.com/jhekasoft/insteadman3/core/utils"
	"github.com/jhekasoft/insteadman3/fyne/screens"
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
	finder := &interpreterfinder.InterpreterFinder{CurrentDir: currentDir}
	mn := &manager.Manager{Config: config, InterpreterFinder: finder}

	app := app.New()
	app.SetIcon(insteadManIcon(c))

	entry := widget.NewEntry()

	w := app.NewWindow("InsteadMan")
	w.SetContent(widget.NewVBox(
		entry,
		widget.NewLabel(config.CalculatedGamesPath),
		widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
			sw := app.NewWindow("Settings")
			sw.SetContent(screens.SettingsScreen(config, &c, mn, insteadManIcon(c), version))
			sw.Show()
		}),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))
	w.SetMaster()

	entry.SetText(config.InterpreterCommand)

	w.ShowAndRun()
}

func ExitIfError(e error) {
	if e == nil {
		return
	}

	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}

func insteadManIcon(configurator configurator.Configurator) fyne.Resource {
	iconFile, e := os.Open(configurator.DataResourcePath("../resources/images/logo.png"))
	ExitIfError(e)

	r := bufio.NewReader(iconFile)

	b, e := ioutil.ReadAll(r)
	ExitIfError(e)

	return fyne.NewStaticResource("insteadman", b)
}
