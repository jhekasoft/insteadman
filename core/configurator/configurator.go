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
	Gtk                      Gtk          `json:"gtk"`
	CalculatedGamesPath      string       `json:"-"`
	CalculatedInsteadManPath string       `json:"-"`
}

func ExpandInterpreterCommand(command string) string {
	if command == "" {
		return ""
	}

	path, e := filepath.Abs(command)
	if e != nil {
		return command
	}

	return path
}

type Repository struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Gtk struct {
	HideSidebar bool `json:"hide_sidebar"`
	MainWidth   int  `json:"main_width"`
	MainHeight  int  `json:"main_height"`
}

const (
	configName        = "config.yml"
	skeletonDir       = "skeleton"
	gamesDirName      = "games"
	insteadManDirName = "insteadman"
)

type Configurator struct {
	FilePath   string
	CurrentDir string
}

func (c *Configurator) insteadManDir() string {
	localPath := filepath.Join(c.CurrentDir, configName)
	_, e := os.Stat(localPath)
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return c.CurrentDir
	}

	insteadManDir := filepath.Join(insteadDir(), insteadManDirName)
	os.MkdirAll(insteadManDir, os.ModePerm)

	return insteadManDir
}

func (c *Configurator) findConfigFileName() string {
	return filepath.Join(c.insteadManDir(), configName)
}

func (c *Configurator) gamesDir() string {
	localPath := filepath.Join(c.CurrentDir, gamesDirName)

	_, e := os.Stat(localPath)
	exists := !os.IsNotExist(e)

	if exists && e == nil {
		return localPath
	}

	gamesDir := filepath.Join(insteadDir(), gamesDirName)
	os.MkdirAll(gamesDir, os.ModePerm)

	return gamesDir
}

func (c *Configurator) ShareResourcePath(relPath string) string {
	// Add curent dir to search
	sharePathList := []string{c.CurrentDir}

	// Add UNIX-path to search
	const unixSharePath = "share/insteadman"
	unixSharedDir, e := filepath.Abs(filepath.Join(c.CurrentDir, "..", unixSharePath))
	if e == nil {
		sharePathList = append(sharePathList, unixSharedDir)
	}

	// Search resource in all the path
	for _, sharePath := range sharePathList {
		absPath := filepath.Join(sharePath, relPath)

		_, e := os.Stat(absPath)
		exists := !os.IsNotExist(e)
		if exists && e == nil {
			return absPath
		}
	}

	// If resource hasn't found then return relative path of resource
	return relPath
}

func (c *Configurator) writeSkeleton() error {
	configData, e := ioutil.ReadFile(c.ShareResourcePath(filepath.Join(skeletonDir, configName)))
	if e != nil {
		return e
	}

	return ioutil.WriteFile(c.FilePath, configData, 0644)
}

func (c *Configurator) GetConfig() (*InsteadmanConfig, error) {
	if c.FilePath == "" {
		c.FilePath = c.findConfigFileName()
	}

	// Write skeleton config if it isn't existing
	_, e := os.Stat(c.FilePath)
	exists := !os.IsNotExist(e)
	if !exists || e != nil {
		e = c.writeSkeleton()
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
		config.CalculatedGamesPath = c.gamesDir()
	}

	config.CalculatedInsteadManPath = config.InsteadManPath
	if config.CalculatedInsteadManPath == "" {
		config.CalculatedInsteadManPath = c.insteadManDir()
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
