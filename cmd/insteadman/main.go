package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/interpreterfinder"
	"github.com/jhekasoft/insteadman3/core/manager"
	"github.com/jhekasoft/insteadman3/core/utils"
)

var version = "3"

func main() {
	m, c := initManagerAndConfigurator()
	needRepositoriesUpdate := !m.HasDownloadedRepositories()
	argsWithoutProg := os.Args[1:]
	command := strings.ToLower(GetCommand(argsWithoutProg))

	switch command {
	case "list":
	case "search":
	case "langs":
		if needRepositoriesUpdate {
			update(m)
		}

	case "run":
	case "install":
		m, _ = checkInterpreterAndReinit(m, c)
	}

	runCommand(command, argsWithoutProg, m, c)
}

func runCommand(command string, args []string, m *manager.Manager, c *configurator.Configurator) {
	switch command {
	case "update":
		update(m)

	case "list":
		list(m, args)

	case "search":
		search(m, args)

	case "show":
		show(m, args)

	case "run":
		run(m, args)

	case "install":
		install(m, args)

	case "remove":
		remove(m, args)

	case "findinterpreter":
		findInterpreter(m, c)

	case "repositories":
		repositories(m)

	case "langs":
		langs(m)

	case "configpath":
		printConfigPath(c)

	case "version":
		printVersion()

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
	games, e := m.GetSortedGamesByDateDesc()
	ExitIfError(e)

	// Parse args without "list" command
	repository, lang, onlyInstalled := getGamesFilterValues(args[1:])

	if repository != nil || lang != nil || onlyInstalled {
		games = manager.FilterGames(games, nil, repository, lang, onlyInstalled)
	}

	printGames(games)
}

func search(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	ExitIfError(e)

	keyword := GetCommandArg(args)
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
	ExitIfError(e)

	keyword := GetCommandArg(args)
	if keyword == nil {
		printHelpAndExit()
	}

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	fmt.Printf("Downloading and installing game %s...", FmtName(game.Title))

	installProgress := func(size uint64) {
		percents := utils.Percents(size, uint64(game.Size))
		fmt.Printf("\rDownloading and installing game %s... %s", FmtName(game.Title), color.GreenString(percents))
	}

	e = m.InstallGame(&game, installProgress)
	ExitIfError(e)

	fmt.Printf("\nGame %s has installed.\n", FmtName(game.Title))
}

func show(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	ExitIfError(e)

	keyword := GetCommandArg(args)
	if keyword == nil {
		printHelpAndExit()
	}

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	installedTxt := ""
	if game.Installed {
		installedTxt = FmtInstalled("[installed]")
	}

	// Print game information
	fmt.Printf(
		"%s (%s) %s %s\n",
		FmtTitle(game.Title), FmtName(game.Name), FmtSize(game.HumanSize()), installedTxt)
	fmt.Printf("Version: %s\n", FmtVersion(game.HumanVersion()))
	if game.Languages != nil {
		fmt.Printf("Languages: %s\n", FmtLang(strings.Join(game.Languages, ", ")))
	}
	if game.RepositoryName != "" {
		fmt.Printf("Repository: %s\n", FmtRepo(game.RepositoryName))
	}
	if game.Descurl != "" {
		fmt.Printf("More: %s\n", FmtURL(game.Descurl))
	}
	if game.Description != "" {
		fmt.Printf("\n"+color.New(color.Bold).Sprint("Descriprion")+":\n%s\n", game.Description)
	}
}

func run(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	ExitIfError(e)

	keyword := GetCommandArg(args)
	if keyword == nil {
		printHelpAndExit()
	}

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	if !game.Installed {
		fmt.Printf("Game %s isn't installed.\n", FmtName(game.Title))
		fmt.Printf("Please run for installation:\n"+
			"insteadman install %s\n", game.Name)
		os.Exit(1)
	}

	e = m.RunGame(&game)
	ExitIfError(e)

	fmt.Printf("Running %s game...\n", FmtName(game.Title))
}

func remove(m *manager.Manager, args []string) {
	games, e := m.GetSortedGames()
	ExitIfError(e)

	keyword := GetCommandArg(args)
	if keyword == nil {
		printHelpAndExit()
	}

	filteredGames := manager.FilterGames(games, keyword, nil, nil, false)

	game := getOrExitIfNoGame(filteredGames, *keyword)

	fmt.Printf("Removing game %s...\n", FmtName(game.Title))

	e = m.RemoveGame(&game)
	ExitIfError(e)

	fmt.Printf("Game %s has removed.\n", FmtName(game.Title))
}

func findInterpreter(m *manager.Manager, c *configurator.Configurator) {
	path := m.InterpreterFinder.Find()

	if path == nil {
		fmt.Println("INSTEAD has not found. Please add it in config.yml (interpreter_command)")
		return
	}

	fmt.Printf("INSTEAD has found: %s\n", *path)

	m.Config.InterpreterCommand = *path
	e := c.SaveConfig(m.Config)
	ExitIfError(e)

	fmt.Println("Path has saved")
}

func repositories(m *manager.Manager) {
	for _, repo := range m.GetRepositories() {
		fmt.Printf("%s (%s)\n", FmtRepo(repo.Name), repo.Url)
	}
}

func langs(m *manager.Manager) {
	games, e := m.GetSortedGames()
	ExitIfError(e)

	for _, lang := range m.FindLangs(games) {
		fmt.Printf("%s\n", FmtLang(lang))
	}
}

func printVersion() {
	fmt.Println(version)
}

func printConfigPath(c *configurator.Configurator) {
	fmt.Println(c.FilePath)
}

func printHelpAndExit() {
	asciiArt := `
 ___           _                 _ __  __             
|_ _|_ __  ___| |_ ___  __ _  __| |  \/  | __ _ _ __  
 | || '_ \/ __| __/ _ \/ _, |/ _, | |\/| |/ _, | '_ \ 
 | || | | \__ \ ||  __/ (_| | (_| | |  | | (_| | | | |
|___|_| |_|___/\__\___|\__,_|\__,_|_|  |_|\__,_|_| |_|
`

	color.Cyan(asciiArt)
	fmt.Printf("\n"+color.New(color.Bold).Sprint("InsteadMan CLI")+" %s — INSTEAD games manager (launcher)\n\n", version)
	fmt.Print(color.New(color.FgCyan, color.Bold).Sprint("Usage") + ":\n" +
		"    insteadman-cli [command] [keyword]\n\n" +

		color.New(color.FgCyan, color.Bold).Sprint("Commands") + ":\n" +

		color.New(color.FgCyan, color.Bold).Sprint("update") +
		"\n    Update game's repositories\n" +

		color.New(color.FgCyan, color.Bold).Sprint("list") + color.CyanString(" --repo=[name] --lang=[lang] --installed") +
		"\n    Print list of games with filtering\n" +

		color.New(color.FgCyan, color.Bold).Sprint("search") + color.CyanString(" [keyword] --repo=[name] --lang=[lang] --installed") +
		"\n    Search game by name and title with filtering\n" +

		color.New(color.FgCyan, color.Bold).Sprint("show") + color.CyanString(" [keyword]") +
		"\n    Show information about game by keyword\n" +

		color.New(color.FgCyan, color.Bold).Sprint("install") + color.CyanString(" [keyword]") +
		"\n    Install game by keyword\n" +

		color.New(color.FgCyan, color.Bold).Sprint("run") + color.CyanString(" [keyword]") +
		"\n    Run game by keyword\n" +

		color.New(color.FgCyan, color.Bold).Sprint("remove") + color.CyanString(" [keyword]") +
		"\n    Remove game by keyword\n" +

		color.New(color.FgCyan, color.Bold).Sprint("findInterpreter") +
		"\n    Find INSTEAD interpreter and save path to the config\n" +

		color.New(color.FgCyan, color.Bold).Sprint("repositories") +
		"\n    Print available repositories\n" +

		color.New(color.FgCyan, color.Bold).Sprint("langs") +
		"\n    Print available game languages\n" +

		color.New(color.FgCyan, color.Bold).Sprint("configPath") +
		"\n    Print config path\n" +

		color.New(color.FgCyan, color.Bold).Sprint("version") +
		"\n    Print current version of the application\n\n" +

		"More info: " + FmtURL("http://jhekasoft.github.io/insteadman/") + "\n")
	os.Exit(1)
}

// -- Commands -----------------------------------

func initManagerAndConfigurator() (*manager.Manager, *configurator.Configurator) {
	executablePath, e := os.Executable()
	ExitIfError(e)

	currentDir, e := utils.BinAbsDir(executablePath)
	ExitIfError(e)

	c := configurator.Configurator{FilePath: "", CurrentDir: currentDir, Version: version}
	config, e := c.GetConfig()
	ExitIfError(e)

	finder := &interpreterfinder.InterpreterFinder{CurrentDir: currentDir}

	m := manager.Manager{Config: config, InterpreterFinder: finder}

	return &m, &c
}

func checkInterpreterAndReinit(m *manager.Manager, c *configurator.Configurator) (*manager.Manager, *configurator.Configurator) {
	if m.InterpreterCommand() == "" {
		findInterpreter(m, c)
		m, c = initManagerAndConfigurator()
	}

	return m, c
}

func printGames(games []manager.Game) {
	for _, game := range games {
		installed := ""
		if game.Installed {
			installed = FmtInstalled("[installed]")
		}
		fmt.Printf(
			"%s, %s, %s "+FmtLang("%v")+" %s\n",
			FmtTitle(game.Title), FmtName(game.Name), FmtRepo(game.RepositoryName), game.Languages, installed)

	}
}

func getOrExitIfNoGame(filteredGames []manager.Game, keyword string) manager.Game {
	if len(filteredGames) < 1 {
		fmt.Printf("Game %s has not found\n", FmtName(keyword))
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
	repository := FindStringArg("--repository", args)
	lang := FindStringArg("--lang", args)
	onlyInstalled := FindBoolArg("--installed", args)

	return repository, lang, onlyInstalled
}
