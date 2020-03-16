package screens

import (
	"bufio"
	"io/ioutil"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

type MainScreen struct {
	Manager      *manager.Manager
	Configurator *configurator.Configurator
	MainIcon     fyne.Resource
	Version      string
	Window       fyne.Window
	Screen       fyne.CanvasObject
}

// NewMainScreen is constructor for main screen
func NewMainScreen(
	manager *manager.Manager,
	configurator *configurator.Configurator,
	mainIcon fyne.Resource,
	version string,
	window fyne.Window,
	showSettings func(),
	showAbout func()) *MainScreen {
	scr := MainScreen{
		Manager:      manager,
		Configurator: configurator,
		MainIcon:     mainIcon,
		Version:      version,
		Window:       window,
	}

	search := widget.NewEntry()
	search.SetPlaceHolder("Search")
	toolbar := fyne.NewContainerWithLayout(
		layout.NewGridLayoutWithColumns(5),
		widget.NewButtonWithIcon("Update", theme.ViewRefreshIcon(), nil),
		search,
		widget.NewCheck("Installed", nil),
		widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), showSettings),
		widget.NewButtonWithIcon("About", theme.InfoIcon(), showAbout),
	)

	// games, e := scr.Manager.GetSortedGamesByDateDesc()
	// if e != nil {
	// 	dialog.ShowError(e, scr.Window)
	// }

	container := fyne.NewContainerWithLayout(
		layout.NewFixedGridLayout(fyne.NewSize(150, 200)),
	)
	// var items []fyne.CanvasObject = nil
	// for _, game := range games {
	// 	container.AddObject(gameItem(&game, scr))
	// 	// items = append(items, gameItem(game.Title, mainIcon))
	// }
	scroll := widget.NewScrollContainer(
		container,
	)
	scroll.Resize(fyne.NewSize(400, 400))

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar,
		scroll,
	)

	return &scr
}

func gameItem(g *manager.Game, scr MainScreen) fyne.CanvasObject {
	var icon fyne.Resource = nil
	var b []byte = nil

	fileName, e := scr.Manager.GetGameImage(g)
	if e == nil {
		iconFile, e := os.Open(scr.Configurator.DataResourcePath(fileName))
		if e == nil {
			r := bufio.NewReader(iconFile)

			b, e = ioutil.ReadAll(r)
		}

		if e != nil {
			dialog.ShowError(e, scr.Window)
			icon = scr.MainIcon
		} else {
			icon = fyne.NewStaticResource("game_"+g.Name, b)
		}
	}

	return widget.NewVBox(
		fyne.NewContainerWithLayout(
			layout.NewFixedGridLayout(fyne.NewSize(140, 140)),
			canvas.NewImageFromResource(icon),
		),
		// widget.NewButton(title, nil),
		widget.NewLabel(g.Title),
		// widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{}),
	)
}
