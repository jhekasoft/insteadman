package manager

import (
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jhekasoft/insteadman/core/configurator"
	"github.com/jhekasoft/insteadman/core/interpreterfinder"
	"github.com/jhekasoft/insteadman/core/utils"
)

const (
	//updateCheckUrl = "https://raw.githubusercontent.com/jhekasoft/insteadman/master/version.json"
	cacheDirName        = "cache"
	repositoriesDirName = "repositories"
	tempGamesDirName    = "temp_games"
	gameImagesDirName   = "game_images"

	SortByTitleAsc = "title"
	SortByDateDesc = "date"
)

type Manager struct {
	Config            *configurator.InsteadmanConfig
	InterpreterFinder *interpreterfinder.InterpreterFinder
	CurrentRunningCmd *exec.Cmd
}

func (m *Manager) HasDownloadedRepositories() bool {
	repositoriesDir := m.repositoriesDir()
	os.MkdirAll(repositoriesDir, os.ModePerm)

	files, e := filepath.Glob(filepath.Join(repositoriesDir, "*.xml"))
	if e != nil || files == nil {
		return false
	}

	return true
}

func (m *Manager) UpdateRepositories() []error {
	repositoriesDir := m.repositoriesDir()
	os.MkdirAll(repositoriesDir, os.ModePerm)

	// Remove all repository files
	files, e := filepath.Glob(filepath.Join(repositoriesDir, "*.xml"))
	if e == nil && files != nil {
		for _, f := range files {
			os.Remove(f)
		}
	}

	var errs []error = nil
	for _, repo := range m.Config.Repositories {
		e := downloadFileSimple(filepath.Join(repositoriesDir, repo.Name+".xml"), repo.Url)

		if e != nil {
			errs = append(errs, e)
		}
	}

	return errs
}

func (m *Manager) GetRepositoryGames() ([]Game, error) {
	repositoriesDir := m.repositoriesDir()
	files, e := filepath.Glob(filepath.Join(repositoriesDir, "*.xml"))
	if e != nil {
		return nil, e
	}

	var games []Game = nil
	for _, fileName := range files {
		// fmt.Printf("File: %v\n", fileName)

		gameList, e := parseRepository(fileName)
		if e == nil {
			repositoryFileName := filepath.Base(fileName)
			repositoryName := strings.TrimSuffix(repositoryFileName, filepath.Ext(repositoryFileName))

			var repositoryGames []Game = nil
			for _, repositoryGame := range gameList.GameList {
				game := Game(repositoryGame)
				game.addGameAdditionalData(repositoryName)
				repositoryGames = append(repositoryGames, game)
			}
			games = append(games, repositoryGames...)
		}
	}

	return games, nil
}

func (m *Manager) CacheDir() string {
	return filepath.Join(m.Config.CalculatedInsteadManPath, cacheDirName)
}

func (m *Manager) repositoriesDir() string {
	return filepath.Join(m.Config.CalculatedInsteadManPath, cacheDirName, repositoriesDirName)
}

func (m *Manager) gameImagesDir() string {
	return filepath.Join(m.Config.CalculatedInsteadManPath, cacheDirName, gameImagesDirName)
}

func parseRepository(fileName string) (*RepositoryGameList, error) {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		return nil, e
	}

	var gameList *RepositoryGameList
	e = xml.Unmarshal(file, &gameList)

	if e != nil {
		return nil, e
	}

	return gameList, nil
}

func (m *Manager) GetInstalledGames() ([]Game, error) {
	files, e := ioutil.ReadDir(m.Config.CalculatedGamesPath)
	if e != nil {
		return nil, e
	}

	var games []Game = nil
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		game := ReadLocalGameInfo(m.Config.CalculatedGamesPath, file)
		games = append(games, game)
	}

	return games, nil
}

func (m *Manager) GetMergedGames() ([]Game, error) {
	games, e := m.GetRepositoryGames()
	if e != nil {
		return nil, e
	}

	installedGames, e := m.GetInstalledGames()
	if e != nil {
		return nil, e
	}

	for i := range installedGames {
		installedGames[i].OnlyInstalled = true
	}

	for i, game := range games {
		for j, installedGame := range installedGames {
			if game.Name == installedGame.Name {
				games[i].Installed = true
				games[i].InstalledVersion = installedGame.InstalledVersion
				installedGames[j].OnlyInstalled = false
			}
		}
	}

	for _, installedGame := range installedGames {
		if installedGame.OnlyInstalled {
			games = append(games, installedGame)
		}
	}

	return games, nil
}

func (m *Manager) GetSortedGames() ([]Game, error) {
	return m.GetSortedGamesBy(SortByTitleAsc)
}

func (m *Manager) GetSortedGamesByDateDesc() ([]Game, error) {
	return m.GetSortedGamesBy(SortByDateDesc)
}

//type GameTitleSorter []Game
//
//func (a GameTitleSorter) Len() int           { return len(a) }
//func (a GameTitleSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a GameTitleSorter) Less(i, j int) bool { return a[i].Title < a[j].Title }

func (m *Manager) GetSortedGamesBy(sortBy string) ([]Game, error) {
	games, e := m.GetMergedGames()

	if e != nil {
		return nil, e
	}

	switch sortBy {
	case SortByTitleAsc:
		//sort.Sort(GameTitleSorter(games))
		sort.Slice(games, func(i, j int) bool {
			return strings.ToLower(games[i].Title) < strings.ToLower(games[j].Title)
		})
	case SortByDateDesc:
		sort.Slice(games, func(i, j int) bool {
			return games[i].Timestamp > games[j].Timestamp
		})
	}

	return games, nil
}

func FilterGames(games []Game, keyword *string, repository *string, lang *string, onlyInstalled bool) []Game {
	if onlyInstalled {
		games = filterGamesBy(games, func(game Game) bool {
			return game.Installed == true
		})
	}

	if repository != nil {
		games = filterGamesBy(games, func(game Game) bool {
			return game.RepositoryName == *repository
		})
	}

	if lang != nil {
		games = filterGamesBy(games, func(game Game) bool {
			return utils.ExistsString(game.Languages, *lang)
		})
	}

	if keyword != nil {
		lowerKeyword := strings.ToLower(*keyword)

		games = filterGamesBy(games, func(game Game) bool {
			return strings.Contains(strings.ToLower(game.Title), lowerKeyword) ||
				strings.Contains(strings.ToLower(game.Name), lowerKeyword)
		})
	}

	return games
}

func filterGamesBy(games []Game, f func(Game) bool) []Game {
	gamesFiltered := make([]Game, 0)
	for _, game := range games {
		if f(game) {
			gamesFiltered = append(gamesFiltered, game)
		}
	}
	return gamesFiltered
}

func FindGameById(games []Game, id string) *Game {
	for _, game := range games {
		if game.Id == id {
			return &game
		}
	}

	return nil
}

func FindGamesByName(games []Game, name string) (foundGames []Game) {
	for _, game := range games {
		if game.Name == name {
			foundGames = append(foundGames, game)
		}
	}

	return
}

func (m *Manager) RunGame(game *Game) error {
	if game == nil {
		return nil
	}

	// Absolute games path
	gamesPath, e := filepath.Abs(m.Config.CalculatedGamesPath)
	if e != nil {
		return e
	}

	interpreterCommand := m.InterpreterCommand()

	// todo: idf
	cmd := exec.Command(interpreterCommand, "-gamespath", gamesPath, "-game", game.Name)
	cmd.Dir = filepath.Dir(interpreterCommand)
	e = cmd.Start()

	// Current running cmd
	if e == nil {
		m.CurrentRunningCmd = cmd
	}

	return e
}

func (m *Manager) StopRunningGame() error {
	if m.CurrentRunningCmd == nil {
		return nil
	}

	e := m.CurrentRunningCmd.Process.Kill()

	return e
}

func downloadFileSimple(fileName, url string) error {
	// Create the file
	out, e := os.Create(fileName)
	if e != nil {
		return e
	}
	defer out.Close()

	// Download the data
	resp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	// Write the data to the file
	_, e = io.Copy(out, resp.Body)

	return e
}

type WriteCounter struct {
	Total     uint64
	progressF func(uint64) // progress function
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	if wc.progressF != nil {
		wc.progressF(wc.Total)
	}
	return n, nil
}

func downloadFile(fileName, url string, progressF func(uint64)) error {
	// Create the file
	out, e := os.Create(fileName)
	if e != nil {
		return e
	}
	defer out.Close()

	// Download the data
	resp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	counter := &WriteCounter{progressF: progressF}
	_, e = io.Copy(out, io.TeeReader(resp.Body, counter))

	return e
}

func (m *Manager) GetGameImage(game *Game) (imagePath string, e error) {
	if game == nil || game.Image == "" || game.Id == "" {
		return
	}

	// Processing image only if URL have extension
	imageExt := filepath.Ext(game.Image)
	if imageExt == "" {
		return
	}

	gameImagesDir := m.gameImagesDir()
	os.MkdirAll(gameImagesDir, os.ModePerm)

	fileName := strings.Replace(game.Id, "/", "_", -1) + imageExt
	imagePath = filepath.Join(gameImagesDir, fileName)

	_, e = os.Stat(imagePath)
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return imagePath, e
	}

	e = downloadFileSimple(imagePath, game.Image)
	if e != nil {
		return "", e
	}

	return imagePath, e
}

func (m *Manager) InstallGame(game *Game, progressF func(uint64)) error {
	// todo: idf

	tempGamesDir := filepath.Join(m.CacheDir(), tempGamesDirName)
	os.MkdirAll(tempGamesDir, os.ModePerm)

	// Absolute filepath
	fileName := filepath.Join(tempGamesDir, path.Base(game.Url))
	fileNameAbs, e := filepath.Abs(fileName)
	if e == nil {
		fileName = fileNameAbs
	}

	e = downloadFile(fileName, game.Url, progressF)
	if e != nil {
		return e
	}

	// Remove downloaded temp file (after installing)
	defer os.Remove(fileName)

	// Absolute games path
	gamesPath, e := filepath.Abs(m.Config.CalculatedGamesPath)
	if e != nil {
		return e
	}

	interpreterCommand := m.InterpreterCommand()

	cmd := exec.Command(interpreterCommand, "-gamespath", gamesPath, "-install", fileName, "-quit")
	cmd.Dir = filepath.Dir(interpreterCommand)
	out, e := cmd.CombinedOutput()
	if e != nil {
		return errors.New(e.Error() + "; " + strings.Replace(string(out), "\n", "", -1))
	}

	return nil
}

func (m *Manager) RemoveGame(game *Game) error {
	// todo: idf

	gameDir := filepath.Join(m.Config.CalculatedGamesPath, game.Name)

	e := os.RemoveAll(gameDir)

	return e
}

func (m *Manager) GetRepositories() []configurator.Repository {
	return m.Config.Repositories
}

func (m *Manager) FindLangs(games []Game) []string {
	var langs []string = nil

	for _, game := range games {
		for _, gameLang := range game.Languages {
			if !utils.ExistsString(langs, gameLang) && strings.Trim(gameLang, " ") != "" {
				langs = append(langs, gameLang)
			}
		}
	}

	return langs
}

func (m *Manager) ClearCache() error {
	return os.RemoveAll(m.CacheDir())
}

func (m *Manager) IsBuiltinInterpreterCommand() bool {
	if m.Config.UseBuiltinInterpreter {
		return m.InterpreterFinder.HaveBuiltIn()
	}

	return false
}

func (m *Manager) InterpreterCommand() string {
	if m.Config.UseBuiltinInterpreter {
		builtInCmd := m.InterpreterFinder.FindBuiltIn()
		if builtInCmd != "" {
			return configurator.ExpandInterpreterCommand(builtInCmd)
		}
	}

	if m.Config.InterpreterCommand != "" {
		return configurator.ExpandInterpreterCommand(m.Config.InterpreterCommand)
	}

	return ""
}

func FilterRepositoryName(name string) (filteredName string, e error) {
	r, e := regexp.Compile("[^a-zA-Z0-9\\-_.]+")
	if e != nil {
		return
	}

	filteredName = r.ReplaceAllString(name, "")
	return
}

// func CheckAppNewVersion() {

// }
