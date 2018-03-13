package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
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
    file, e := ioutil.ReadFile("./instead-manager-settings.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }
    // fmt.Printf("%s\n", string(file))

    var config InsteadmanConfigType
    json.Unmarshal(file, &config)
    fmt.Printf("Results: %v\n", config)
}


// func cli() {
//     argsWithProg := os.Args
//     argsWithoutProg := os.Args[1:]

//     fmt.Println(argsWithProg)
//     fmt.Println(argsWithoutProg)
//     fmt.Println(arg)
// }
