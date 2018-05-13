package manager

import (
	"../utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func readLocalGameInfo(path string, info os.FileInfo) Game {
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
