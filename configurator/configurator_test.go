package configurator

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const configFilePath = "../resources/testdata/insteadman/config.yml"

func TestGetConfig(t *testing.T) {
	configurator := Configurator{FilePath: configFilePath}
	config, e := configurator.GetConfig()
	expectedGamePath := "../games"

	assert.NoError(t, e)
	assert.Equal(t, config.GamesPath, expectedGamePath)
}
