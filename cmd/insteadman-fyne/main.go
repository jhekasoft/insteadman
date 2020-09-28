package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"

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

	if mn.InterpreterCommand() == "" {
		findInterpreter(mn, c, w)
	}

	w.ShowAndRun()
}

func newMainWin(app fyne.App, mn *manager.Manager, c *configurator.Configurator) fyne.Window {
	var sw fyne.Window
	var settingsScreen *screen.SettingsScreen

	w := app.NewWindow("InsteadMan")
	// TODO: improve settings open functions
	mainScreen := screen.NewMainScreen(
		mn,
		c,
		data.InsteadManLogo,
		w,
		func() {
			if sw == nil {
				sw, settingsScreen = newSettingsWin(app, mn, c)
				sw.SetOnClosed(func() {
					w.RequestFocus()
					sw = nil
				})
			}
			settingsScreen.SetMainTab()
			sw.Show()
		},
		func() {
			if sw == nil {
				sw, settingsScreen = newSettingsWin(app, mn, c)
				sw.SetOnClosed(func() {
					w.RequestFocus()
					sw = nil
				})
			}
			settingsScreen.SetAboutTab()
			sw.Show()
		},
	)
	w.SetContent(mainScreen.Screen)
	w.Resize(fyne.NewSize(800, 500))
	w.CenterOnScreen()

	return w
}

func newSettingsWin(app fyne.App, mn *manager.Manager, c *configurator.Configurator) (fyne.Window, *screen.SettingsScreen) {
	w := app.NewWindow("Settings")
	settingsScreen := screen.NewSettingsScreen(mn, c, data.InsteadManLogo, w)
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

// func insteadManIcon(configurator *configurator.Configurator) fyne.Resource {
// 	iconFile, e := os.Open(configurator.DataResourcePath("../resources/images/logo.png"))
// 	exitIfError(e)

// 	r := bufio.NewReader(iconFile)

// 	b, e := ioutil.ReadAll(r)
// 	exitIfError(e)

// 	return fyne.NewStaticResource("insteadman", b)
// }

func findInterpreter(m *manager.Manager, c *configurator.Configurator, w fyne.Window) {
	path := m.InterpreterFinder.Find()

	if path == nil {
		e := errors.New("INSTEAD has not found. Please add INSTEAD in the Settings.")
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
