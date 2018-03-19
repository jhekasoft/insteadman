package manager

import (
	"../configurator"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type RepositoryGameList struct {
	// XMLName xml.Name `xml:"game_list"`
	GameList []RepositoryGame `xml:"game"`
}

type RepositoryGame struct {
	// XMLName xml.Name `xml:"game"`
	Name             string   `xml:"name"`
	Title            string   `xml:"title"`
	Version          string   `xml:"version"`
	Url              string   `xml:"url"`
	Size             int      `xml:"size"`
	Lang             string   `xml:"lang"`
	Descurl          string   `xml:"descurl"`
	Author           string   `xml:"author"`
	Description      string   `xml:"description"`
	Image            string   `xml:"image"`
	Langs            []string `xml:"langs>lang"`
	InstalledVersion string   `xml:"-"`
	RepositoryName   string   `xml:"-"`
	Installed        bool     `xml:"-"`
	OnlyInstalled    bool     `xml:"-"`
	IsUpdateExist    bool     `xml:"-"`
	Languages        []string `xml:"-"`
	Id               string   `xml:"-"`
}

type Game RepositoryGame

func (g *Game) addGameAdditionalData(repositoryName string) {
	g.Id = repositoryName + "/" + g.Name

	if len(g.Langs) > 0 {
		g.Languages = g.Langs
	} else {
		g.Languages = strings.Split(g.Lang, ",")
	}

	g.RepositoryName = repositoryName
}

const (
	//updateCheckUrl = "https://raw.githubusercontent.com/jhekasoft/insteadman/master/version.json"
	repositoriesDirName = "repositories"
	tempGamesDirName    = "temp_games"
)

func downloadRepository(fileName, url string) error {
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
	if e != nil {
		return e
	}

	return nil
}

type Manager struct {
	Config            *configurator.InsteadmanConfig
	CurrentRunningCmd *exec.Cmd
}

func (m *Manager) DownloadRepositories() []error {
	repositoriesDir := filepath.Join(m.Config.InsteadManPath, repositoriesDirName)
	os.MkdirAll(repositoriesDir, os.ModePerm)

	var errors []error = nil
	for _, repo := range m.Config.Repositories {
		// fmt.Printf("%v %v\n", repo.Name, repo.Url)
		e := downloadRepository(filepath.Join(repositoriesDir, repo.Name+".xml"), repo.Url)

		if e != nil {
			errors = append(errors, e)
		}
	}

	return errors
}

func (m *Manager) GetRepositoryGames() ([]Game, error) {
	repositoriesDir := filepath.Join(m.Config.InsteadManPath, repositoriesDirName)
	files, e := filepath.Glob(filepath.Join(repositoriesDir, "*.xml"))
	if e != nil {
		return nil, e
	}

	var games []Game = nil
	for _, fileName := range files {
		// fmt.Printf("File: %v\n", fileName)

		gameList, e := parseRepository(filepath.Join(".", fileName))
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
	files, e := ioutil.ReadDir(m.Config.GamesPath)
	if e != nil {
		return nil, e
	}

	var games []Game = nil
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		gameName := file.Name()
		if !file.IsDir() {
			gameName = strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		}

		games = append(games, Game{
			Name:      gameName,
			Title:     gameName,
			Installed: true,
		})
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
				// todo: installed version
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

//type GameTitleSorter []Game
//
//func (a GameTitleSorter) Len() int           { return len(a) }
//func (a GameTitleSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a GameTitleSorter) Less(i, j int) bool { return a[i].Title < a[j].Title }

func (m *Manager) GetSortedGames() ([]Game, error) {
	games, e := m.GetMergedGames()

	if e != nil {
		return nil, e
	}

	//sort.Sort(GameTitleSorter(games))
	sort.Slice(games, func(i, j int) bool {
		return games[i].Title < games[j].Title
	})

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
			for _, l := range game.Languages {
				if l == *lang {
					return true
				}
			}
			return false
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

func (m *Manager) RunGame(game *Game) error {
	if game == nil {
		return nil
	}

	// Absolute games path
	gamesPath, e := filepath.Abs(m.Config.GamesPath)
	if e != nil {
		return e
	}

	// todo: idf
	cmd := exec.Command(m.Config.InterpreterCommand, "-gamespath", gamesPath, "-game", game.Name)
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

func (m *Manager) InstallGame(game *Game) error {
	// todo: idf

	tempGamesDir := filepath.Join(m.Config.InsteadManPath, tempGamesDirName)
	os.MkdirAll(tempGamesDir, os.ModePerm)

	// Absolute filepath
	fileName := filepath.Join(tempGamesDir, path.Base(game.Url))
	fileNameAbs, e := filepath.Abs(fileName)
	if e == nil {
		fileName = fileNameAbs
	}

	// Create the file
	out, e := os.Create(fileName)
	if e != nil {
		return e
	}
	defer out.Close()

	// Download the data
	resp, e := http.Get(game.Url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	// Write the data to the file
	_, e = io.Copy(out, resp.Body)
	if e != nil {
		return e
	}

	// Absolute games path
	gamesPath, e := filepath.Abs(m.Config.GamesPath)
	if e != nil {
		return e
	}

	e = exec.Command(m.Config.InterpreterCommand, "-gamespath", gamesPath, "-install", fileName, "-quit").Run()
	if e != nil {
		return e
	}

	// Remove downloaded temp file
	e = os.Remove(fileName)
	if e != nil {
		return e
	}

	return nil
}

func (m *Manager) RemoveGame(game *Game) error {
	// todo: idf

	gameDir := filepath.Join(m.Config.GamesPath, game.Name)

	e := os.RemoveAll(gameDir)

	return e
}

// func CheckAppNewVersion() {

// }
