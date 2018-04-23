package configurator

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

type InsteadmanConfig struct {
	Repositories             []Repository `json:"repositories"`
	InterpreterCommand       string       `json:"interpreter_command"`
	Version                  string       `json:"version"`
	UseBuiltinInterpreter    bool         `json:"use_builtin_interpreter"`
	Lang                     string       `json:"lang"`
	CheckUpdateOnStart       bool         `json:"check_update_on_start"`
	GamesPath                string       `json:"games_path"`
	InsteadManPath           string       `json:"insteadman_path"`
	CalculatedGamesPath      string       `json:"-"`
	CalculatedInsteadManPath string       `json:"-"`
}

func (c *InsteadmanConfig) GetInterpreterCommand() string {
	path, e := filepath.Abs(c.InterpreterCommand)
	if e != nil {
		return c.InterpreterCommand
	}

	return path
}

type Repository struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

const (
	configName        = "config.yml"
	skeletonDir       = "skeleton"
	gamesDirName      = "games"
	insteadManDirName = "insteadman"
)

type Configurator struct {
	FilePath string
}

func insteadManDir() string {
	localPath := filepath.Join(".", configName)
	_, e := os.Stat(localPath)
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return "."
	}

	insteadManDir := filepath.Join(insteadDir(), insteadManDirName)
	os.MkdirAll(insteadManDir, os.ModePerm)

	return insteadManDir
}

func findConfigFileName() string {
	return filepath.Join(insteadManDir(), configName)
}

func gamesDir() string {
	localPath := filepath.Join(".", gamesDirName)

	_, e := os.Stat(localPath)
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return localPath
	}

	gamesDir := filepath.Join(insteadDir(), gamesDirName)
	os.MkdirAll(gamesDir, os.ModePerm)

	return gamesDir
}

func writeSkeleton(c *Configurator) error {
	configData, e := ioutil.ReadFile(filepath.Join(skeletonDir, configName))
	if e != nil {
		return e
	}

	return ioutil.WriteFile(c.FilePath, configData, 0644)
}

func (c *Configurator) GetConfig() (*InsteadmanConfig, error) {
	if c.FilePath == "" {
		c.FilePath = findConfigFileName()
	}

	// Write skeleton config if it isn't existing
	_, e := os.Stat(c.FilePath)
	exists := !os.IsNotExist(e)
	if !exists || e != nil {
		e = writeSkeleton(c)
		if e != nil {
			return nil, e
		}
	}

	file, e := ioutil.ReadFile(c.FilePath)
	if e != nil {
		return nil, e
	}
	// fmt.Printf("%s\n", string(file))

	var config *InsteadmanConfig
	yaml.Unmarshal(file, &config)

	// TODO: make Calculated* fields like GetInterpreterCommand() func, but like "lazy vars"

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
