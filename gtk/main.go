package main

import (
	"../core/configurator"
	"../core/interpreter_finder"
	"../core/manager"
	"../core/utils"
	"./i18n"
	"./os_integration"
	"./ui"
	gtkutils "./utils"
	"github.com/gosexy/gettext"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"runtime"
)

const (
	Title = "InsteadMan"

	LogoFilePath     = "resources/images/logo.png"
	MainFormFilePath = "resources/gtk/main.glade"

	ComboBoxColumnId    = 0
	ComboBoxColumnTitle = 1

	GameColumnId         = 0
	GameColumnTitle      = 1
	GameColumnVersion    = 2
	GameColumnSizeHuman  = 3
	GameColumnFontWeight = 4
	GameColumnSize       = 5

	FontWeightNormal = 400
	FontWeightBold   = 700

	EnvDataPath = "DATA_PATH"

	I18nDomain = "insteadman"
)

var (
	version = "3"

	M            *manager.Manager
	Configurator *configurator.Configurator
	Games        []manager.Game
	CurGame      *manager.Game
	IsRefreshing bool

	WndMain *gtk.Window

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
)

func main() {
	runtime.LockOSThread()

	// OS integrations
	os_integration.OsIntegrate()

	gtk.Init(nil)

	b, e := gtk.BuilderNew()
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	executablePath, e := os.Executable()
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
	}

	currentDir, e := utils.BinAbsDir(executablePath)
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
	}

	dataPath := os.Getenv(EnvDataPath)

	Configurator = &configurator.Configurator{FilePath: "", CurrentDir: currentDir, DataPath: dataPath,
		Version: version}

	config, e := Configurator.GetConfig()
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
	}

	finder := &interpreterFinder.InterpreterFinder{CurrentDir: currentDir}

	M = &manager.Manager{Config: config, InterpreterFinder: finder}

	// I18n init
	i18n.Init(Configurator, I18nDomain, config.Lang)

	e = b.AddFromFile(Configurator.DataResourcePath(MainFormFilePath))
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
	}

	obj, e := b.GetObject("window_main")
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
	}

	var ok bool
	WndMain, ok = obj.(*gtk.Window)
	if !ok {
		ui.ShowErrorDlgFatal("No main window", WndMain)
	}

	ListStoreRepo = gtkutils.GetListStore(b, "liststore_repo")
	ListStoreLang = gtkutils.GetListStore(b, "liststore_lang")
	ListStoreGames = gtkutils.GetListStore(b, "liststore_games")

	BtnUpdate = gtkutils.GetButton(b, "button_update")
	EntryKeyword = gtkutils.GetEntry(b, "entry_keyword")
	CmbBoxRepo = gtkutils.GetComboBox(b, "combobox_repo")
	CmbBoxLang = gtkutils.GetComboBox(b, "combobox_lang")
	ChckBtnInstalled = gtkutils.GetCheckButton(b, "checkutton_installed")
	BtnClear = gtkutils.GetButton(b, "button_clear")

	ScrWndGames = gtkutils.GetScrolledWindow(b, "scrolledwindow_games")
	SpinnerGames = gtkutils.GetSpinner(b, "spinner_games")

	treeViewGames := gtkutils.GetTreeView(b, "treeview_games")
	GamesSelection, e = treeViewGames.GetSelection()
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
	}

	LblGameTitle = gtkutils.GetLabel(b, "label_game_title")
	ImgGame = gtkutils.GetImage(b, "image_game")
	LblGameRepo = gtkutils.GetLabel(b, "label_game_repo")
	LblGameLang = gtkutils.GetLabel(b, "label_game_lang")
	LblGameVersion = gtkutils.GetLabel(b, "label_game_version")

	ScrWndGameDesc = gtkutils.GetScrolledWindow(b, "scrolledwindow_game_desc")
	LblGameDesc = gtkutils.GetLabel(b, "label_game_desc")

	BtnGameRun = gtkutils.GetButton(b, "button_game_run")
	BtnGameInstall = gtkutils.GetButton(b, "button_game_install")
	BtnGameUpdate = gtkutils.GetButton(b, "button_game_update")
	BtnGameRemove = gtkutils.GetButton(b, "button_game_remove")

	SprtrSideBox = gtkutils.GetSeparator(b, "separator_side")
	BxSideBox = gtkutils.GetBox(b, "box_side")

	MenuItmSortingReset = gtkutils.GetMenuItem(b, "menuitem_sorting_reset")
	ChckMenuItmSideBar = gtkutils.GetCheckMenuItem(b, "checkmenuitem_sidebar")
	MenuItmSettings = gtkutils.GetMenuItem(b, "menuitem_settings")
	MenuItmAbout = gtkutils.GetMenuItem(b, "menuitem_about")

	if M.InterpreterCommand() == "" {
		findInterpreter(M, Configurator)
	}

	// Update repositories
	updateClicked(BtnUpdate)

	PixBufGameDefaultImage, e = gdk.PixbufNewFromFileAtScale(
		Configurator.DataResourcePath(LogoFilePath), 210, 210, true)

	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
	}

	showSideBar := !M.Config.Gtk.HideSidebar
	ChckMenuItmSideBar.SetActive(showSideBar)
	ToggleSidebox(showSideBar)

	BtnUpdate.Connect("clicked", updateClicked)
	EntryKeyword.Connect("changed", func(s *gtk.Entry) {
		if !s.IsSensitive() {
			return
		}
		RefreshGames()
	})
	CmbBoxRepo.Connect("changed", func(s *gtk.ComboBox) {
		if !s.IsSensitive() {
			return
		}
		RefreshGames()
	})
	CmbBoxLang.Connect("changed", func(s *gtk.ComboBox) {
		if !s.IsSensitive() {
			return
		}
		RefreshGames()
	})
	ChckBtnInstalled.Connect("clicked", func(s *gtk.CheckButton) {
		if !s.IsSensitive() {
			return
		}
		RefreshGames()
	})
	BtnClear.Connect("clicked", func() {
		ClearFilter()
	})

	treeViewGames.Connect("row_activated", gameRowActivated)
	GamesSelection.Connect("changed", gameChanged)

	BtnGameRun.Connect("clicked", runGameClicked)
	BtnGameInstall.Connect("clicked", installGameClicked)
	BtnGameUpdate.Connect("clicked", updateGameClicked)
	BtnGameRemove.Connect("clicked", removeGameClicked)

	MenuItmSortingReset.Connect("activate", func() {
		ListStoreGames.SetSortColumnId(gtk.SORT_COLUMN_UNSORTED, gtk.SORT_ASCENDING)
		RefreshGames()
	})
	ChckMenuItmSideBar.Connect("toggled", sideBarToggled)
	MenuItmSettings.Connect("activate", func() {
		ui.ShowSettingWin(M, Configurator, version, WndMain)
	})
	MenuItmAbout.Connect("activate", func() {
		ui.ShowAboutWin(M, Configurator, version, WndMain)
	})

	WndMain.Connect("destroy", gtk.MainQuit)
	WndMain.Connect("delete_event", mainDeleted)

	resetGameInfo()

	width, height := GetDefaultWindowSize(M.Config)
	WndMain.SetDefaultSize(width, height)

	WndMain.SetTitle(Title)
	WndMain.Show()

	// OS integrations for window
	os_integration.OsIntegrateWindow(WndMain)

	gtk.Main()
}

func updateClicked(s *gtk.Button) {
	ScrWndGames.Hide()
	SpinnerGames.Show()
	s.SetSensitive(false)

	log.Print("Updating repositories...")

	go func() {
		errors := M.UpdateRepositories()
		for _, e := range errors {
			log.Printf("Update repository error: %s", e.Error())
		}
		log.Print("Repositories have updated.")

		_, e := glib.IdleAdd(func() {
			ClearFilterValues()
			RefreshGames()
			RefreshFilterValues()

			ScrWndGames.Show()
			SpinnerGames.Hide()
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("Updating repositories. IdleAdd() failed:", e)
		}
	}()
}

func gameChanged(s *gtk.TreeSelection) {
	if IsRefreshing {
		return
	}

	iter, e := gtkutils.FindFirstIterInTreeSelection(ListStoreGames, s)
	if e != nil {
		log.Printf("Error: %v", e)
		return
	}
	if iter == nil {
		return
	}

	value, e := ListStoreGames.GetValue(iter, GameColumnId)
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
		return
	}

	id, e := value.GetString()
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
		return
	}

	CurGame = manager.FindGameById(Games, id)
	if CurGame == nil {
		ui.ShowErrorDlgFatal(gettext.Sprintf(i18n.T("Game %s has not found."), id), WndMain)
		return
	}

	updateGameInfo(CurGame)
}

func gameRowActivated() {
	if CurGame == nil {
		return
	}

	if !CurGame.Installed {
		if BtnGameInstall.IsSensitive() {
			InstallGame(CurGame, BtnGameInstall)
		}
	} else if CurGame.IsUpdateAvailable() {
		if BtnGameUpdate.IsSensitive() {
			InstallGame(CurGame, BtnGameUpdate)
		}
	} else {
		RunGame(CurGame)
	}
}

func runGameClicked() {
	RunGame(CurGame)
}

func installGameClicked(s *gtk.Button) {
	// todo: CurGame as parameter

	InstallGame(CurGame, s)
}

func updateGameClicked(s *gtk.Button) {
	InstallGame(CurGame, s)
}

func removeGameClicked(s *gtk.Button) {
	// todo: CurGame as parameter

	s.SetSensitive(false)
	log.Printf("Removing %s (%s) game...", CurGame.Title, CurGame.Name)

	// Set removing status in the list
	iter, e := gtkutils.FindFirstIterInTreeSelection(ListStoreGames, GamesSelection)
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
		return
	}
	ListStoreGames.SetValue(iter, GameColumnSizeHuman,
		gettext.Sprintf(i18n.T("%s Removing..."), CurGame.HumanSize()))

	go func() {
		rmGame := CurGame
		M.RemoveGame(rmGame)
		log.Print("Game has removed.")

		_, e := glib.IdleAdd(func() {
			RefreshSeveralGames([]manager.Game{*rmGame})
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("Removing game. IdleAdd() failed:", e)
		}
	}()
}

func sideBarToggled(s *gtk.CheckMenuItem) {
	showSideBar := s.GetActive()
	ToggleSidebox(showSideBar)
	M.Config.Gtk.HideSidebar = !showSideBar
	Configurator.SaveConfig(M.Config)
}

func mainDeleted() {
	width, height := WndMain.GetSize()

	M.Config.Gtk.MainWidth = width
	M.Config.Gtk.MainHeight = height
	Configurator.SaveConfig(M.Config)
}
