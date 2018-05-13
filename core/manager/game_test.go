package manager

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

const (
	gamesPath          = "../../resources/testdata/games"
	stead2gameFileName = "stead2testgame"
	stead3gameFileName = "stead3testgame"
)

func TestReadLocalGameInfoStead2(t *testing.T) {
	stead2GameInfo, e := os.Stat(filepath.Join(gamesPath, stead2gameFileName))
	assert.NoError(t, e)

	stead2Game := ReadLocalGameInfo(gamesPath, stead2GameInfo)

	assert.Equal(t, "stead2testgame", stead2Game.Name)
	assert.Equal(t, "Stead2 Test game", stead2Game.Title)
	assert.Equal(t, "0.4", stead2Game.InstalledVersion)
	assert.Equal(t, "0.4", stead2Game.Version)
}

func TestReadLocalGameInfoStead3(t *testing.T) {
	stead2GameInfo, e := os.Stat(filepath.Join(gamesPath, stead3gameFileName))
	assert.NoError(t, e)

	stead2Game := ReadLocalGameInfo(gamesPath, stead2GameInfo)

	assert.Equal(t, "stead3testgame", stead2Game.Name)
	assert.Equal(t, "Stead3 Test game", stead2Game.Title)
	assert.Equal(t, "0.1", stead2Game.InstalledVersion)
	assert.Equal(t, "0.1", stead2Game.Version)
}
