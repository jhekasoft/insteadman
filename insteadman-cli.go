package main

import (
	"./configurator"
	"./interpreter_finder"
	"./manager"
	"./utils/cli"
	"fmt"
	"os"
	"strings"
)

const version = "3.0.2"

func main() {
	m, c := initManagerAndConfigurator()

	needRepositoriesUpdate := !m.HasDownloadedRepositories()

	argsWithoutProg := os.Args[1:]

	switch strings.ToLower(cli.GetCommand(argsWithoutProg)) {
	case "update":
		update(m)

	case "list":
		if needRepositoriesUpdate {
			update(m)
		}

		list(m, argsWithoutProg)

	case "search":
		if needRepositoriesUpdate {
			update(m)
		}

		search(m, argsWithoutProg)

	case "run":
		m, _ = checkInterpreterAndReinit(m, c)

		run(m, argsWithoutProg)

	case "install":
		m, _ = checkInterpreterAndReinit(m, c)

		install(m, argsWithoutProg)

	case "remove":
		remove(m, argsWithoutProg)

	case "findinterpreter":
		findInterpreter(m, c)

	case "repositories":
		repositories(m)

	case "langs":
		if needRepositoriesUpdate {
			update(m)
		}

		langs(m)

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

func list(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	cli.ExitIfError(e)

	// Parse args without "list" command
	repository, lang, onlyInstalled := getGamesFilterValues(args[1:])

	if repository != nil || lang != nil || onlyInstalled {
		games = manager.FilterGames(games, nil, repository, lang, onlyInstalled)
	}

	printGames(games)
}

func search(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	cli.ExitIfError(e)

	keyword := cli.GetCommandArg(args)
	if keyword == nil {
		printHelpAndExit()
	}

	// Parse args without "search [keyword]" command
	repository, lang, onlyInstalled := getGamesFilterValues(args[2:])

	filteredGames := manager.FilterGames(games, keyword, repository, lang, onlyInstalled)

	printGames(filteredGames)
}

func install(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	cli.ExitIfError(e)

	keyword := cli.GetCommandArg(args)
	if keyword == nil {
		printHelpAndExit()
	}

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	fmt.Printf("Downloading and installing game %s...\n", game.Title)

	e = m.InstallGame(&game)
	cli.ExitIfError(e)

	fmt.Printf("Game %s has installed.\n", game.Title)
}

func run(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	cli.ExitIfError(e)

	keyword := cli.GetCommandArg(args)
	if keyword == nil {
		printHelpAndExit()
	}

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	e = m.RunGame(&game)
	cli.ExitIfError(e)

	fmt.Printf("Running %s game...\n", game.Title)
}

func remove(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	cli.ExitIfError(e)

	keyword := cli.GetCommandArg(args)
	if keyword == nil {
		printHelpAndExit()
	}

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	fmt.Printf("Removing game %s...\n", game.Title)

	e = m.RemoveGame(&game)
	cli.ExitIfError(e)

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
	cli.ExitIfError(e)

	fmt.Println("Path has saved")
}

func repositories(m *manager.Manager) {
	for _, repo := range m.GetRepositories() {
		fmt.Printf("%s (%s)\n", repo.Name, repo.Url)
	}
}

func langs(m *manager.Manager) {
	games, e := m.GetSortedGames()
	cli.ExitIfError(e)

	for _, lang := range m.FindLangs(games) {
		fmt.Printf("%s\n", lang)
	}
}

func printHelpAndExit() {
	fmt.Printf("InsteadMan CLI %s — INSTEAD games manager\n\n"+
		"Usage:\n"+
		"    insteadman-cli [command] [keyword]\n\n"+
		"Commands:\n"+
		"update\n    Update game's repositories\n"+
		"list --repo=[name] --lang=[lang] --installed\n    Print list of games with filtering\n"+
		"search [keyword] --repo=[name] --lang=[lang] --installed\n    Search game by name and title with filtering\n"+
		"install [keyword]\n    Install game by keyword\n"+
		"run [keywork]\n    Run game by keyword\n"+
		"remove [keywork]\n    Remove game by keyword\n"+
		"findInterpreter\n    Find INSTEAD interpreter and save path to the config\n"+
		"repositories\n    Pring available repositories\n"+
		"langs\n    Pring available game languages\n\n"+
		"More info: https://github.com/jhekasoft/insteadman3\n", version)
	os.Exit(1)
}

// -- Commands -----------------------------------

func initManagerAndConfigurator() (*manager.Manager, *configurator.Configurator) {
	c := configurator.Configurator{FilePath: ""}
	config, e := c.GetConfig()
	cli.ExitIfError(e)

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

func getGamesFilterValues(args []string) (*string, *string, bool) {
	repository := cli.FindStringArg("--repository", args)
	lang := cli.FindStringArg("--lang", args)
	onlyInstalled := cli.FindBoolArg("--installed", args)

	return repository, lang, onlyInstalled
}
