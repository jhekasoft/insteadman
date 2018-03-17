package main

import (
    "fmt"
    "os"
    "./configurator"
    "./manager"
    "./interpreter_finder"
)

func main() {
    config, e := configurator.GetConfig()
    if e != nil {
        fmt.Printf("Error: %v\n", e)
        os.Exit(1)
    }

    man := manager.Manager{Config: config}

    //man.DownloadRepositories()

    games, e := man.GetMergedGames()
    fmt.Printf("Games: %v\n", games)

    interpreterPath := interpreterFinder.Find()
    if interpreterPath != nil {
        fmt.Printf("INSTEAD path: %v\n", *interpreterPath)
    }

    version, e := interpreterFinder.CheckInterpreter(config)
    fmt.Printf("INSTEAD error: %v\n", e)
    fmt.Printf("INSTEAD version: %v\n", version)

    e = man.RunGame("lifter2")
    fmt.Printf("Run error: %v\n", e)
}

// Cli
// func main() {
//     argsWithProg := os.Args
//     argsWithoutProg := os.Args[1:]

//     fmt.Println(argsWithProg)
//     fmt.Println(argsWithoutProg)
//     fmt.Println(arg)
// }
