package main

import (
	"./configurator"
	//"./interpreter_finder"
	"./manager"
	"fmt"
	"os"
)

func main() {
	conf := configurator.Configurator{FilePath: ""}
	config, e := conf.GetConfig()
	if e != nil {
		fmt.Printf("Error: %v\n", e)
		os.Exit(1)
	}

	m := manager.Manager{Config: config}

	argsWithoutProg := os.Args[1:]
	command := argsWithoutProg[0]
	switch command {
	case "list":
		list(&m)
	case "search":
		keyword := argsWithoutProg[1]
		search(&m, &keyword)
	case "run":
		keyword := argsWithoutProg[1]
		run(&m, &keyword)
	case "install":
		keyword := argsWithoutProg[1]
		install(&m, &keyword)
	case "remove":
		keyword := argsWithoutProg[1]
		remove(&m, &keyword)
	}
}

func list(m *manager.Manager) {
	games, e := m.GetSortedGames()

	exitIfError(e)

	printGames(games)
}

func search(m *manager.Manager, keyword *string) {
	games, e := m.GetSortedGames()

	exitIfError(e)

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	printGames(filteredGames)
}

func install(m *manager.Manager, keyword *string) {
	games, e := m.GetSortedGames()

	exitIfError(e)

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := filteredGames[0]
	m.InstallGame(&game)
}

func run(m *manager.Manager, keyword *string) {
	games, e := m.GetSortedGames()

	exitIfError(e)

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := filteredGames[0]
	m.RunGame(&game)
}

func remove(m *manager.Manager, keyword *string) {
	games, e := m.GetSortedGames()

	exitIfError(e)

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := filteredGames[0]
	m.RemoveGame(&game)
}

func printGames(games []manager.Game) {
	for _, game := range games {
		fmt.Printf("%v, %v, %v, %v\n", game.Title, game.Name, game.RepositoryName, game.Languages)
	}
}

func exitIfError(e error) {
	if e == nil {
		return
	}

	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}
