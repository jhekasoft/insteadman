package main

import (
    "fmt"
    "os"
    "./configurator"
    "./manager"
)

func main() {
    config, e := configurator.GetConfig()
    if e != nil {
        fmt.Printf("Error: %v\n", e)
        os.Exit(1)
    }

    fmt.Printf("Config: %v\n", *config)

    manager.DownloadRepositories(config)

    games, e := manager.ParseRepositories()
    fmt.Printf("Config: %v\n", games)
}

// Cli
// func main() {
//     argsWithProg := os.Args
//     argsWithoutProg := os.Args[1:]

//     fmt.Println(argsWithProg)
//     fmt.Println(argsWithoutProg)
//     fmt.Println(arg)
// }
