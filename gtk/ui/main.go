package ui

import (
	"../../core/configurator"
	"../../core/manager"
	"../../core/utils"
	"../i18n"
	"../os_integration"
	gtkutils "../utils"
	"fmt"
	"github.com/gosexy/gettext"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"log"
	"strings"
)

const (
	mainFormFilePath = "resources/gtk/main.glade"
	logoFilePath     = "resources/images/logo.png"

	comboBoxColumnId    = 0
	comboBoxColumnTitle = 1

	gameColumnId         = 0
	gameColumnTitle      = 1
	gameColumnVersion    = 2
	gameColumnSizeHuman  = 3
	gameColumnFontWeight = 4
	gameColumnSize       = 5

	fontWeightNormal = pango.WEIGHT_NORMAL
	fontWeightBold   = pango.WEIGHT_BOLD
)

var (
	MainWin *MainWindow
)

func ShowMainWindow(manager *manager.Manager, configurator *configurator.Configurator,
	title, version string) *MainWindow {

	GetMain(manager, configurator, title, version)
	return ShowExistingMainWindow(true)
}

func ShowExistingMainWindow(updateRepositories bool) *MainWindow {
	MainWin.Window.Show()
	MainWin.Window.Present()

	if updateRepositories {
		MainWin.updateRepositories()
	}

	return MainWin
}

// Singleton
func GetMain(manager *manager.Manager, configurator *configurator.Configurator,
	title, version string) *MainWindow {

	if MainWin != nil && MainWin.Window.IsVisible() {
		return MainWin
	}
	MainWin = MainWindowNew(manager, configurator, title, version)

	return MainWin
}

type MainWindow struct {
	Window *gtk.Window

	ListStoreGames *gtk.ListStore
	ListStoreRepo  *gtk.ListStore
	ListStoreLang  *gtk.ListStore

	GamesSelection *gtk.TreeSelection

	BtnUpdate        *gtk.Button
	EntryKeyword     *gtk.Entry
	CmbBoxRepo       *gtk.ComboBox
	CmbBoxLang       *gtk.ComboBox
	ChckBtnInstalled *gtk.CheckButton
	BtnClear         *gtk.Button

	ScrWndGames  *gtk.ScrolledWindow
	SpinnerGames *gtk.Spinner

	LblGameTitle   *gtk.Label
	ImgGame        *gtk.Image
	LblGameRepo    *gtk.Label
	LblGameLang    *gtk.Label
	LblGameVersion *gtk.Label
	ScrWndGameDesc *gtk.ScrolledWindow
	LblGameDesc    *gtk.Label
	BtnGameRun     *gtk.Button
	BtnGameInstall *gtk.Button
	BtnGameUpdate  *gtk.Button
	BtnGameRemove  *gtk.Button

	SprtrSideBox *gtk.Separator
	BxSideBox    *gtk.Box

	MenuItmSortingReset *gtk.MenuItem
	ChckMenuItmSideBar  *gtk.CheckMenuItem
	MenuItmSettings     *gtk.MenuItem
	MenuItmAbout        *gtk.MenuItem

	PixBufGameDefaultImage *gdk.Pixbuf
	PixBufGameImage        *gdk.Pixbuf

	Games        []manager.Game
	CurGame      *manager.Game // current selected game
	IsRefreshing bool

	Title   string
	Version string

	Manager      *manager.Manager
	Configurator *configurator.Configurator
}

func MainWindowNew(manager *manager.Manager, configurator *configurator.Configurator,
	title, version string) *MainWindow {

	win := new(MainWindow)

	win.Manager = manager
	win.Configurator = configurator
	win.Title = title
	win.Version = version

	b, e := gtk.BuilderNew()
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	e = b.AddFromFile(configurator.DataResourcePath(mainFormFilePath))
	if e != nil {
		ShowErrorDlgFatal(e.Error(), win.Window)
	}

	obj, e := b.GetObject("window_main")
	if e != nil {
		ShowErrorDlgFatal(e.Error(), win.Window)
	}

	var ok bool
	win.Window, ok = obj.(*gtk.Window)
	if !ok {
		ShowErrorDlgFatal("No main window", win.Window)
	}

	win.ListStoreRepo = gtkutils.GetListStore(b, "liststore_repo")
	win.ListStoreLang = gtkutils.GetListStore(b, "liststore_lang")
	win.ListStoreGames = gtkutils.GetListStore(b, "liststore_games")

	win.BtnUpdate = gtkutils.GetButton(b, "button_update")
	win.EntryKeyword = gtkutils.GetEntry(b, "entry_keyword")
	win.CmbBoxRepo = gtkutils.GetComboBox(b, "combobox_repo")
	win.CmbBoxLang = gtkutils.GetComboBox(b, "combobox_lang")
	win.ChckBtnInstalled = gtkutils.GetCheckButton(b, "checkutton_installed")
	win.BtnClear = gtkutils.GetButton(b, "button_clear")

	win.ScrWndGames = gtkutils.GetScrolledWindow(b, "scrolledwindow_games")
	win.SpinnerGames = gtkutils.GetSpinner(b, "spinner_games")

	treeViewGames := gtkutils.GetTreeView(b, "treeview_games")
	win.GamesSelection, e = treeViewGames.GetSelection()
	if e != nil {
		ShowErrorDlgFatal(e.Error(), win.Window)
	}

	win.LblGameTitle = gtkutils.GetLabel(b, "label_game_title")
	win.ImgGame = gtkutils.GetImage(b, "image_game")
	win.LblGameRepo = gtkutils.GetLabel(b, "label_game_repo")
	win.LblGameLang = gtkutils.GetLabel(b, "label_game_lang")
	win.LblGameVersion = gtkutils.GetLabel(b, "label_game_version")

	win.ScrWndGameDesc = gtkutils.GetScrolledWindow(b, "scrolledwindow_game_desc")
	win.LblGameDesc = gtkutils.GetLabel(b, "label_game_desc")

	win.BtnGameRun = gtkutils.GetButton(b, "button_game_run")
	win.BtnGameInstall = gtkutils.GetButton(b, "button_game_install")
	win.BtnGameUpdate = gtkutils.GetButton(b, "button_game_update")
	win.BtnGameRemove = gtkutils.GetButton(b, "button_game_remove")

	win.SprtrSideBox = gtkutils.GetSeparator(b, "separator_side")
	win.BxSideBox = gtkutils.GetBox(b, "box_side")

	win.MenuItmSortingReset = gtkutils.GetMenuItem(b, "menuitem_sorting_reset")
	win.ChckMenuItmSideBar = gtkutils.GetCheckMenuItem(b, "checkmenuitem_sidebar")
	win.MenuItmSettings = gtkutils.GetMenuItem(b, "menuitem_settings")
	win.MenuItmAbout = gtkutils.GetMenuItem(b, "menuitem_about")

	// todo: to constant sizes
	win.PixBufGameDefaultImage, e = gdk.PixbufNewFromFileAtScale(
		configurator.DataResourcePath(logoFilePath), 210, 210, true)

	if e != nil {
		ShowErrorDlgFatal(e.Error(), win.Window)
	}

	showSideBar := !manager.Config.Gtk.HideSidebar
	win.ChckMenuItmSideBar.SetActive(showSideBar)
	win.toggleSideBar(showSideBar)

	win.resetGameInfo()

	// Handlers
	handlers := &MainWindowHandlers{win: win}
	win.BtnUpdate.Connect("clicked", handlers.updateClicked)
	win.EntryKeyword.Connect("changed", handlers.keywordChanged)
	win.CmbBoxRepo.Connect("changed", handlers.repoChanged)
	win.CmbBoxLang.Connect("changed", handlers.langChanged)
	win.ChckBtnInstalled.Connect("clicked", handlers.installedClicked)
	win.BtnClear.Connect("clicked", handlers.clearClicked)
	treeViewGames.Connect("row_activated", handlers.gameRowActivated)
	win.GamesSelection.Connect("changed", handlers.gameChanged)
	win.BtnGameRun.Connect("clicked", handlers.runGameClicked)
	win.BtnGameInstall.Connect("clicked", handlers.installGameClicked)
	win.BtnGameUpdate.Connect("clicked", handlers.updateGameClicked)
	win.BtnGameRemove.Connect("clicked", handlers.removeGameClicked)
	win.MenuItmSortingReset.Connect("activate", handlers.sortingResetActivated)
	win.ChckMenuItmSideBar.Connect("toggled", handlers.sideBarToggled)
	win.MenuItmSettings.Connect("activate", handlers.settingsActivated)
	win.MenuItmAbout.Connect("activate", handlers.aboutActivated)
	win.Window.Connect("destroy", handlers.windowDestroyed)
	win.Window.Connect("delete_event", handlers.mainDeleted)

	width, height := win.getDefaultWindowSize(manager.Config)
	win.Window.SetDefaultSize(width, height)

	win.Window.SetTitle(title)

	// OS integrations for window
	os_integration.OsIntegrateWindow(win.Window)

	return win
}

func (win *MainWindow) toggleSideBar(show bool) {
	if show {
		win.SprtrSideBox.Show()
		win.BxSideBox.Show()
	} else {
		win.SprtrSideBox.Hide()
		win.BxSideBox.Hide()
	}
}

func (win *MainWindow) getDefaultWindowSize(config *configurator.InsteadmanConfig) (width, height int) {
	width = config.Gtk.MainWidth
	height = config.Gtk.MainHeight
	if width < 1 {
		width = 770 // todo: to constant
	}
	if height < 1 {
		height = 500 // todo: to constant
	}
	return
}

func (win *MainWindow) resetGameInfo() {
	win.LblGameTitle.SetText(win.Title + " " + win.Version)

	win.ImgGame.SetFromPixbuf(win.PixBufGameDefaultImage)

	win.ScrWndGameDesc.Hide()
	win.LblGameRepo.Hide()
	win.LblGameLang.Hide()
	win.LblGameVersion.Hide()
	win.BtnGameRun.Hide()
	win.BtnGameInstall.Hide()
	win.BtnGameUpdate.Hide()
	win.BtnGameRemove.Hide()
}

func (win *MainWindow) clearFilterValues() {
	win.CmbBoxRepo.SetSensitive(false)
	win.CmbBoxLang.SetSensitive(false)

	win.ListStoreRepo.Clear()
	iter := win.ListStoreRepo.Append()
	win.ListStoreRepo.Set(iter, []int{comboBoxColumnId, comboBoxColumnTitle}, []interface{}{"", i18n.T("Repository")})
	win.CmbBoxRepo.SetActiveID("")

	win.ListStoreLang.Clear()
	iter = win.ListStoreLang.Append()
	win.ListStoreLang.Set(iter, []int{comboBoxColumnId, comboBoxColumnTitle}, []interface{}{"", i18n.T("Language")})
	win.CmbBoxLang.SetActiveID("")

	win.CmbBoxRepo.SetSensitive(true)
	win.CmbBoxLang.SetSensitive(true)
}

func (win *MainWindow) refreshFilterValues() {
	repositories := win.Manager.GetRepositories()
	langs := win.Manager.FindLangs(win.Games)

	for _, repo := range repositories {
		iter := win.ListStoreRepo.Append()
		win.ListStoreRepo.Set(iter, []int{comboBoxColumnId, comboBoxColumnTitle}, []interface{}{repo.Name, repo.Name})
	}

	for _, lang := range langs {
		iter := win.ListStoreLang.Append()
		win.ListStoreLang.Set(iter, []int{comboBoxColumnId, comboBoxColumnTitle}, []interface{}{lang, lang})
	}
}

func (win *MainWindow) gameListStoreColumns() []int {
	return []int{gameColumnId, gameColumnTitle, gameColumnVersion, gameColumnSizeHuman, gameColumnFontWeight,
		gameColumnSize}
}

func (win *MainWindow) gameListStoreValues(g manager.Game) []interface{} {
	fontWeight := fontWeightNormal
	if g.Installed {
		fontWeight = fontWeightBold
	}

	return []interface{}{g.Id, g.Title, g.HumanVersion(), g.HumanSize(), fontWeight, g.Size}
}

func (win *MainWindow) refreshGames() {
	var e error

	log.Print("Refreshing games...")

	win.Games, e = win.Manager.GetSortedGamesByDateDesc()
	if e != nil {
		ShowErrorDlgFatal(e.Error(), win.Window)
		return
	}

	keywordP, repoP, langP, onlyInstalled := gtkutils.GetFilterValues(win.EntryKeyword, win.CmbBoxRepo,
		win.CmbBoxLang, win.ChckBtnInstalled)

	filteredGames := manager.FilterGames(win.Games, keywordP, repoP, langP, onlyInstalled)

	win.IsRefreshing = true

	win.ListStoreGames.Clear()

	for _, game := range filteredGames {
		win.ListStoreGames.InsertWithValues(nil, -1, win.gameListStoreColumns(), win.gameListStoreValues(game))
	}

	win.CurGame = nil
	win.resetGameInfo()

	log.Print("Refreshing games has finished.")

	win.IsRefreshing = false
}

func (win *MainWindow) refreshSeveralGames(upGames []manager.Game) {
	var e error

	log.Print("Refreshing several games...")

	win.Games, e = win.Manager.GetSortedGamesByDateDesc()
	if e != nil {
		ShowErrorDlgFatal(e.Error(), win.Window)
		return
	}

	// Update current (selected) game info
	if win.CurGame != nil {
		win.CurGame = manager.FindGameById(win.Games, win.CurGame.Id)
		win.updateGameInfo(win.CurGame)
	}

	var foundGames []manager.Game = nil
	for _, game := range upGames {
		foundGames = append(foundGames, manager.FindGamesByName(win.Games, game.Name)...)
	}

	for _, game := range foundGames {
		iter, _ := win.ListStoreGames.GetIterFirst()

		for iter != nil {
			value, e := win.ListStoreGames.GetValue(iter, gameColumnId)
			if e != nil {
				ShowErrorDlgFatal(e.Error(), win.Window)
				return
			}

			id, e := value.GetString()
			if e != nil {
				ShowErrorDlgFatal(e.Error(), win.Window)
				return
			}

			if id == game.Id {
				win.ListStoreGames.Set(iter, win.gameListStoreColumns(), win.gameListStoreValues(game))
			}

			if !win.ListStoreGames.IterNext(iter) {
				iter = nil
			}
		}
	}
}

func (win *MainWindow) clearFilter() {
	win.EntryKeyword.SetSensitive(false)
	win.CmbBoxRepo.SetSensitive(false)
	win.CmbBoxLang.SetSensitive(false)
	win.ChckBtnInstalled.SetSensitive(false)

	win.EntryKeyword.SetText("")
	win.CmbBoxRepo.SetActiveID("")
	win.CmbBoxLang.SetActiveID("")
	win.ChckBtnInstalled.SetActive(false)

	win.refreshGames()

	win.EntryKeyword.SetSensitive(true)
	win.CmbBoxRepo.SetSensitive(true)
	win.CmbBoxLang.SetSensitive(true)
	win.ChckBtnInstalled.SetSensitive(true)
}

func (win *MainWindow) updateGameInfo(g *manager.Game) {
	if g == nil {
		return
	}

	win.LblGameTitle.SetText(g.Title)

	if g.Description != "" {
		win.LblGameDesc.SetText(g.Description)
		win.ScrWndGameDesc.Show()
	} else {
		win.ScrWndGameDesc.Hide()
	}

	if g.RepositoryName != "" {
		win.LblGameRepo.SetText(g.RepositoryName)
		win.LblGameRepo.Show()
	} else {
		win.LblGameRepo.Hide()
	}

	if g.Languages != nil {
		win.LblGameLang.SetText(strings.Join(g.Languages, ", "))
		win.LblGameLang.Show()
	} else {
		win.LblGameLang.Hide()
	}

	if g.Version != "" {
		win.LblGameVersion.SetText(g.Version)
		win.LblGameVersion.Show()
	} else {
		win.LblGameVersion.Hide()
	}

	if g.Installed {
		win.BtnGameRun.Show()
		win.BtnGameInstall.Hide()
		win.BtnGameRemove.Show()
		if g.IsUpdateAvailable() {
			win.BtnGameUpdate.Show()
		} else {
			win.BtnGameUpdate.Hide()
		}
	} else {
		win.BtnGameRun.Hide()
		win.BtnGameInstall.Show()
		win.BtnGameRemove.Hide()
		win.BtnGameUpdate.Hide()
	}

	// Image
	go func() {
		win.updateGameImage(g)
	}()
}

func (win *MainWindow) updateGameImage(g *manager.Game) {
	gameImagePath, e := win.Manager.GetGameImage(g)
	if e == nil && gameImagePath != "" {
		win.PixBufGameImage, e = gdk.PixbufNewFromFileAtScale(gameImagePath, 210, 210, true) // todo: size to constants
		if e == nil {
			_, e := glib.IdleAdd(func() {
				// Set image if there is current game (user hasn't changed selected game)
				if win.CurGame != nil && g.Id == win.CurGame.Id {
					win.ImgGame.SetFromPixbuf(win.PixBufGameImage)
				}
			})

			if e != nil {
				log.Fatal("Change game image. IdleAdd() failed:", e)
			}
		}
	}

	if e != nil {
		log.Printf("Image error: %s", e)
	}

	if e != nil || gameImagePath == "" {
		_, e := glib.IdleAdd(func() {
			win.ImgGame.SetFromPixbuf(win.PixBufGameDefaultImage)
		})

		if e != nil {
			log.Fatal("Change game image. IdleAdd() failed:", e)
		}
	}
}

func (win *MainWindow) runGame(g *manager.Game) {
	if win.Manager.InterpreterCommand() == "" {
		ShowErrorDlg(i18n.T("INSTEAD has not found. Please add INSTEAD in the Settings."), win.Window)
		return
	}

	if win.CurGame == nil {
		ShowErrorDlg(i18n.T("No running. No game selected."), win.Window)
		return
	}

	win.Manager.RunGame(g)
	log.Printf("Running %s (%s) game...", g.Title, g.Name)
}

func (win *MainWindow) installGame(g *manager.Game, instBtn *gtk.Button) {
	if win.Manager.InterpreterCommand() == "" {
		ShowErrorDlg(i18n.T("INSTEAD has not found. Please add INSTEAD in the Settings."), win.Window)
		return
	}

	if win.CurGame == nil {
		ShowErrorDlg(i18n.T("No installing. No game selected."), win.Window)
		return
	}

	if instBtn != nil {
		instBtn.SetSensitive(false)
	}
	log.Printf("Installing %s (%s) game...", g.Title, g.Name)

	// Set installing status in the list
	iter, e := gtkutils.FindFirstIterInTreeSelection(win.ListStoreGames, win.GamesSelection)
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
	win.ListStoreGames.SetValue(iter, gameColumnSizeHuman, fmt.Sprintf(i18n.T("%s Installing..."),
		g.HumanSize()))

	installProgress := func(size uint64) {
		percents := utils.Percents(size, uint64(g.Size))
		glib.IdleAdd(func() {
			win.ListStoreGames.SetValue(iter, gameColumnSizeHuman, fmt.Sprintf(i18n.T("%s %s Installing..."),
				g.HumanSize(), percents))
		})
	}

	go func() {
		instGame := g
		instErr := win.Manager.InstallGame(instGame, installProgress)

		if instErr == nil {
			log.Print("Game has installed.")
		}

		_, e := glib.IdleAdd(func() {
			if instErr != nil {
				ShowErrorDlg(
					fmt.Sprintf(i18n.T("Game hasn't installed (%s). Please check INSTEAD in the Settings."),
						instErr.Error()), win.Window)
			}
			win.refreshSeveralGames([]manager.Game{*instGame})

			if instBtn != nil {
				instBtn.SetSensitive(true)
			}
		})

		if e != nil {
			log.Fatal("Installing game. IdleAdd() failed:", e)
		}
	}()
}

func (win *MainWindow) updateRepositories() {
	win.ScrWndGames.Hide()
	win.SpinnerGames.Show()
	win.BtnUpdate.SetSensitive(false)

	log.Print("Updating repositories...")

	go func() {
		errors := win.Manager.UpdateRepositories()
		for _, e := range errors {
			log.Printf("Update repository error: %s", e.Error())
		}
		log.Print("Repositories have updated.")

		_, e := glib.IdleAdd(func() {
			win.clearFilterValues()
			win.refreshGames()
			win.refreshFilterValues()

			win.ScrWndGames.Show()
			win.SpinnerGames.Hide()
			win.BtnUpdate.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("Updating repositories. IdleAdd() failed:", e)
		}
	}()
}

/* Handlers */
type MainWindowHandlers struct {
	win *MainWindow
}

func (h *MainWindowHandlers) updateClicked() {
	h.win.updateRepositories()
}

func (h *MainWindowHandlers) keywordChanged(s *gtk.Entry) {
	if !s.IsSensitive() {
		return
	}
	h.win.refreshGames()
}

func (h *MainWindowHandlers) repoChanged(s *gtk.ComboBox) {
	if !s.IsSensitive() {
		return
	}
	h.win.refreshGames()
}

func (h *MainWindowHandlers) langChanged(s *gtk.ComboBox) {
	if !s.IsSensitive() {
		return
	}
	h.win.refreshGames()
}

func (h *MainWindowHandlers) installedClicked(s *gtk.CheckButton) {
	if !s.IsSensitive() {
		return
	}
	h.win.refreshGames()
}

func (h *MainWindowHandlers) clearClicked(s *gtk.Button) {
	h.win.clearFilter()
}

func (h *MainWindowHandlers) gameRowActivated() {
	if h.win.CurGame == nil {
		return
	}

	if !h.win.CurGame.Installed {
		if h.win.BtnGameInstall.IsSensitive() {
			h.win.installGame(h.win.CurGame, h.win.BtnGameInstall)
		}
	} else if h.win.CurGame.IsUpdateAvailable() {
		if h.win.BtnGameUpdate.IsSensitive() {
			h.win.installGame(h.win.CurGame, h.win.BtnGameUpdate)
		}
	} else {
		h.win.runGame(h.win.CurGame)
	}
}

func (h *MainWindowHandlers) gameChanged(s *gtk.TreeSelection) {
	if h.win.IsRefreshing {
		return
	}

	iter, e := gtkutils.FindFirstIterInTreeSelection(h.win.ListStoreGames, s)
	if e != nil {
		log.Printf("Error: %v", e)
		return
	}
	if iter == nil {
		return
	}

	value, e := h.win.ListStoreGames.GetValue(iter, gameColumnId)
	if e != nil {
		ShowErrorDlgFatal(e.Error(), h.win.Window)
		return
	}

	id, e := value.GetString()
	if e != nil {
		ShowErrorDlgFatal(e.Error(), h.win.Window)
		return
	}

	h.win.CurGame = manager.FindGameById(h.win.Games, id)
	if h.win.CurGame == nil {
		ShowErrorDlgFatal(gettext.Sprintf(i18n.T("Game %s has not found."), id), h.win.Window)
		return
	}

	h.win.updateGameInfo(h.win.CurGame)
}

func (h *MainWindowHandlers) runGameClicked() {
	h.win.runGame(h.win.CurGame)
}

func (h *MainWindowHandlers) installGameClicked(s *gtk.Button) {
	// todo: CurGame as parameter

	h.win.installGame(h.win.CurGame, s)
}

func (h *MainWindowHandlers) updateGameClicked(s *gtk.Button) {
	h.win.installGame(h.win.CurGame, s)
}

func (h *MainWindowHandlers) removeGameClicked(s *gtk.Button) {
	// todo: CurGame as parameter

	s.SetSensitive(false)
	log.Printf("Removing %s (%s) game...", h.win.CurGame.Title, h.win.CurGame.Name)

	// Set removing status in the list
	iter, e := gtkutils.FindFirstIterInTreeSelection(h.win.ListStoreGames, h.win.GamesSelection)
	if e != nil {
		ShowErrorDlgFatal(e.Error(), h.win.Window)
		return
	}
	h.win.ListStoreGames.SetValue(iter, gameColumnSizeHuman,
		gettext.Sprintf(i18n.T("%s Removing..."), h.win.CurGame.HumanSize()))

	go func() {
		rmGame := h.win.CurGame
		h.win.Manager.RemoveGame(rmGame)
		log.Print("Game has removed.")

		_, e := glib.IdleAdd(func() {
			h.win.refreshSeveralGames([]manager.Game{*rmGame})
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("Removing game. IdleAdd() failed:", e)
		}
	}()
}

func (h *MainWindowHandlers) sortingResetActivated() {
	h.win.ListStoreGames.SetSortColumnId(gtk.SORT_COLUMN_UNSORTED, gtk.SORT_ASCENDING)
	h.win.refreshGames()
}

func (h *MainWindowHandlers) sideBarToggled(s *gtk.CheckMenuItem) {
	showSideBar := s.GetActive()
	h.win.toggleSideBar(showSideBar)
	h.win.Manager.Config.Gtk.HideSidebar = !showSideBar
	h.win.Configurator.SaveConfig(h.win.Manager.Config)
}

func (h *MainWindowHandlers) settingsActivated() {
	ShowSettingWin(h.win.Manager, h.win.Configurator, h.win.Version, h.win.Window)
}

func (h *MainWindowHandlers) aboutActivated() {
	ShowAboutWin(h.win.Manager, h.win.Configurator, h.win.Version, h.win.Window)
}

func (h *MainWindowHandlers) windowDestroyed() {
	gtk.MainQuit()
}

func (h *MainWindowHandlers) mainDeleted() {
	width, height := h.win.Window.GetSize()

	h.win.Manager.Config.Gtk.MainWidth = width
	h.win.Manager.Config.Gtk.MainHeight = height
	h.win.Configurator.SaveConfig(h.win.Manager.Config)
}
