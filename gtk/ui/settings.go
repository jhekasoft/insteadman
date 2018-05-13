package ui

import (
	"../../core/configurator"
	"../../core/manager"
	gtkutils "../utils"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
)

const (
	settingsFormFilePath = "resources/gtk/settings.glade"
	aboutTabNum          = 2
	RepositoryColumnName = 0
	RepositoryColumnUrl  = 1
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

func ShowSettingWin(manager *manager.Manager, configurator *configurator.Configurator, version string,
	parent *gtk.Window) {

	SettingsWin = GetSettings(manager, configurator, version)
	if parent != nil {
		SettingsWin.Window.SetTransientFor(parent)
	}
	SettingsWin.Window.Show()
	SettingsWin.Window.Present()
}

func ShowAboutWin(manager *manager.Manager, configurator *configurator.Configurator, version string,
	parent *gtk.Window) {

	SettingsWin = GetSettings(manager, configurator, version)
	if parent != nil {
		SettingsWin.Window.SetTransientFor(parent)
	}
	SettingsWin.Window.Show()
	SettingsWin.Window.Present()

	SettingsWin.NtbkCategories.SetCurrentPage(aboutTabNum)
}

type SettingsWindow struct {
	Window *gtk.Window
	e      error

	NtbkCategories *gtk.Notebook

	EntryInstead         *gtk.Entry
	BtnInsteadBrowse     *gtk.Button
	TglBtnInsteadBuiltin *gtk.ToggleButton
	BtnInsteadDetect     *gtk.Button
	BtnInsteadCheck      *gtk.Button
	LblInsteadInf        *gtk.Label

	BtnCacheClear *gtk.Button
	LblCacheInf   *gtk.Label

	LblConfigPath *gtk.Label

	LblVersion *gtk.Label

	ListStoreRepositories   *gtk.ListStore
	TrSlctnRepositories     *gtk.TreeSelection
	CllRndrTxtName          *gtk.CellRendererText
	CllRndrTxtUrl           *gtk.CellRendererText
	BtnRepositoriesAdd      *gtk.Button
	BtnRepositoriesRemove   *gtk.Button
	BtnRepositoriesUp       *gtk.Button
	BtnRepositoriesDown     *gtk.Button
	BtnRepositoriesDefaults *gtk.Button

	BtnClose *gtk.Button

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

	// Main tab
	win.EntryInstead = gtkutils.GetEntry(b, "entry_instead")
	win.BtnInsteadBrowse = gtkutils.GetButton(b, "button_instead_browse")
	win.TglBtnInsteadBuiltin = gtkutils.GetToggleButton(b, "togglebutton_instead_builtin")
	win.BtnInsteadDetect = gtkutils.GetButton(b, "button_instead_detect")
	win.BtnInsteadCheck = gtkutils.GetButton(b, "button_instead_check")
	win.LblInsteadInf = gtkutils.GetLabel(b, "label_instead_inf")

	win.BtnCacheClear = gtkutils.GetButton(b, "button_cache_clear")
	win.LblCacheInf = gtkutils.GetLabel(b, "label_cache_inf")

	win.LblConfigPath = gtkutils.GetLabel(b, "label_config_path")

	// Repositories tab
	win.ListStoreRepositories = gtkutils.GetListStore(b, "liststore_repositories")
	win.CllRndrTxtName = gtkutils.GetCellRendererText(b, "cellrenderertext_repositories_name")
	win.CllRndrTxtUrl = gtkutils.GetCellRendererText(b, "cellrenderertext_repositories_url")
	win.BtnRepositoriesAdd = gtkutils.GetButton(b, "button_repositories_add")
	win.BtnRepositoriesRemove = gtkutils.GetButton(b, "button_repositories_remove")
	win.BtnRepositoriesUp = gtkutils.GetButton(b, "button_repositories_up")
	win.BtnRepositoriesDown = gtkutils.GetButton(b, "button_repositories_down")
	win.BtnRepositoriesDefaults = gtkutils.GetButton(b, "button_repositories_defaults")
	treeViewRepositories := gtkutils.GetTreeView(b, "treeview_repositories")
	win.TrSlctnRepositories, e = treeViewRepositories.GetSelection()
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	// About tab
	win.LblVersion = gtkutils.GetLabel(b, "label_version")
	win.LblVersion.SetText(version)

	win.BtnClose = gtkutils.GetButton(b, "button_close")

	win.readSettings()

	// Handlers
	handlers := &SettingsWindowHandlers{win: win}
	win.EntryInstead.Connect("changed", handlers.insteadChanged)
	win.BtnInsteadBrowse.Connect("clicked", handlers.insteadBrowseClicked)
	win.TglBtnInsteadBuiltin.Connect("clicked", handlers.insteadBuiltinClicked)
	win.BtnInsteadDetect.Connect("clicked", handlers.insteadDetectClicked)
	win.BtnInsteadCheck.Connect("clicked", handlers.insteadCheckClicked)
	win.BtnCacheClear.Connect("clicked", handlers.cacheClearClicked)
	//win.TrSlctnRepositories.Connect("changed", handlers.repositoriesChanged)
	win.CllRndrTxtName.Connect("edited", handlers.repositoriesNameEdited)
	win.CllRndrTxtUrl.Connect("edited", handlers.repositoriesUrlEdited)
	win.BtnRepositoriesAdd.Connect("clicked", handlers.repositoryAddClicked)
	win.BtnRepositoriesRemove.Connect("clicked", handlers.repositoryRemoveClicked)
	win.BtnRepositoriesUp.Connect("clicked", handlers.repositoryUpClicked)
	win.BtnRepositoriesDown.Connect("clicked", handlers.repositoryDownClicked)
	win.BtnRepositoriesDefaults.Connect("clicked", handlers.repositoryDefaultsClicked)
	win.BtnClose.Connect("clicked", handlers.closeClicked)
	win.Window.Connect("delete_event", handlers.settingsDeleted)

	win.Window.SetTitle("Settings")

	return win
}

func (win *SettingsWindow) readSettings() {
	config := win.Manager.Config

	// INSTEAD
	win.EntryInstead.SetText(config.InterpreterCommand)
	haveBuiltInInstead := win.Manager.InterpreterFinder.HaveBuiltIn()
	win.TglBtnInsteadBuiltin.SetSensitive(haveBuiltInInstead)
	if !haveBuiltInInstead {
		win.TglBtnInsteadBuiltin.SetTooltipText("Built-in INSTEAD hasn't found")
	}
	win.TglBtnInsteadBuiltin.SetActive(config.UseBuiltinInterpreter)
	win.toggleBuiltin(!config.UseBuiltinInterpreter || !win.Manager.InterpreterFinder.HaveBuiltIn())

	// Cache
	win.BtnCacheClear.SetTooltipText("Cache directory: " + win.Manager.CacheDir())

	// Config path
	win.LblConfigPath.SetText(win.Configurator.FilePath)

	// Repositories
	win.ListStoreRepositories.Clear()
	for _, repo := range win.Manager.Config.Repositories {
		addToListStoreRepositories(win.ListStoreRepositories, repo.Name, repo.Url)
	}
}

func (win *SettingsWindow) setRepositoriesConfigFromListStore() {
	repos, e := configRepositoriesFromListStore(win.ListStoreRepositories)
	if e != nil {
		ShowErrorDlg(e.Error())
		return
	}
	win.Manager.Config.Repositories = repos
}

func addToListStoreRepositories(ls *gtk.ListStore, name, url string) (iter *gtk.TreeIter) {
	iter = new(gtk.TreeIter)
	ls.InsertWithValues(iter, -1, []int{RepositoryColumnName, RepositoryColumnUrl}, []interface{}{name, url})
	return iter
}

func (win *SettingsWindow) toggleBuiltin(active bool) {
	win.EntryInstead.SetSensitive(active)
	win.BtnInsteadBrowse.SetSensitive(active)
	win.BtnInsteadDetect.SetSensitive(active)
}

func configRepositoriesFromListStore(ls *gtk.ListStore) (repositories []configurator.Repository, e error) {
	iter, _ := ls.GetIterFirst()

	var (
		value     *glib.Value
		name, url string
	)

	// Collect repositories from list store
	for iter != nil {
		value, e = ls.GetValue(iter, RepositoryColumnName)
		if e != nil {
			return
		}
		name, e = value.GetString()
		if e != nil {
			return
		}

		value, e = ls.GetValue(iter, RepositoryColumnUrl)
		if e != nil {
			return
		}
		url, e = value.GetString()
		if e != nil {
			return
		}

		// Add non-empty repositories
		if name != "" && url != "" {
			repositories = append(repositories, configurator.Repository{Name: name, Url: url})
		}

		if !ls.IterNext(iter) {
			iter = nil
		}
	}

	return
}

type SettingsWindowHandlers struct {
	win *SettingsWindow
}

func (h *SettingsWindowHandlers) insteadChanged(s *gtk.Entry) {
	value, e := s.GetText()
	if e == nil {
		h.win.Manager.Config.InterpreterCommand = value
	}
}

func (h *SettingsWindowHandlers) insteadBrowseClicked(s *gtk.Button) {
	s.SetSensitive(false)

	dlg, _ := gtk.FileChooserDialogNewWith2Buttons("Choose INSTEAD", h.win.Window, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL, "Open", gtk.RESPONSE_ACCEPT)

	response := dlg.Run()
	if response == int(gtk.RESPONSE_ACCEPT) {
		h.win.EntryInstead.SetText(dlg.GetFilename())
	}

	dlg.Destroy()

	s.SetSensitive(true)
}

func (h *SettingsWindowHandlers) insteadBuiltinClicked(s *gtk.ToggleButton) {
	h.win.Manager.Config.UseBuiltinInterpreter = s.GetActive()
	h.win.readSettings()
}

func (h *SettingsWindowHandlers) insteadDetectClicked(s *gtk.Button) {
	s.SetSensitive(false)
	h.win.LblInsteadInf.Hide()

	go func() {
		command := h.win.Manager.InterpreterFinder.Find()

		_, e := glib.IdleAdd(func() {
			if command != nil {
				h.win.EntryInstead.SetText(*command)
				h.win.LblInsteadInf.SetText("INSTEAD has detected!")
				h.win.LblInsteadInf.Show()
			} else {
				h.win.LblInsteadInf.SetText("INSTEAD hasn't detected!")
				h.win.LblInsteadInf.Show()
			}

			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("INSTEAD detect. IdleAdd() failed:", e)
		}
	}()
}

func (h *SettingsWindowHandlers) insteadCheckClicked(s *gtk.Button) {
	s.SetSensitive(false)
	h.win.LblInsteadInf.Hide()

	go func() {
		version, checkErr := h.win.Manager.InterpreterFinder.Check(h.win.Manager.InterpreterCommand())

		_, e := glib.IdleAdd(func() {
			var txt string
			if checkErr != nil {
				if h.win.Manager.IsBuiltinInterpreterCommand() {
					txt = "INSTEAD built-in check failed!"
				} else {
					txt = "INSTEAD check failed!"
				}

			} else {
				txt = "INSTEAD " + version + " has found!"
			}
			h.win.LblInsteadInf.SetText(txt)

			h.win.LblInsteadInf.Show()
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("INSTEAD check. IdleAdd() failed:", e)
		}
	}()
}

func (h *SettingsWindowHandlers) cacheClearClicked(s *gtk.Button) {
	s.SetSensitive(false)
	h.win.LblCacheInf.Hide()

	go func() {
		cacheErr := h.win.Manager.ClearCache()
		_, e := glib.IdleAdd(func() {
			if cacheErr != nil {
				ShowErrorDlg(cacheErr.Error())
			} else {
				h.win.LblCacheInf.SetText("Cache has been cleared!")
			}

			h.win.LblCacheInf.Show()
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("Cache clear. IdleAdd() failed:", e)
		}
	}()
}

//func (h *SettingsWindowHandlers) repositoriesChanged(s *gtk.TreeSelection) {
//}

func (h *SettingsWindowHandlers) repositoriesNameEdited(s *gtk.CellRendererText, path, newText string) {
	iter, e := gtkutils.GetIterFromTextPathInListStore(h.win.ListStoreRepositories, path)
	if e != nil {
		ShowErrorDlg(e.Error())
		return
	}

	h.win.ListStoreRepositories.SetValue(iter, RepositoryColumnName, newText)

	h.win.setRepositoriesConfigFromListStore()
}

func (h *SettingsWindowHandlers) repositoriesUrlEdited(s *gtk.CellRendererText, path, newText string) {
	iter, e := gtkutils.GetIterFromTextPathInListStore(h.win.ListStoreRepositories, path)
	if e != nil {
		ShowErrorDlg(e.Error())
		return
	}

	h.win.ListStoreRepositories.SetValue(iter, RepositoryColumnUrl, newText)

	h.win.setRepositoriesConfigFromListStore()
}

func (h *SettingsWindowHandlers) repositoryAddClicked() {
	iter := addToListStoreRepositories(h.win.ListStoreRepositories, "", "")
	h.win.TrSlctnRepositories.SelectIter(iter)

	h.win.setRepositoriesConfigFromListStore()
}

func (h *SettingsWindowHandlers) repositoryRemoveClicked() {
	iter, e := gtkutils.FindFirstIterInTreeSelection(h.win.ListStoreRepositories, h.win.TrSlctnRepositories)
	if e != nil {
		log.Printf("Error: %v", e)
		return
	}
	if iter == nil {
		return
	}

	h.win.ListStoreRepositories.Remove(iter)

	h.win.setRepositoriesConfigFromListStore()
}

func (h *SettingsWindowHandlers) repositoryUpClicked() {
	iter, e := gtkutils.FindFirstIterInTreeSelection(h.win.ListStoreRepositories, h.win.TrSlctnRepositories)
	if e != nil {
		log.Printf("Error: %v", e)
		return
	}
	curIter := *iter

	if h.win.ListStoreRepositories.IterPrevious(iter) {
		prevIter := *iter

		h.win.ListStoreRepositories.MoveBefore(&curIter, &prevIter)
	}

	h.win.setRepositoriesConfigFromListStore()
}

func (h *SettingsWindowHandlers) repositoryDownClicked() {
	iter, e := gtkutils.FindFirstIterInTreeSelection(h.win.ListStoreRepositories, h.win.TrSlctnRepositories)
	if e != nil {
		log.Printf("Error: %v", e)
		return
	}
	curIter := *iter

	if h.win.ListStoreRepositories.IterNext(iter) {
		nextIter := *iter

		h.win.ListStoreRepositories.MoveAfter(&curIter, &nextIter)
	}

	h.win.setRepositoriesConfigFromListStore()
}

func (h *SettingsWindowHandlers) repositoryDefaultsClicked() {
	skeletonConfig, e := h.win.Configurator.GetSkeletonConfig()
	if e != nil {
		ShowErrorDlg(e.Error())
		return
	}

	h.win.Manager.Config.Repositories = skeletonConfig.Repositories

	h.win.readSettings()
}

func (h *SettingsWindowHandlers) closeClicked() {
	h.win.Window.Close()
}

func (h *SettingsWindowHandlers) settingsDeleted() {
	// Auto save
	e := h.win.Configurator.SaveConfig(h.win.Manager.Config)
	if e != nil {
		ShowErrorDlg(e.Error())
		return
	}

	h.win.Window.Destroy()
}
