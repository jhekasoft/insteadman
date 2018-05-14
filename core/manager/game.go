package manager

import (
	"../utils"
	"github.com/pyk/byten"
	"html"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
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
	Date             string   `xml:"date"`
	Timestamp        int64    `xml:"-"`
	InstalledVersion string   `xml:"-"`
	RepositoryName   string   `xml:"-"`
	Installed        bool     `xml:"-"`
	OnlyInstalled    bool     `xml:"-"`
	//IsUpdateExist    bool     `xml:"-"`
	Languages []string `xml:"-"`
	Id        string   `xml:"-"`
}

type Game RepositoryGame

func generateGameId(repository string, g *Game) string {
	return repository + "/" + g.Name + "/" + strings.Join(g.Languages, "_")
}

func (g *Game) addGameAdditionalData(repositoryName string) {
	date, e := time.Parse("2006-01-02", g.Date)
	if e == nil {
		g.Timestamp = date.Unix()
	}

	if len(g.Langs) > 0 {
		g.Languages = g.Langs
	} else {
		g.Languages = strings.Split(g.Lang, ",")
	}

	g.RepositoryName = repositoryName

	g.Title = html.UnescapeString(g.Title)

	if g.Description != "" {
		g.Description = html.UnescapeString(g.Description)
	}

	g.Id = generateGameId(repositoryName, g)
}

func (g *Game) HumanSize() string {
	if g.Size > 0 {
		return byten.Size(int64(g.Size))
	}

	return ""
}

func (g *Game) HumanVersion() string {
	if g.IsUpdateAvailable() {
		return g.InstalledVersion + " (" + g.Version + ")"
	}

	return g.Version
}

func (g *Game) IsUpdateAvailable() bool {
	return g.InstalledVersion != "" && g.InstalledVersion != g.Version
}

func ReadLocalGameInfo(path string, info os.FileInfo) Game {
	var e error

	// Eval possible symlink
	gameName := info.Name()
	evalPath, e := filepath.EvalSymlinks(filepath.Join(path, info.Name()))
	if e == nil {
		newInfo, e := os.Stat(evalPath)
		if e == nil {
			info = newInfo
		}
	}

	if !info.IsDir() { // IDF
		gameName = strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
	}

	game := Game{
		Name:      gameName,
		Title:     gameName,
		Installed: true,
	}

	if info.IsDir() {
		game, _ = appendGameInfo(game, path, info)
	}

	game.addGameAdditionalData("")

	return game
}

func appendGameInfo(game Game, path string, info os.FileInfo) (newGame Game, e error) {
	// TODO: idf

	newGame = game

	mainLuaFilePath := filepath.Join(path, info.Name(), "main.lua") // STEAD2 main file path

	if !utils.PathExist(mainLuaFilePath) {
		mainLuaFilePath = filepath.Join(path, info.Name(), "main3.lua") // STEAD3 main file path
	}

	file, e := ioutil.ReadFile(mainLuaFilePath)
	if e != nil {
		return
	}

	// Title
	r, e := regexp.Compile("(?i)--\\s*\\$Name:\\s*(.*)\\$")
	if e == nil {
		matches := r.FindStringSubmatch(string(file))
		if len(matches) > 1 {
			newGame.Title = matches[1]
		}
	}

	// Version
	r, e = regexp.Compile("(?i)--\\s*\\$Version:\\s*(.*)\\$")
	if e == nil {
		matches := r.FindStringSubmatch(string(file))
		if len(matches) > 1 {
			newGame.InstalledVersion = matches[1]
			newGame.Version = matches[1]
		}
	}

	return
}
