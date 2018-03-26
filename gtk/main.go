package main

import (
	"github.com/gotk3/gotk3/gtk"
	"fmt"
	"os"
	"../core/configurator"
	"../core/manager"
)

func main() {
	c := configurator.Configurator{FilePath: ""}
	config, e := c.GetConfig()

	m := manager.Manager{Config: config}

	repositories := m.GetRepositories()
	games, e := m.GetSortedGames()
	langs := m.FindLangs(games)

	gtk.Init(nil)

	//_, e := gtk.BuilderNewFromFile("./resources/gtk/about.glade")
	b, e := gtk.BuilderNew()
	if e != nil {
		fmt.Printf("E: %v\n", e)
	}
	e = b.AddFromFile("./resources/gtk/main.glade")
	if e != nil {
		fmt.Printf("E: %v\n", e)
	}

	obj, e := b.GetObject("window_main")
	if e != nil {
		os.Exit(1)
	}
	window, ok := obj.(*gtk.Window)
	if !ok {
		os.Exit(1)
	}

	obj, e = b.GetObject("liststore_repo")
	listStoreRepo, ok := obj.(*gtk.ListStore)
	for _, repo := range repositories {
		iter := listStoreRepo.Append()
		listStoreRepo.Set(iter, []int{0, 1}, []interface{}{repo.Name, repo.Name})
	}

	obj, e = b.GetObject("liststore_lang")
	listStoreLang, ok := obj.(*gtk.ListStore)
	for _, lang := range langs {
		iter := listStoreLang.Append()
		listStoreLang.Set(iter, []int{0, 1}, []interface{}{lang, lang})
	}

	obj, e = b.GetObject("liststore_games")
	listStoreGames, ok := obj.(*gtk.ListStore)
	for _, game := range games {
		iter := listStoreGames.Append()
		listStoreGames.Set(iter, []int{0, 1, 2}, []interface{}{game.Title, game.Version, game.GetHumanSize()})
	}


	window.SetTitle("InsteadMan 3")
	window.SetDefaultSize(770, 500)
	window.Connect("destroy", destroy)
	window.ShowAll()

	gtk.Main()
}

func destroy() {
	gtk.MainQuit()
}

