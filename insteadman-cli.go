package main

import (
	"./configurator"
	"./interpreter_finder"
	"./manager"
	"fmt"
	"os"
	"strings"
)

const version = "3.0.1"

func main() {
	m, c := initManagerAndConfigurator()

	needRepositoriesUpdate := !m.HasDownloadedRepositories()

	argsWithoutProg := os.Args[1:]

	switch strings.ToLower(getCommand(argsWithoutProg)) {
	case "update":
		update(m)

	case "list":
		if needRepositoriesUpdate {
			update(m)
		}

		list(m)

	case "search":
		if needRepositoriesUpdate {
			update(m)
		}

		keyword := getCommandArg(argsWithoutProg)
		if keyword == "" {
			printHelpAndExit()
		}

		search(m, &keyword)

	case "run":
		m, c = checkInterpreterAndReinit(m, c)

		keyword := getCommandArg(argsWithoutProg)
		if keyword == "" {
			printHelpAndExit()
		}

		run(m, &keyword)

	case "install":
		m, c = checkInterpreterAndReinit(m, c)

		keyword := getCommandArg(argsWithoutProg)
		if keyword == "" {
			printHelpAndExit()
		}

		install(m, &keyword)

	case "remove":
		keyword := getCommandArg(argsWithoutProg)
		if keyword == "" {
			printHelpAndExit()
		}

		remove(m, &keyword)

	case "findinterpreter":
		findInterpreter(m, c)

	default:
		printHelpAndExit()
	}
}

// -- Commands -----------------------------------
func update(m *manager.Manager) {
	fmt.Println("Updating repositories...")
	errors := m.UpdateRepositories()

	if errors != nil {
		fmt.Println("There are errors:")
	}
	for _, e := range errors {
		fmt.Printf("%s\n", e)
	}

	fmt.Println("Repositories have updated.")
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

	game := getOrExitIfNoGame(filteredGames, *keyword)

	fmt.Printf("Downloading and installing game %s...\n", game.Title)

	e = m.InstallGame(&game)
	exitIfError(e)

	fmt.Printf("Game %s has installed.\n", game.Title)
}

func run(m *manager.Manager, keyword *string) {
	games, e := m.GetSortedGames()

	exitIfError(e)

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	e = m.RunGame(&game)
	exitIfError(e)

	fmt.Printf("Running %s game...\n", game.Title)
}

func remove(m *manager.Manager, keyword *string) {
	games, e := m.GetSortedGames()

	exitIfError(e)

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	fmt.Printf("Removing game %s...\n", game.Title)

	e = m.RemoveGame(&game)
	exitIfError(e)

	fmt.Printf("Game %s has removed.\n", game.Title)
}

func findInterpreter(m *manager.Manager, c *configurator.Configurator) {
	finder := interpreterFinder.InterpreterFinder{Config: m.Config}
	path := finder.Find()

	if path == nil {
		fmt.Println("INSTEAD has not found. Please add it in config.yml (interpreter_command)")
		return
	}

	fmt.Printf("INSTEAD has found: %s\n", *path)

	m.Config.InterpreterCommand = *path
	e := c.SaveConfig(m.Config)
	exitIfError(e)

	fmt.Println("Path has saved")
}

func printHelpAndExit() {
	fmt.Printf("InsteadMan CLI %s — INSTEAD games manager\n\n"+
		"Usage:\n"+
		"    insteadman-cli [command] [keyword]\n\n"+
		"Commands:\n"+
		"update\n    Update game's repositories\n"+
		"list\n    Print list of games\n"+
		"search [keyword]\n    Search game by name and title\n"+
		"install [keyword]\n    Install game by keyword\n"+
		"run [keywork]\n    Run game by keyword\n"+
		"remove [keywork]\n    Remove game by keyword\n"+
		"findInterpreter\n    Find INSTEAD interpreter and save path to the config\n\n"+
		"More info: https://github.com/jhekasoft/insteadman3\n", version)
	os.Exit(1)
}

// -- Commands -----------------------------------

func initManagerAndConfigurator() (*manager.Manager, *configurator.Configurator) {
	c := configurator.Configurator{FilePath: ""}
	config, e := c.GetConfig()
	exitIfError(e)

	m := manager.Manager{Config: config}

	return &m, &c
}

func checkInterpreterAndReinit(m *manager.Manager, c *configurator.Configurator) (*manager.Manager, *configurator.Configurator) {
	if m.Config.InterpreterCommand == "" {
		findInterpreter(m, c)
		m, c = initManagerAndConfigurator()
	}

	return m, c
}

func getCommand(argsWithoutProg []string) string {
	if len(argsWithoutProg) > 0 {
		return argsWithoutProg[0]
	}

	return ""
}

func getCommandArg(argsWithoutProg []string) string {
	if len(argsWithoutProg) > 1 {
		return argsWithoutProg[1]
	}

	return ""
}

func printGames(games []manager.Game) {
	fmt.Println("Games:")

	for _, game := range games {
		installed := ""
		if game.Installed {
			installed = "[installed]"
		}

		fmt.Printf("%v, %v, %v %v %v\n", game.Title, game.Name, game.RepositoryName, game.Languages, installed)
	}
}

func getOrExitIfNoGame(filteredGames []manager.Game, keyword string) manager.Game {
	if len(filteredGames) < 1 {
		fmt.Printf("Game %s has not found\n", keyword)
		os.Exit(1)
	}

	for _, game := range filteredGames {
		if strings.ToLower(game.Name) == strings.ToLower(keyword) {
			return game
		}
	}

	return filteredGames[0]
}

func exitIfError(e error) {
	if e == nil {
		return
	}

	fmt.Printf("Error: %v\n", e)
	os.Exit(1)
}
