package manager

import (
	"testing"

	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/interpreterfinder"
	"github.com/stretchr/testify/assert"
)

const (
	configFilePath     = "../../resources/testdata/insteadman/config.yml"
	testGameName       = "crossworlds"
	testGameUrl        = "http://instead-games.ru/download/instead-crossworlds-0.7.zip"
	testGameId         = "instead-games/crossworlds"
	testGameImage      = "http://instead-games.ru/games/screenshots/20130204203941245.png"
	filterKeyword      = "test"
	filterWrongKeyword = "878fdfd----fdsfdsftest"
)

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
	finder := new(interpreterfinder.InterpreterFinder)
	interpreterPath := finder.Find()
	assert.NotNil(t, interpreterPath)
	config.InterpreterCommand = *interpreterPath

	man := Manager{Config: config, InterpreterFinder: finder}

	e = man.InstallGame(&Game{Name: testGameName, Url: testGameUrl}, nil)

	assert.NoError(t, e)
}

func TestFindGameById(t *testing.T) {
	games := []Game{
		{Id: "official/game1"},
		{Id: "official/game2"},
	}

	assert.NotNil(t, FindGameById(games, "official/game2"))
	assert.Nil(t, FindGameById(games, "fdfdfdfd"))
}

func TestRunGame(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	// Find installed INSTEAD and use it in config
	finder := new(interpreterfinder.InterpreterFinder)
	interpreterPath := finder.Find()
	assert.NotNil(t, interpreterPath)
	config.InterpreterCommand = *interpreterPath

	man := Manager{Config: config, InterpreterFinder: finder}

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

func TestGetGameImage(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	// Find installed INSTEAD and use it in config
	finder := new(interpreterfinder.InterpreterFinder)
	interpreterPath := finder.Find()
	assert.NotNil(t, interpreterPath)
	config.InterpreterCommand = *interpreterPath

	man := Manager{Config: config, InterpreterFinder: finder}

	imageFilePath, e := man.GetGameImage(&Game{Id: testGameId, Image: testGameImage})

	assert.NoError(t, e)
	assert.NotEmpty(t, imageFilePath)
}

func TestClearCache(t *testing.T) {
	conf := configurator.Configurator{FilePath: configFilePath}
	config, e := conf.GetConfig()
	assert.NoError(t, e)

	man := Manager{Config: config}

	e = man.ClearCache()
	assert.NoError(t, e)
}

func TestFilterRepositoryName(t *testing.T) {
	names := map[string]string{
		"test/test12":      "testtest12",
		"TesT.///2_gg":     "TesT.2_gg",
		"TestПривет/2Пока": "Test2",
	}

	for name, mustBeName := range names {
		result, e := FilterRepositoryName(name)
		assert.NoError(t, e)
		assert.Equal(t, result, mustBeName)
	}
}
