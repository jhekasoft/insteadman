package ui

import (
	"../../core/configurator"
	"../../core/interpreter_finder"
	"../../core/manager"
	gtkutils "../utils"
	"github.com/gotk3/gotk3/glib"
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

	BtnClose *gtk.Button

	Manager      *manager.Manager
	Configurator *configurator.Configurator
	Finder       *interpreterFinder.InterpreterFinder
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

	win.BtnClose = gtkutils.GetButton(b, "button_close")

	win.Finder = &interpreterFinder.InterpreterFinder{Config: win.Manager.Config}
	win.readSettings()

	// Handlers
	handlers := &SettingWindowHandlers{win: win}
	win.BtnInsteadDetect.Connect("clicked", handlers.insteadDetectClicked)
	win.BtnInsteadCheck.Connect("clicked", handlers.insteadCheckClicked)
	win.BtnCacheClear.Connect("clicked", handlers.cacheClearClicked)
	win.BtnClose.Connect("clicked", handlers.closeClicked)
	win.Window.Connect("delete_event", handlers.settingsDeleted)

	win.Window.SetTitle("Settings")
	win.Window.SetPosition(gtk.WIN_POS_CENTER)

	return win
}

func (win *SettingsWindow) readSettings() {
	config := win.Manager.Config

	// INSTEAD
	win.EntryInstead.SetText(config.InterpreterCommand)
	win.TglBtnInsteadBuiltin.SetSensitive(win.Finder.HaveBuiltIn())
	win.TglBtnInsteadBuiltin.SetActive(config.UseBuiltinInterpreter)

	// Cache
	win.BtnCacheClear.SetTooltipText("Cache directory: " + win.Manager.CacheDir())

	// Config path
	win.LblConfigPath.SetText(win.Configurator.FilePath)
}

type SettingWindowHandlers struct {
	win *SettingsWindow
}

func (h *SettingWindowHandlers) insteadDetectClicked(s *gtk.Button) {
	s.SetSensitive(false)
	h.win.LblInsteadInf.Hide()

	go func() {
		_, e := glib.IdleAdd(func() {
			command := h.win.Finder.Find()
			if command == nil {
				h.win.LblInsteadInf.SetText("INSTEAD hasn't detected!")
				h.win.LblInsteadInf.Show()
				s.SetSensitive(true)
				return
			}

			h.win.Manager.Config.InterpreterCommand = *command
			e := h.win.Configurator.SaveConfig(h.win.Manager.Config)
			if e != nil {
				ShowErrorDlg(e.Error())
				s.SetSensitive(true)
				return
			}

			h.win.readSettings()
			h.win.LblInsteadInf.SetText("INSTEAD has detected!")
			h.win.LblInsteadInf.Show()
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("INSTEAD detect. IdleAdd() failed:", e)
		}
	}()
}

func (h *SettingWindowHandlers) insteadCheckClicked(s *gtk.Button) {
	s.SetSensitive(false)
	h.win.LblInsteadInf.Hide()

	go func() {
		_, e := glib.IdleAdd(func() {
			command, e := h.win.EntryInstead.GetText()
			if e != nil {
				ShowErrorDlg(e.Error())
				s.SetSensitive(true)
				return
			}

			version, e := h.win.Finder.Check(command)
			if e != nil {
				h.win.LblInsteadInf.SetText("INSTEAD check failed!")
				h.win.LblInsteadInf.Show()
				s.SetSensitive(true)
				return
			}

			h.win.LblInsteadInf.SetText("INSTEAD " + version + " has found!")
			h.win.LblInsteadInf.Show()
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("INSTEAD check. IdleAdd() failed:", e)
		}
	}()
}

func (h *SettingWindowHandlers) cacheClearClicked(s *gtk.Button) {
	s.SetSensitive(false)
	h.win.LblCacheInf.Hide()

	go func() {
		_, e := glib.IdleAdd(func() {
			e := h.win.Manager.ClearCache()
			if e != nil {
				ShowErrorDlg(e.Error())
				s.SetSensitive(true)
				return
			}

			h.win.LblCacheInf.SetText("Cache has been cleared!")
			h.win.LblCacheInf.Show()
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("Cache clear. IdleAdd() failed:", e)
		}
	}()
}

func (h *SettingWindowHandlers) closeClicked() {
	h.win.Window.Close()
}

func (h *SettingWindowHandlers) settingsDeleted() {
	needSave := false

	interpreteNewValue, _ := h.win.EntryInstead.GetText()
	if interpreteNewValue != "" && interpreteNewValue != h.win.Manager.Config.InterpreterCommand {
		h.win.Manager.Config.InterpreterCommand = interpreteNewValue
		needSave = true
	}

	// Autosave
	if needSave {
		e := h.win.Configurator.SaveConfig(h.win.Manager.Config)
		if e != nil {
			ShowErrorDlg(e.Error())
			return
		}
	}
}
