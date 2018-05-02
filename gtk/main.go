package main

import (
	"../core/configurator"
	"../core/manager"
	"../core/utils"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"runtime"
)

const (
	Title   = "InsteadMan"

	LogoFilePath     = "resources/images/logo.png"
	MainFormFilePath = "resources/gtk/main.glade"

	ComboBoxColumnId    = 0
	ComboBoxColumnTitle = 1

	GameColumnId         = 0
	GameColumnTitle      = 1
	GameColumnVersion    = 2
	GameColumnSize       = 3
	GameColumnFontWeight = 4

	FontWeightNormal = 400
	FontWeightBold   = 700
)

var (
	version string = "3"

	M            *manager.Manager
	C            *configurator.Configurator
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
	//BtnGameUpdate  *gtk.Button
	BtnGameRemove *gtk.Button

	SprtrSideBox *gtk.Separator
	BxSideBox *gtk.Box

	ChckMenuItmSideBar *gtk.CheckMenuItem
	MenuItmSettings *gtk.MenuItem
	MenuItmAbout *gtk.MenuItem

	PixBufGameDefaultImage *gdk.Pixbuf
	PixBufGameImage        *gdk.Pixbuf
)

func main() {
	runtime.LockOSThread()

	gtk.Init(nil)

	b, e := gtk.BuilderNew()
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	executablePath, e := os.Executable()
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	currentDir, e := utils.BinAbsDir(executablePath)
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	C = &configurator.Configurator{FilePath: "", CurrentDir: currentDir}
	config, e := C.GetConfig()
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	M = &manager.Manager{Config: config}

	e = b.AddFromFile(C.ShareResourcePath(MainFormFilePath))
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	obj, e := b.GetObject("window_main")
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	var ok bool
	WndMain, ok = obj.(*gtk.Window)
	if !ok {
		ShowErrorDlgFatal("No main window")
	}

	ListStoreRepo = GetListStore(b, "liststore_repo")
	ListStoreLang = GetListStore(b, "liststore_lang")
	ListStoreGames = GetListStore(b, "liststore_games")

	BtnUpdate = GetButton(b, "button_update")
	EntryKeyword = GetEntry(b, "entry_keyword")
	CmbBoxRepo = GetComboBox(b, "combobox_repo")
	CmbBoxLang = GetComboBox(b, "combobox_lang")
	ChckBtnInstalled = GetCheckButton(b, "checkutton_installed")
	BtnClear = GetButton(b, "button_clear")

	ScrWndGames = GetScrolledWindow(b, "scrolledwindow_games")
	SpinnerGames = GetSpinner(b, "spinner_games")

	treeViewGames := GetTreeView(b, "treeview_games")
	GamesSelection, e = treeViewGames.GetSelection()
	if e != nil {
		ShowErrorDlgFatal(e.Error())
	}

	LblGameTitle = GetLabel(b, "label_game_title")
	ImgGame = GetImage(b, "image_game")
	LblGameRepo = GetLabel(b, "label_game_repo")
	LblGameLang = GetLabel(b, "label_game_lang")
	LblGameVersion = GetLabel(b, "label_game_version")

	ScrWndGameDesc = GetScrolledWindow(b, "scrolledwindow_game_desc")
	LblGameDesc = GetLabel(b, "label_game_desc")

	BtnGameRun = GetButton(b, "button_game_run")
	BtnGameInstall = GetButton(b, "button_game_install")
	BtnGameRemove = GetButton(b, "button_game_remove")

	SprtrSideBox = GetSeparator(b,"separator_side")
	BxSideBox = GetBox(b,"box_side")

	ChckMenuItmSideBar = GetCheckMenuItem(b, "checkmenuitem_sidebar")
	MenuItmSettings = GetMenuItem(b, "menuitem_settings")
	MenuItmAbout = GetMenuItem(b, "menuitem_about")

	if M.Config.InterpreterCommand == "" {
		findInterpreter(M, C)
	}

	if M.HasDownloadedRepositories() {
		ClearFilterValues()
		RefreshGames()
		RefreshFilterValues()
	} else {
		// Update repositories
		updateClicked(BtnUpdate)
	}

	PixBufGameDefaultImage, e = gdk.PixbufNewFromFileAtScale(
		C.ShareResourcePath(LogoFilePath), 210, 210, true)

	if e != nil {
		ShowErrorDlgFatal(e.Error())
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
	BtnGameRemove.Connect("clicked", removeGameClicked)

	ChckMenuItmSideBar.Connect("toggled", sideBarToggled)
	MenuItmSettings.Connect("activate", func() {
		log.Println("Settings...")
	})
	MenuItmAbout.Connect("activate", func() {
		log.Println("About...")
	})

	WndMain.Connect("destroy", gtk.MainQuit)
	WndMain.Connect("delete_event", mainDeleted)

	resetGameInfo()

	width, height := GetDefaultWindowSize(M.Config)
	WndMain.SetDefaultSize(width, height)

	WndMain.SetTitle(Title + " " + version)
	WndMain.SetPosition(gtk.WIN_POS_CENTER)
	WndMain.Show()


	gtk.Main()
}

func updateClicked(s *gtk.Button) {
	ScrWndGames.Hide()
	SpinnerGames.Show()
	s.SetSensitive(false)

	log.Print("Updating repositories...")

	go func() {
		M.UpdateRepositories()
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

	iter, e := FindFirstIterInTreeSelection(ListStoreGames, s)
	if e != nil {
		log.Printf("Error: %v", e)
		return
	}
	if iter == nil {
		return
	}

	value, e := ListStoreGames.GetValue(iter, GameColumnId)
	if e != nil {
		ShowErrorDlgFatal(e.Error())
		return
	}

	id, e := value.GetString()
	if e != nil {
		ShowErrorDlgFatal(e.Error())
		return
	}

	CurGame = manager.FindGameById(Games, id)
	if CurGame == nil {
		ShowErrorDlgFatal("Game " + id + " has not found")
		return
	}

	updateGameInfo(CurGame)
}

func gameRowActivated() {
	if CurGame == nil {
		return
	}

	if CurGame.Installed {
		RunGame(CurGame)
	} else {
		InstallGame(CurGame, BtnGameInstall)
	}
}

func runGameClicked() {
	RunGame(CurGame)
}

func installGameClicked(s *gtk.Button) {
	// todo: CurGame as parameter

	InstallGame(CurGame, s)
}

func removeGameClicked(s *gtk.Button) {
	// todo: CurGame as parameter

	s.SetSensitive(false)
	log.Printf("Removing %s (%s) game...", CurGame.Title, CurGame.Name)

	// Set removing status in the list
	iter, e := FindFirstIterInTreeSelection(ListStoreGames, GamesSelection)
	if e != nil {
		ShowErrorDlgFatal(e.Error())
		return
	}
	ListStoreGames.SetValue(iter, GameColumnSize, CurGame.GetHumanSize()+" Removing...")

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
	C.SaveConfig(M.Config)
}

func mainDeleted() {
	width, height := WndMain.GetSize()

	M.Config.Gtk.MainWidth = width
	M.Config.Gtk.MainHeight = height
	C.SaveConfig(M.Config)
}

func ShowErrorDlgFatal(txt string) {
	showErrorDlg(txt, true)
}

func ShowErrorDlg(txt string) {
	showErrorDlg(txt, false)
}

func showErrorDlg(txt string, fatal bool) {
	log.Printf("Error: %v", txt)

	dlg, _ := gtk.DialogNew()
	dlg.SetTitle(Title + " " + version)
	dlg.AddButton("Close" ,gtk.RESPONSE_ACCEPT)
	dlgBox, _ := dlg.GetContentArea()
	dlgBox.SetSpacing(6)

	lbl, _ := gtk.LabelNew(txt)
	lbl.SetMarginStart(6)
	lbl.SetMarginEnd(6)
	dlgBox.Add(lbl)
	lbl.Show()

	dlg.SetModal(true)
	dlg.SetPosition(gtk.WIN_POS_CENTER)
	dlg.SetResizable(false)
	//dlg.SetTransientFor(window)

	response := dlg.Run()
	dlg.SetKeepAbove(true)
	if response == int(gtk.RESPONSE_ACCEPT) {
		dlg.Destroy()
		if fatal {
			os.Exit(1)
		}
	}

	if fatal {
		dlg.Connect("destroy", gtk.MainQuit)
	}
}
