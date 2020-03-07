package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
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
	exitIfError(e)

	currentDir, e := utils.BinAbsDir(executablePath)
	exitIfError(e)

	c := &configurator.Configurator{FilePath: "", CurrentDir: currentDir, Version: version}
	config, e := c.GetConfig()
	exitIfError(e)
	finder := &interpreterfinder.InterpreterFinder{CurrentDir: currentDir}
	mn := &manager.Manager{Config: config, InterpreterFinder: finder}

	app := app.New()
	app.SetIcon(insteadManIcon(c))
	// app.Settings().SetTheme(theme.LightTheme())

	w := newMainWin(app, mn, c)
	w.SetMaster()
	w.ShowAndRun()
}

func newMainWin(app fyne.App, mn *manager.Manager, c *configurator.Configurator) fyne.Window {
	var sw fyne.Window = nil
	var settingsScreen *screens.SettingsScreen = nil

	w := app.NewWindow("InsteadMan")
	entry := widget.NewEntry()
	w.SetContent(widget.NewVBox(
		entry,
		widget.NewLabel(mn.Config.CalculatedGamesPath),
		widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
			if sw == nil {
				sw, settingsScreen = newSettingsWin(app, mn, c, version)
				sw.SetOnClosed(func() {
					w.RequestFocus()
					sw = nil
				})
			}
			settingsScreen.SetMainTab()
			sw.Show()
		}),
		widget.NewButtonWithIcon("About", theme.InfoIcon(), func() {
			if sw == nil {
				sw, settingsScreen = newSettingsWin(app, mn, c, version)
				sw.SetOnClosed(func() {
					w.RequestFocus()
					sw = nil
				})
			}
			settingsScreen.SetAboutTab()
			sw.Show()
		}),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))
	w.CenterOnScreen()

	return w
}

func newSettingsWin(app fyne.App, mn *manager.Manager, c *configurator.Configurator, version string) (fyne.Window, *screens.SettingsScreen) {
	w := app.NewWindow("Settings")
	settingsScreen := screens.NewSettingsScreen(mn, c, insteadManIcon(c), version, w)
	w.SetContent(settingsScreen.Screen)
	w.CenterOnScreen()

	return w, settingsScreen
}

func exitIfError(e error) {
	if e == nil {
		return
	}

	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}

func insteadManIcon(configurator *configurator.Configurator) fyne.Resource {
	iconFile, e := os.Open(configurator.DataResourcePath("../resources/images/logo.png"))
	exitIfError(e)

	r := bufio.NewReader(iconFile)

	b, e := ioutil.ReadAll(r)
	exitIfError(e)

	return fyne.NewStaticResource("insteadman", b)
}
