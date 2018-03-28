package main

import (
	"../core/configurator"
	"../core/manager"
	"github.com/gotk3/gotk3/gtk"
	"html"
	"log"
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
	M              *manager.Manager
	Games          []manager.Game
	CurGame        *manager.Game
	ListStoreGames *gtk.ListStore
	LblGameTitle   *gtk.Label
	BtnGameRun     *gtk.Button
	BtnGameInstall *gtk.Button
	BtnGameUpdate  *gtk.Button
	BtnGameRemove  *gtk.Button
)

func main() {
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

	c := configurator.Configurator{FilePath: ""}
	config, e := c.GetConfig()

	M = &manager.Manager{Config: config}

	repositories := M.GetRepositories()
	Games, e = M.GetSortedGames()
	langs := M.FindLangs(Games)

	listStoreRepo := GetListStore(b, "liststore_repo")
	for _, repo := range repositories {
		iter := listStoreRepo.Append()
		listStoreRepo.Set(iter, []int{ComboBoxColumnId, ComboBoxColumnTitle}, []interface{}{repo.Name, repo.Name})
	}

	listStoreLang := GetListStore(b, "liststore_lang")
	for _, lang := range langs {
		iter := listStoreLang.Append()
		listStoreLang.Set(iter, []int{ComboBoxColumnId, ComboBoxColumnTitle}, []interface{}{lang, lang})
	}

	ListStoreGames = GetListStore(b, "liststore_games")
	for _, game := range Games {
		iter := ListStoreGames.Append()

		title := html.UnescapeString(game.Title)

		fontWeight := FontWeightNormal
		if game.Installed {
			fontWeight = FontWeightBold
		}
		ListStoreGames.Set(
			iter,
			[]int{GameColumnId, GameColumnTitle, GameColumnVersion, GameColumnSize, GameColumnFontWeight},
			[]interface{}{game.Id, title, game.Version, game.GetHumanSize(), fontWeight})
	}

	treeViewGames := GetTreeView(b, "treeview_games")
	gamesSelection, e := treeViewGames.GetSelection()
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
	gamesSelection.Connect("changed", gameChanged)

	LblGameTitle = GetLabel(b, "label_game_title")

	BtnGameRun = GetButton(b, "button_game_run")
	BtnGameRun.Connect("clicked", runGameClicked)

	BtnGameInstall = GetButton(b, "button_game_install")

	BtnGameRemove = GetButton(b, "button_game_remove")

	window.SetTitle("InsteadMan 3")
	window.SetDefaultSize(770, 500)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.Connect("destroy", gtk.MainQuit)
	window.Show()

	gtk.Main()
}

func gameChanged(s *gtk.TreeSelection) {
	rows := s.GetSelectedRows(ListStoreGames)
	path := rows.Data().(*gtk.TreePath)
	iter, e := ListStoreGames.GetIter(path)
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	value, e := ListStoreGames.GetValue(iter, GameColumnId)
	if e != nil {
		log.Fatalf("Error: %v", e)
	}

	id, _ := value.GetString()
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

func updateGameInfo(g *manager.Game) {
	LblGameTitle.SetText(g.Title)

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
