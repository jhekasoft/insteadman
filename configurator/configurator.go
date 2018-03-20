package configurator

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type InsteadmanConfig struct {
	Repositories          []Repository `json:"repositories"`
	InterpreterCommand    string       `json:"interpreter_command"`
	Version               string       `json:"version"`
	UseBuiltinInterpreter bool         `json:"use_builtin_interpreter"`
	Lang                  string       `json:"lang"`
	CheckUpdateOnStart    bool         `json:"check_update_on_start"`
	GamesPath             string       `json:"games_path"`
	InsteadManPath        string       `json:"insteadman_path"`
	CalculatedGamesPath   string       `json:"-"`
	CalculatedInsteadManPath   string       `json:"-"`
}

type Repository struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

const configName = "config.yml"

type Configurator struct {
	FilePath string
}

func insteadDir() string {
	homeDir := ""

	u, e := user.Current()
	if e == nil {
		homeDir = u.HomeDir
	}

	return filepath.Join(homeDir, ".instead")
}

func insteadManDir() string {
	localPath := filepath.Join(".", configName)
	_, e := os.Stat(localPath)
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return "."
	}

	return filepath.Join(insteadDir(), "insteadman")
}

func findConfigFileName() string {
	return filepath.Join(insteadManDir(), configName)
}

func gamesDir() string {
	localPath := filepath.Join(".", "games")

	_, e := os.Stat(localPath)
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return localPath
	}

	return filepath.Join(insteadDir(), "games")
}

func (c *Configurator) GetConfig() (*InsteadmanConfig, error) {
	if c.FilePath == "" {
		c.FilePath = findConfigFileName()
	}

	file, e := ioutil.ReadFile(c.FilePath)
	if e != nil {
		return nil, e
	}
	// fmt.Printf("%s\n", string(file))

	var config *InsteadmanConfig
	yaml.Unmarshal(file, &config)

	config.CalculatedGamesPath = config.GamesPath
	if config.CalculatedGamesPath == "" {
		config.CalculatedGamesPath = gamesDir()
	}

	config.CalculatedInsteadManPath = config.InsteadManPath
	if config.CalculatedInsteadManPath == "" {
		config.CalculatedInsteadManPath = insteadManDir()
	}

	return config, nil
}

func (c *Configurator) SaveConfig(config *InsteadmanConfig) error {
	 bytes, e := yaml.Marshal(config)
	 if e != nil {
	     return e
	 }

	 return ioutil.WriteFile(c.FilePath, bytes, 0644)
}
