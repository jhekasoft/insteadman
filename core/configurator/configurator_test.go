package configurator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const configFilePath = "../../resources/testdata/insteadman/config.yml"

func TestGetConfig(t *testing.T) {
	configurator := Configurator{FilePath: configFilePath}
	config, e := configurator.GetConfig()

	assert.NoError(t, e)
	assert.NotEmpty(t, config.CalculatedGamesPath)
}

func TestSaveConfig(t *testing.T) {
	configurator := Configurator{FilePath: configFilePath}
	config, e := configurator.GetConfig()

	assert.NoError(t, e)
	assert.NotEmpty(t, config.CalculatedGamesPath)

	e = configurator.SaveConfig(config)
	assert.NoError(t, e)
}

func TestGetSkeletonConfig(t *testing.T) {
	configurator := Configurator{FilePath: configFilePath, CurrentDir: "../../"}
	config, e := configurator.GetSkeletonConfig()

	assert.NoError(t, e)
	assert.NotEmpty(t, config.Repositories)
}
