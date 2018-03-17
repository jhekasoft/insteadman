package configurator

import (
    //"encoding/json"
    "github.com/ghodss/yaml"
    "io/ioutil"
    "path/filepath"
    "os"
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
}

type Repository struct {
    Name string `json:"name"`
    Url  string `json:"url"`
}

const configName = "config.yml"

func insteadDir() string {
    return filepath.Join(userHomeDir(), ".instead")
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

func configFileName() string {
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

func GetConfig() (*InsteadmanConfig, error) {
    file, e := ioutil.ReadFile(configFileName())
    if e != nil {
        return nil, e
    }
    // fmt.Printf("%s\n", string(file))

    var config *InsteadmanConfig
    yaml.Unmarshal(file, &config)

    if config.GamesPath == "" {
        config.GamesPath = gamesDir()
    }

    if config.InsteadManPath == "" {
        config.InsteadManPath = insteadManDir()
    }

    return config, nil

    // write config
    // config.Lang = "ru"

    // bytes, e := json.MarshalIndent(config, "", "  ")
    // if e != nil {
    //     fmt.Printf("Config error: %v\n", e)
    //     os.Exit(1)
    // }

    // ioutil.WriteFile(configFileName, bytes, 0644)
}
