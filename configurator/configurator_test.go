package configurator

import "testing"

const configFilePath = "../resources/testdata/insteadman/config.yml"

func TestGetConfig(t *testing.T) {
	configurator := Configurator{FilePath: configFilePath}
	config, e := configurator.GetConfig()

	if e != nil {
		t.Error("Expected no errors but got ", e)
	}

	expectedGamePath := "../games"
	if config.GamesPath != expectedGamePath {
		t.Fatalf("Expected %s but got %s", expectedGamePath, config.GamesPath)
	}
}
