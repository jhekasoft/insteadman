package main

import (
	"../core/configurator"
	"../core/interpreter_finder"
	"../core/manager"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strings"
)

func RunGame(g *manager.Game) {
	if CurGame == nil {
		log.Print("No running. No game selected.")
		return
	}

	M.RunGame(g)
	log.Printf("Running %s (%s) game...", g.Title, g.Name)
}

func InstallGame(g *manager.Game, instBtn *gtk.Button) {
	if CurGame == nil {
		log.Print("No installing. No game selected.")
		return
	}

	if instBtn != nil {
		instBtn.SetSensitive(false)
	}
	log.Printf("Installing %s (%s) game...", g.Title, g.Name)

	// Set installing status in the list
	iter, e := FindFirstIterInTreeSelection(ListStoreGames, GamesSelection)
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
	ListStoreGames.SetValue(iter, GameColumnSize, g.GetHumanSize()+" Installing...")

	go func() {
		instGame := g
		M.InstallGame(instGame)
		log.Print("Game has installed.")

		_, e := glib.IdleAdd(func() {
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

	Games, e = M.GetSortedGames()
	if e != nil {
		log.Fatalf("Error: %s", e)
	}

	keywordP, repoP, langP, onlyInstalled := GetFilterValues(EntryKeyword, CmbBoxRepo, CmbBoxLang, ChckBtnInstalled)
	filteredGames := manager.FilterGames(Games, keywordP, repoP, langP, onlyInstalled)

	IsRefreshing = true

	ListStoreGames.Clear()

	for _, game := range filteredGames {
		ListStoreGames.InsertWithValues(nil, -1, gameListStoreColumns(), gameListStoreValues(game))
	}

	CurGame = nil
	resetGameInfo()

	IsRefreshing = false
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

func resetGameInfo() {
	LblGameTitle.SetText(Title + " " + Version)

	ImgGame.SetFromPixbuf(PixBufGameDefaultImage)

	ScrWndGameDesc.Hide()
	LblGameRepo.Hide()
	LblGameLang.Hide()
	LblGameVersion.Hide()
	BtnGameRun.Hide()
	BtnGameInstall.Hide()
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
	} else {
		BtnGameRun.Hide()
		BtnGameInstall.Show()
		BtnGameRemove.Hide()
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
	return []int{GameColumnId, GameColumnTitle, GameColumnVersion, GameColumnSize, GameColumnFontWeight}
}

func gameListStoreValues(g manager.Game) []interface{} {
	fontWeight := FontWeightNormal
	if g.Installed {
		fontWeight = FontWeightBold
	}

	return []interface{}{g.Id, g.Title, g.Version, g.GetHumanSize(), fontWeight}
}
