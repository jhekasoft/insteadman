package configurator

import (
    "encoding/json"
    "io/ioutil"
    "path/filepath"
)

type InsteadmanConfigType struct {
    Repositories []RepositoryType `json:"repositories"`
    GamesPath string `json:"games_path"`
    InterpreterCommand string `json:"interpreter_command"`
    Version string `json:"version"`
    UseBuiltinInterpreter bool `json:"use_builtin_interpreter"`
    Lang string `json:"lang"`
    CheckUpdateOnStart bool `json:"check_update_on_start"`
}

type RepositoryType struct {
    Name string `json:"name"`
    Url string `json:"url"`
}

func GetConfig() (*InsteadmanConfigType, error) {
    configFileName := filepath.Join(".", "instead-manager-settings.json")

    file, e := ioutil.ReadFile(configFileName)
    if e != nil {
        return nil, e
    }
    // fmt.Printf("%s\n", string(file))

    var config *InsteadmanConfigType
    json.Unmarshal(file, &config)

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