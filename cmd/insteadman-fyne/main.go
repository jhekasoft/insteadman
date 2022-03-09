package main

import (
	"errors"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"

	"github.com/jhekasoft/insteadman3/cmd/insteadman-fyne/data"
	"github.com/jhekasoft/insteadman3/cmd/insteadman-fyne/screen"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/interpreterfinder"
	"github.com/jhekasoft/insteadman3/core/manager"
	"github.com/jhekasoft/insteadman3/core/utils"
)

func main() {
	executablePath, e := os.Executable()
	exitIfError(e)

	currentDir, e := utils.BinAbsDir(executablePath)
	exitIfError(e)

	c := &configurator.Configurator{FilePath: "", CurrentDir: currentDir, Version: manager.Version}
	config, e := c.GetConfig()
	exitIfError(e)
	finder := &interpreterfinder.InterpreterFinder{CurrentDir: currentDir}
	mn := &manager.Manager{Config: config, InterpreterFinder: finder}

	app := app.NewWithID("insteadman3-fyne")
	app.SetIcon(data.InsteadManLogo)
	// app.Settings().SetTheme(theme.LightTheme())

	w := newMainWin(app, mn, c)
	w.SetMaster()

	// log.Println(w.Canvas().Scale())

	if mn.InterpreterCommand() == "" {
		findInterpreter(mn, c, w)
	}

	w.ShowAndRun()
}

func newMainWin(app fyne.App, m *manager.Manager, c *configurator.Configurator) fyne.Window {
	w := app.NewWindow("InsteadMan")
	mainScreen := screen.NewMainScreen(w, m, c,
		func() {
			sw, settingsScreen := newSettingsWin(app, m, c)
			settingsScreen.SetMainTab()
			sw.Show()
		},
		func() {
			sw, settingsScreen := newSettingsWin(app, m, c)
			settingsScreen.SetAboutTab()
			sw.Show()
		},
	)
	w.SetContent(mainScreen.Screen)
	w.Resize(fyne.NewSize(800, 500))
	w.CenterOnScreen()

	return w
}

func newSettingsWin(app fyne.App, m *manager.Manager, c *configurator.Configurator) (
	fyne.Window, *screen.SettingsScreen) {
	w := app.NewWindow("Settings")
	settingsScreen := screen.NewSettingsScreen(w, m, c)
	w.SetContent(settingsScreen.Screen)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(600, 400))

	return w, settingsScreen
}

func exitIfError(e error) {
	if e == nil {
		return
	}

	log.Printf("Error: %v\n", e)
	os.Exit(1)
}

func findInterpreter(m *manager.Manager, c *configurator.Configurator, w fyne.Window) {
	path := m.InterpreterFinder.Find()

	if path == nil {
		e := errors.New("INSTEAD has not found. Please add INSTEAD in the Settings")
		dialog.ShowError(e, w)
		return
	}

	log.Printf("INSTEAD has found: %s", *path)

	m.Config.InterpreterCommand = *path
	e := c.SaveConfig(m.Config)
	if e != nil {
		dialog.ShowError(e, w)
		return
	}

	log.Print("Path has saved")
}
