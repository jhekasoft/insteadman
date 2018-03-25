package manager

import (
	"../configurator"
	"../interpreter_finder"
	"github.com/stretchr/testify/assert"
	"testing"
)

const configFilePath = "../../resources/testdata/insteadman/config.yml"
const testGameName = "crossworlds"
const testGameUrl = "http://instead-games.ru/download/instead-crossworlds-0.7.zip"
const filterKeyword = "test"
const filterWrongKeyword = "878fdfd----fdsfdsftest"

func TestUpdateRepositories(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	man := Manager{Config: config}
	errors := man.UpdateRepositories()

	assert.Empty(t, errors)
}

func TestGetSortedGamesAndFilterGames(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	man := Manager{Config: config}

	// Sorted games
	games, e := man.GetSortedGames()
	assert.NoError(t, e)
	assert.NotEmpty(t, games)

	// Filter games by keyword
	keyword := filterKeyword
	filteredGames := FilterGames(games, &keyword, nil, nil, false)
	assert.NotEmpty(t, filteredGames)

	// Filter games by wrong keyword (not found games)
	wrongKeyword := filterWrongKeyword
	emptyFilteredGames := FilterGames(games, &wrongKeyword, nil, nil, false)
	assert.Empty(t, emptyFilteredGames)
}

func TestInstallGame(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	// Find installed INSTEAD and use it in config
	finder := interpreterFinder.InterpreterFinder{Config: config}
	interpreterPath := finder.Find()
	assert.NotNil(t, interpreterPath)
	config.InterpreterCommand = *interpreterPath

	man := Manager{Config: config}

	e = man.InstallGame(&Game{Name: testGameName, Url: testGameUrl})

	assert.NoError(t, e)
}

func TestRunGame(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	// Find installed INSTEAD and use it in config
	finder := interpreterFinder.InterpreterFinder{Config: config}
	interpreterPath := finder.Find()
	assert.NotNil(t, interpreterPath)
	config.InterpreterCommand = *interpreterPath

	man := Manager{Config: config}

	// Run game
	e = man.RunGame(&Game{Name: testGameName, Url: testGameUrl})
	assert.NoError(t, e)

	// Stop running
	e = man.StopRunningGame()
	assert.NoError(t, e)
}

func TestRemoveGame(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	man := Manager{Config: config}
	e = man.RemoveGame(&Game{Name: testGameName, Url: testGameUrl})

	assert.NoError(t, e)
}

func TestRepositories(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	man := Manager{Config: config}
	repositories := man.GetRepositories()

	assert.NotEmpty(t, repositories)
}

func TestLangs(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	man := Manager{Config: config}

	games, e := man.GetSortedGames()
	assert.NoError(t, e)

	langs := man.FindLangs(games)
	assert.NotEmpty(t, langs)
}
