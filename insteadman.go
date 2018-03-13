package main

import (
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "net/http"
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


func main() {
    config := getConfig()
    fmt.Printf("Config: %v\n", config)

    downloadRepositories(config)
}

func downloadRepository(fileName, url string) error {
    // Create the file
    out, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer out.Close()

    // Download the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the data to the file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}

func downloadRepositories(config InsteadmanConfigType) {
    repositoriesDir := filepath.Join(".", "repositories")
    os.MkdirAll(repositoriesDir, os.ModePerm)

    for _, repo := range config.Repositories {
        fmt.Printf("%v %v\n", repo.Name, repo.Url)
        downloadRepository(filepath.Join(repositoriesDir, repo.Name + ".xml"), repo.Url)
    }
}

func getConfig() InsteadmanConfigType {
    configFileName := filepath.Join(".", "instead-manager-settings.json")

    file, e := ioutil.ReadFile(configFileName)
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }
    // fmt.Printf("%s\n", string(file))

    var config InsteadmanConfigType
    json.Unmarshal(file, &config)
    // fmt.Printf("Results: %v\n", config)

    return config

    // write config
    // config.Lang = "ru"

    // bytes, e := json.MarshalIndent(config, "", "  ")
    // if e != nil {
    //     fmt.Printf("Config error: %v\n", e)
    //     os.Exit(1)
    // }

    // ioutil.WriteFile(configFileName, bytes, 0644)
}


// Cli
// func main() {
//     argsWithProg := os.Args
//     argsWithoutProg := os.Args[1:]

//     fmt.Println(argsWithProg)
//     fmt.Println(argsWithoutProg)
//     fmt.Println(arg)
// }
