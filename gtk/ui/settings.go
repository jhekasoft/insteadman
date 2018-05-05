package ui

import (
	"../../core/configurator"
	"../../core/interpreter_finder"
	"../../core/manager"
	gtkutils "../utils"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

const (
	settingsFormFilePath = "resources/gtk/settings.glade"
	aboutTabNum          = 2
)

var (
	SettingsWin *SettingsWindow
)

func GetSettings(manager *manager.Manager, configurator *configurator.Configurator, version string) *SettingsWindow {
	if SettingsWin != nil && SettingsWin.Window.IsVisible() {
		return SettingsWin
	}
	return SettingsWindowNew(manager, configurator, version)
}

func ShowSettingWin(manager *manager.Manager, configurator *configurator.Configurator, version string) {
	SettingsWin = GetSettings(manager, configurator, version)
	SettingsWin.Window.Show()
	SettingsWin.Window.Present()
}

func ShowAboutWin(manager *manager.Manager, configurator *configurator.Configurator, version string) {
	SettingsWin = GetSettings(manager, configurator, version)
	SettingsWin.Window.Show()
	SettingsWin.Window.Present()

	SettingsWin.NtbkCategories.SetCurrentPage(aboutTabNum)
}

type SettingsWindow struct {
	Window *gtk.Window
	e      error

	NtbkCategories *gtk.Notebook

	LblVersion *gtk.Label

	EntryInstead         *gtk.Entry
	BtnInsteadBrowse     *gtk.Button
	TglBtnInsteadBuiltin *gtk.ToggleButton
	BtnInsteadDetect     *gtk.Button
	BtnInsteadCheck      *gtk.Button
	LblInsteadInf        *gtk.Label

	BtnCacheClear *gtk.Button
	LblCacheInf   *gtk.Label

	LblConfigPath *gtk.Label

	Manager      *manager.Manager
	Configurator *configurator.Configurator
}

func SettingsWindowNew(manager *manager.Manager, configurator *configurator.Configurator, version string) *SettingsWindow {
	win := new(SettingsWindow)

	b, e := gtk.BuilderNew()
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	e = b.AddFromFile(configurator.ShareResourcePath(settingsFormFilePath))
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	obj, e := b.GetObject("window_settings")
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	var ok bool
	win.Window, ok = obj.(*gtk.Window)
	if !ok {
		ShowErrorDlgFatal("No settings window")
	}

	win.Manager = manager
	win.Configurator = configurator

	win.NtbkCategories = gtkutils.GetNotebook(b, "notebook_categories")

	win.LblVersion = gtkutils.GetLabel(b, "label_version")
	win.LblVersion.SetText(version)

	win.EntryInstead = gtkutils.GetEntry(b, "entry_instead")
	win.BtnInsteadBrowse = gtkutils.GetButton(b, "button_instead_browse")
	win.TglBtnInsteadBuiltin = gtkutils.GetToggleButton(b, "togglebutton_instead_builtin")
	win.BtnInsteadDetect = gtkutils.GetButton(b, "button_instead_detect")
	win.BtnInsteadCheck = gtkutils.GetButton(b, "button_instead_check")
	win.LblInsteadInf = gtkutils.GetLabel(b, "label_instead_inf")

	win.BtnCacheClear = gtkutils.GetButton(b, "button_cache_clear")
	win.LblCacheInf = gtkutils.GetLabel(b, "label_cache_inf")

	win.LblConfigPath = gtkutils.GetLabel(b, "label_config_path")

	win.readSettings()

	win.Window.SetTitle("Settings")
	win.Window.SetPosition(gtk.WIN_POS_CENTER)

	return win
}

func (win *SettingsWindow) readSettings() {
	finder := interpreterFinder.InterpreterFinder{Config: win.Manager.Config}
	config := win.Manager.Config

	// INSTEAD
	win.EntryInstead.SetText(config.InterpreterCommand)
	win.TglBtnInsteadBuiltin.SetSensitive(finder.HaveBuiltIn())
	win.TglBtnInsteadBuiltin.SetActive(config.UseBuiltinInterpreter)

	// Cache
	win.BtnCacheClear.SetTooltipText("Cache directory: " + win.Manager.CacheDir())

	// Config path
	win.LblConfigPath.SetText(win.Configurator.FilePath)
}
