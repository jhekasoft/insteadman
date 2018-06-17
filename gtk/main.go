package main

import (
	"../core/configurator"
	"../core/interpreterfinder"
	"../core/manager"
	"../core/utils"
	"./i18n"
	"./osintegration"
	"./ui"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"runtime"
)

const (
	title = "InsteadMan"

	envDataPath   = "DATA_PATH"
	envLocalePath = "LOCALE_PATH"

	i18nDomain = "insteadman"
)

var (
	version = "3"
)

func main() {
	runtime.LockOSThread()

	// OS integrations
	osintegration.OsIntegrate()

	gtk.Init(nil)

	executablePath, e := os.Executable()
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), nil)
	}

	currentDir, e := utils.BinAbsDir(executablePath)
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), nil)
	}

	dataPath := os.Getenv(envDataPath)
	localePath := os.Getenv(envLocalePath)

	cf := &configurator.Configurator{FilePath: "", CurrentDir: currentDir, DataPath: dataPath,
		LocalePath: localePath, Version: version}

	config, e := cf.GetConfig()
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), nil)
	}

	finder := &interpreterfinder.InterpreterFinder{CurrentDir: currentDir}

	mn := &manager.Manager{Config: config, InterpreterFinder: finder}

	// I18n init
	i18n.Init(cf.DataLocalePath(), i18nDomain, config.Lang)

	mainWindow := ui.GetMain(mn, cf, title, version)

	if mn.InterpreterCommand() == "" {
		findInterpreter(mn, cf, mainWindow.Window)
	}

	ui.ShowExistingMainWindow(true)

	gtk.Main()
}

func findInterpreter(m *manager.Manager, c *configurator.Configurator, wnd *gtk.Window) {
	path := m.InterpreterFinder.Find()

	if path == nil {
		ui.ShowErrorDlg(i18n.T("INSTEAD has not found. Please add INSTEAD in the Settings."), wnd)
		return
	}

	log.Printf("INSTEAD has found: %s", *path)

	m.Config.InterpreterCommand = *path
	e := c.SaveConfig(m.Config)
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), wnd)
		return
	}

	log.Print("Path has saved")
}
