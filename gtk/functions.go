package main

import (
	"../core/configurator"
	"../core/manager"
	"../core/utils"
	"./ui"
	gtkutils "./utils"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strings"
)

func RunGame(g *manager.Game) {
	if M.InterpreterCommand() == "" {
		ui.ShowErrorDlg("INSTEAD has not found. Please add INSTEAD in the Settings.", WndMain)
		return
	}

	if CurGame == nil {
		ui.ShowErrorDlg("No running. No game selected.", WndMain)
		return
	}

	M.RunGame(g)
	log.Printf("Running %s (%s) game...", g.Title, g.Name)
}

func InstallGame(g *manager.Game, instBtn *gtk.Button) {
	if M.InterpreterCommand() == "" {
		ui.ShowErrorDlg("INSTEAD has not found. Please add INSTEAD in the Settings.", WndMain)
		return
	}

	if CurGame == nil {
		ui.ShowErrorDlg("No installing. No game selected.", WndMain)
		return
	}

	if instBtn != nil {
		instBtn.SetSensitive(false)
	}
	log.Printf("Installing %s (%s) game...", g.Title, g.Name)

	// Set installing status in the list
	iter, e := gtkutils.FindFirstIterInTreeSelection(ListStoreGames, GamesSelection)
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
	ListStoreGames.SetValue(iter, GameColumnSizeHuman, g.HumanSize()+" Installing...")

	installProgress := func(size uint64) {
		percents := utils.Percents(size, uint64(g.Size))
		glib.IdleAdd(func() {
			ListStoreGames.SetValue(iter, GameColumnSizeHuman, g.HumanSize()+" Installing... "+percents)
		})
	}

	go func() {
		instGame := g
		instErr := M.InstallGame(instGame, installProgress)

		if instErr == nil {
			log.Print("Game has installed.")
		}

		_, e := glib.IdleAdd(func() {
			if instErr != nil {
				ui.ShowErrorDlg("Game hasn't installed ("+instErr.Error()+
					"). Please check INSTEAD in the Settings.", WndMain)
			}
			RefreshSeveralGames([]manager.Game{*instGame})

			if instBtn != nil {
				instBtn.SetSensitive(true)
			}
		})

		if e != nil {
			log.Fatal("Installing game. IdleAdd() failed:", e)
		}
	}()
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

	Games, e = M.GetSortedGamesByDateDesc()
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
		return
	}

	keywordP, repoP, langP, onlyInstalled := gtkutils.GetFilterValues(EntryKeyword, CmbBoxRepo, CmbBoxLang, ChckBtnInstalled)
	filteredGames := manager.FilterGames(Games, keywordP, repoP, langP, onlyInstalled)

	IsRefreshing = true

	ListStoreGames.Clear()

	for _, game := range filteredGames {
		ListStoreGames.InsertWithValues(nil, -1, gameListStoreColumns(), gameListStoreValues(game))
	}

	CurGame = nil
	resetGameInfo()

	log.Print("Refreshing games has finished.")

	IsRefreshing = false
}

func RefreshSeveralGames(upGames []manager.Game) {
	var e error

	log.Print("Refreshing several games...")

	Games, e = M.GetSortedGamesByDateDesc()
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
		return
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
				ui.ShowErrorDlgFatal(e.Error(), WndMain)
				return
			}

			id, e := value.GetString()
			if e != nil {
				ui.ShowErrorDlgFatal(e.Error(), WndMain)
				return
			}

			if id == game.Id {
				ListStoreGames.Set(iter, gameListStoreColumns(), gameListStoreValues(game))
			}

			if !ListStoreGames.IterNext(iter) {
				iter = nil
			}
		}
	}
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

func findInterpreter(m *manager.Manager, c *configurator.Configurator) {
	path := m.InterpreterFinder.Find()

	if path == nil {
		ui.ShowErrorDlg("INSTEAD has not found. Please add it in config.yml (interpreter_command)", WndMain)
		return
	}

	log.Printf("INSTEAD has found: %s", *path)

	m.Config.InterpreterCommand = *path
	e := c.SaveConfig(m.Config)
	if e != nil {
		ui.ShowErrorDlgFatal(e.Error(), WndMain)
		return
	}

	log.Print("Path has saved")
}

func resetGameInfo() {
	LblGameTitle.SetText(Title + " " + version)

	ImgGame.SetFromPixbuf(PixBufGameDefaultImage)

	ScrWndGameDesc.Hide()
	LblGameRepo.Hide()
	LblGameLang.Hide()
	LblGameVersion.Hide()
	BtnGameRun.Hide()
	BtnGameInstall.Hide()
	BtnGameUpdate.Hide()
	BtnGameRemove.Hide()
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
		if g.IsUpdateAvailable() {
			BtnGameUpdate.Show()
		} else {
			BtnGameUpdate.Hide()
		}
	} else {
		BtnGameRun.Hide()
		BtnGameInstall.Show()
		BtnGameRemove.Hide()
		BtnGameUpdate.Hide()
	}

	// Image
	go func() {
		gameImagePath, e := M.GetGameImage(g)
		if e == nil && gameImagePath != "" {
			PixBufGameImage, e = gdk.PixbufNewFromFileAtScale(gameImagePath, 210, 210, true)
			if e == nil {
				_, e := glib.IdleAdd(func() {
					// Set image if there is current game (user hasn't changed selected game)
					if CurGame != nil && g.Id == CurGame.Id {
						ImgGame.SetFromPixbuf(PixBufGameImage)
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
				ImgGame.SetFromPixbuf(PixBufGameDefaultImage)
			})

			if e != nil {
				log.Fatal("Change game image. IdleAdd() failed:", e)
			}
		}
	}()
}

func gameListStoreColumns() []int {
	return []int{GameColumnId, GameColumnTitle, GameColumnVersion, GameColumnSizeHuman, GameColumnFontWeight,
		GameColumnSize}
}

func gameListStoreValues(g manager.Game) []interface{} {
	fontWeight := FontWeightNormal
	if g.Installed {
		fontWeight = FontWeightBold
	}

	return []interface{}{g.Id, g.Title, g.HumanVersion(), g.HumanSize(), fontWeight, g.Size}
}

func ToggleSidebox(show bool) {
	if show {
		SprtrSideBox.Show()
		BxSideBox.Show()
	} else {
		SprtrSideBox.Hide()
		BxSideBox.Hide()
	}
}

func GetDefaultWindowSize(config *configurator.InsteadmanConfig) (width, height int) {
	width = config.Gtk.MainWidth
	height = config.Gtk.MainHeight
	if width < 1 {
		width = 770
	}
	if height < 1 {
		height = 500
	}
	return
}
