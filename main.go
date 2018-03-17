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

    games, e := man.GetSortedGames()
    keyword := "кот"
    //repo := "official"
    //lang := "en"
    games = manager.FilterGames(games, &keyword, nil, nil, false)

    for _, game := range games {
        fmt.Printf("%v, %v, %v, %v\n", game.Title, game.Name, game.RepositoryName, game.Languages)
    }

    fmt.Println("-------------------")

    interpreterPath := interpreterFinder.Find()
    if interpreterPath != nil {
        fmt.Printf("INSTEAD path: %v\n", *interpreterPath)
    }

    version, e := interpreterFinder.CheckInterpreter(config)
    fmt.Printf("INSTEAD error: %v\n", e)
    fmt.Printf("INSTEAD version: %v\n", version)

    //e = man.InstallGame(&manager.Game{Name:"lifter2",Url:"http://instead-games.ru//download/instead-lifter2-0.3.zip"})
    //fmt.Printf("Install error: %v\n", e)

    //e = man.RunGame(&manager.Game{Name:"lifter2"})
    //fmt.Printf("Run error: %v\n", e)

    //e = man.RemoveGame(&manager.Game{Name:"lifter2"})
    //fmt.Printf("Remove error: %v\n", e)
}

// Cli
// func main() {
//     argsWithProg := os.Args
//     argsWithoutProg := os.Args[1:]

//     fmt.Println(argsWithProg)
//     fmt.Println(argsWithoutProg)
//     fmt.Println(arg)
// }
