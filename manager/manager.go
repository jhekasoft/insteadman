package manager

import (
    "io"
    "os"
    "path/filepath"
    "net/http"
    "encoding/xml"
    "io/ioutil"
    // "fmt"
    "../configurator"
)

type RepositoryGameList struct {
    // XMLName xml.Name `xml:"game_list"`
    GameList []RepositoryGame `xml:"game"`
}

type RepositoryGame struct {
    // XMLName xml.Name `xml:"game"`
    Name string `xml:"name"`
    Title string `xml:"title"`
    Version string `xml:"version"`
    Url string `xml:"url"`
    Size int `xml:"size"`
    Lang string `xml:"lang"`
    Descurl string `xml:"descurl"`
    Author string `xml:"author"`
    Description string `xml:"description"`
    Image string `xml:"image"`
    Langs []string `xml:"langs>lang"`
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

func DownloadRepositories(config *configurator.InsteadmanConfigType) {
    repositoriesDir := filepath.Join(".", "repositories")
    os.MkdirAll(repositoriesDir, os.ModePerm)

    for _, repo := range config.Repositories {
        // fmt.Printf("%v %v\n", repo.Name, repo.Url)
        downloadRepository(filepath.Join(repositoriesDir, repo.Name + ".xml"), repo.Url)
    }
}

func ParseRepositories() ([]RepositoryGame, error) {
    repositoriesDir := filepath.Join(".", "repositories")
    files, e := filepath.Glob(filepath.Join(repositoriesDir, "*.xml"))
    if e != nil {
        return nil, e
    }

    games := []RepositoryGame{}
    for _, fileName := range files {
        // fmt.Printf("File: %v\n", fileName)

        gameList, e := parseRepository(filepath.Join(".", fileName))
        if e == nil {
            games = append(games, gameList.GameList...)
            // fmt.Printf("Games: %v\n", *gameList)
        }
    }

    return games, nil
}

func parseRepository(fileName string) (*RepositoryGameList, error) {
    file, e := ioutil.ReadFile(fileName)
    if e != nil {
        return nil, e
    }
    // fmt.Printf("File: %v\n", string(file))

    var gameList *RepositoryGameList
    e = xml.Unmarshal(file, &gameList)
    // fmt.Printf("Games: %v\n", *gameList)
    if e != nil {
        return nil, e
    }

    return gameList, nil
}
