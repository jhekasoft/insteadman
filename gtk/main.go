package main

import (
	"../core/configurator"
	"../core/interpreter_finder"
	"../core/manager"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"runtime"
	"strings"
)

const (
	ComboBoxColumnId    = 0
	ComboBoxColumnTitle = 1

	GameColumnId         = 0
	GameColumnTitle      = 1
	GameColumnVersion    = 2
	GameColumnSize       = 3
	GameColumnFontWeight = 4

	FontWeightNormal = 0
	FontWeightBold   = 700
)

var (
	M       *manager.Manager
	Games   []manager.Game
	CurGame *manager.Game

	ListStoreGames *gtk.ListStore
	ListStoreRepo  *gtk.ListStore
	ListStoreLang  *gtk.ListStore

	BtnUpdate        *gtk.Button
	EntryKeyword     *gtk.Entry
	CmbBoxRepo       *gtk.ComboBox
	CmbBoxLang       *gtk.ComboBox
	ChckBtnInstalled *gtk.CheckButton
	BtnClear         *gtk.Button

	ScrWndGames  *gtk.ScrolledWindow
	SpinnerGames *gtk.Spinner

	LblGameTitle   *gtk.Label
	LblGameRepo    *gtk.Label
	LblGameLang    *gtk.Label
	LblGameVersion *gtk.Label
	ScrWndGameDesc *gtk.ScrolledWindow
	LblGameDesc    *gtk.Label
	BtnGameRun     *gtk.Button
	BtnGameInstall *gtk.Button
	//BtnGameUpdate  *gtk.Button
	BtnGameRemove *gtk.Button
)

func main() {
	runtime.LockOSThread()

	gtk.Init(nil)

	b, e := gtk.BuilderNew()
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
	e = b.AddFromFile("./resources/gtk/main.glade")
	if e != nil {
		log.Fatalf("Error: %v\n", e)
	}

	obj, e := b.GetObject("window_main")
	if e != nil {
		log.Fatalf("Error: %s", e)
	}
	window, ok := obj.(*gtk.Window)
	if !ok {
		log.Fatalf("No main window")
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
	gamesSelection, e := treeViewGames.GetSelection()
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	LblGameTitle = GetLabel(b, "label_game_title")
	LblGameRepo = GetLabel(b, "label_game_repo")
	LblGameLang = GetLabel(b, "label_game_lang")
	LblGameVersion = GetLabel(b, "label_game_version")

	ScrWndGameDesc = GetScrolledWindow(b, "scrolledwindow_game_desc")
	LblGameDesc = GetLabel(b, "label_game_desc")

	BtnGameRun = GetButton(b, "button_game_run")
	BtnGameInstall = GetButton(b, "button_game_install")
	BtnGameRemove = GetButton(b, "button_game_remove")

	c := configurator.Configurator{FilePath: ""}
	config, e := c.GetConfig()

	M = &manager.Manager{Config: config}

	if M.Config.InterpreterCommand == "" {
		findInterpreter(M, &c)
	}

	if M.HasDownloadedRepositories() {
		ClearFilterValues()
		RefreshGames()
		RefreshFilterValues()
	} else {
		// Update repositories
		updateClicked(BtnUpdate)
	}

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

	gamesSelection.Connect("changed", gameChanged)

	BtnGameRun.Connect("clicked", runGameClicked)
	BtnGameInstall.Connect("clicked", installGameClicked)
	BtnGameRemove.Connect("clicked", removeGameClicked)

	window.SetTitle("InsteadMan 3")
	window.SetDefaultSize(770, 500)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.Connect("destroy", gtk.MainQuit)
	window.Show()

	gtk.Main()
}

func findInterpreter(m *manager.Manager, c *configurator.Configurator) {
	finder := interpreterFinder.InterpreterFinder{Config: m.Config}
	path := finder.Find()

	if path == nil {
		log.Print("INSTEAD has not found. Please add it in config.yml (interpreter_command)")
		return
	}

	log.Printf("INSTEAD has found: %s", *path)

	m.Config.InterpreterCommand = *path
	e := c.SaveConfig(m.Config)
	if e != nil {
		log.Fatalf("Error: %v\n", e)
	}

	log.Print("Path has saved")
}

func ClearFilter() {
	EntryKeyword.SetSensitive(false)
	CmbBoxRepo.SetSensitive(false)
	CmbBoxLang.SetSensitive(false)
	ChckBtnInstalled.SetSensitive(false)

	EntryKeyword.SetText("")
	CmbBoxRepo.SetActiveID("")
	CmbBoxLang.SetActiveID("")
	ChckBtnInstalled.SetActive(false)

	RefreshGames()

	EntryKeyword.SetSensitive(true)
	CmbBoxRepo.SetSensitive(true)
	CmbBoxLang.SetSensitive(true)
	ChckBtnInstalled.SetSensitive(true)
}

func ClearFilterValues() {
	CmbBoxRepo.SetSensitive(false)
	CmbBoxLang.SetSensitive(false)

	ListStoreRepo.Clear()
	iter := ListStoreRepo.Append()
	ListStoreRepo.Set(iter, []int{ComboBoxColumnId, ComboBoxColumnTitle}, []interface{}{"", "Repository"})
	CmbBoxRepo.SetActiveID("")

	ListStoreLang.Clear()
	iter = ListStoreLang.Append()
	ListStoreLang.Set(iter, []int{ComboBoxColumnId, ComboBoxColumnTitle}, []interface{}{"", "Language"})
	CmbBoxLang.SetActiveID("")

	CmbBoxRepo.SetSensitive(true)
	CmbBoxLang.SetSensitive(true)
}

func RefreshFilterValues() {
	repositories := M.GetRepositories()
	langs := M.FindLangs(Games)

	for _, repo := range repositories {
		iter := ListStoreRepo.Append()
		ListStoreRepo.Set(iter, []int{ComboBoxColumnId, ComboBoxColumnTitle}, []interface{}{repo.Name, repo.Name})
	}

	for _, lang := range langs {
		iter := ListStoreLang.Append()
		ListStoreLang.Set(iter, []int{ComboBoxColumnId, ComboBoxColumnTitle}, []interface{}{lang, lang})
	}
}

func RefreshGames() {
	var e error

	log.Print("Refreshing games...")

	Games, e = M.GetSortedGames()
	if e != nil {
		log.Fatalf("Error: %s", e)
	}

	// Update current (selected) game info
	if CurGame != nil {
		CurGame = manager.FindGameById(Games, CurGame.Id)
		updateGameInfo(CurGame)
	}

	keywordP, repoP, langP, onlyInstalled := GetFilterValues(EntryKeyword, CmbBoxRepo, CmbBoxLang, ChckBtnInstalled)
	filteredGames := manager.FilterGames(Games, keywordP, repoP, langP, onlyInstalled)

	ListStoreGames.Clear()

	for _, game := range filteredGames {
		iter := ListStoreGames.Append()
		setListStoreGamesItem(ListStoreGames, iter, game)
	}
}

func RefreshSeveralGames(upGames []manager.Game) {
	var e error

	log.Print("Refreshing several games...")

	Games, e = M.GetSortedGames()
	if e != nil {
		log.Fatalf("Error: %s", e)
	}

	// Update current (selected) game info
	if CurGame != nil {
		CurGame = manager.FindGameById(Games, CurGame.Id)
		updateGameInfo(CurGame)
	}

	var foundGames []manager.Game = nil
	for _, game := range upGames {
		foundGames = append(foundGames, manager.FindGamesByName(Games, game.Name)...)
	}

	for _, game := range foundGames {
		iter, _ := ListStoreGames.GetIterFirst()

		for iter != nil {
			value, e := ListStoreGames.GetValue(iter, GameColumnId)
			if e != nil {
				log.Fatalf("Error: %v", e)
			}

			id, e := value.GetString()
			if e != nil {
				log.Fatalf("Error: %v", e)
			}

			if id == game.Id {
				setListStoreGamesItem(ListStoreGames, iter, game)
			}

			if !ListStoreGames.IterNext(iter) {
				iter = nil
			}
		}
	}
}

func setListStoreGamesItem(ls *gtk.ListStore, iter *gtk.TreeIter, game manager.Game) {
	fontWeight := FontWeightNormal
	if game.Installed {
		fontWeight = FontWeightBold
	}
	ls.Set(
		iter,
		[]int{GameColumnId, GameColumnTitle, GameColumnVersion, GameColumnSize, GameColumnFontWeight},
		[]interface{}{game.Id, game.Title, game.Version, game.GetHumanSize(), fontWeight})
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
	rows := s.GetSelectedRows(ListStoreGames)
	if rows.Length() < 1 {
		return
	}

	path := rows.Data().(*gtk.TreePath)
	iter, e := ListStoreGames.GetIter(path)
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	value, e := ListStoreGames.GetValue(iter, GameColumnId)
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	id, e := value.GetString()
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	CurGame = manager.FindGameById(Games, id)
	if CurGame == nil {
		log.Printf("Game %s has not found", id)
		return
	}

	updateGameInfo(CurGame)
}

func runGameClicked() {
	M.RunGame(CurGame)
	log.Printf("Running %s (%s) game...", CurGame.Title, CurGame.Name)
}

func installGameClicked(s *gtk.Button) {
	// todo: CurGame as parameter

	s.SetSensitive(false)
	log.Printf("Installing %s (%s) game...", CurGame.Title, CurGame.Name)

	go func() {
		instGame := CurGame
		M.InstallGame(instGame)
		log.Print("Game has installed.")

		_, e := glib.IdleAdd(func() {
			RefreshSeveralGames([]manager.Game{*instGame})
			s.SetSensitive(true)
		})

		if e != nil {
			log.Fatal("Installing game. IdleAdd() failed:", e)
		}
	}()
}

func removeGameClicked(s *gtk.Button) {
	// todo: CurGame as parameter

	s.SetSensitive(false)
	log.Printf("Removing %s (%s) game...", CurGame.Title, CurGame.Name)

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

func updateGameInfo(g *manager.Game) {
	if g == nil {
		return
	}

	LblGameTitle.SetText(g.Title)

	if g.Description != "" {
		LblGameDesc.SetText(g.Description)
		ScrWndGameDesc.Show()
	} else {
		ScrWndGameDesc.Hide()
	}

	if g.RepositoryName != "" {
		LblGameRepo.SetText(g.RepositoryName)
		LblGameRepo.Show()
	} else {
		LblGameRepo.Hide()
	}

	if g.Languages != nil {
		LblGameLang.SetText(strings.Join(g.Languages, ", "))
		LblGameLang.Show()
	} else {
		LblGameLang.Hide()
	}

	if g.Version != "" {
		LblGameVersion.SetText(g.Version)
		LblGameVersion.Show()
	} else {
		LblGameVersion.Hide()
	}

	if g.Installed {
		BtnGameRun.Show()
		BtnGameInstall.Hide()
		BtnGameRemove.Show()
	} else {
		BtnGameRun.Hide()
		BtnGameInstall.Show()
		BtnGameRemove.Hide()
	}
}
