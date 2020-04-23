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
	Window       fyne.Window
	Screen       fyne.CanvasObject
}

// NewMainScreen is constructor for main screen
func NewMainScreen(
	m *manager.Manager,
	c *configurator.Configurator,
	mainIcon fyne.Resource,
	window fyne.Window,
	showSettings func(),
	showAbout func()) *MainScreen {
	scr := MainScreen{
		Manager:      m,
		Configurator: c,
		MainIcon:     mainIcon,
		Window:       window,
	}

	search := widget.NewEntry()
	search.SetPlaceHolder("Search")
	// buttons := fyne.NewContainerWithLayout(
	// 	layout.NewGridLayoutWithColumns(1),
	// 	// widget.NewCheck("Installed", nil),
	// 	// widget.NewButtonWithIcon("Update", theme.ViewRefreshIcon(), nil),
	// 	widget.NewButtonWithIcon("", theme.SettingsIcon(), showSettings),
	// 	// widget.NewButtonWithIcon("About", theme.InfoIcon(), showAbout),
	// )
	// toolbar := fyne.NewContainerWithLayout(
	// 	layout.NewBorderLayout(nil, nil, nil, buttons),
	// 	buttons,
	// )

	games, e := scr.Manager.GetSortedGames()
	if e != nil {
		dialog.ShowError(e, scr.Window)
	}
	games = manager.FilterGames(games, nil, nil, nil, false)

	var items []fyne.CanvasObject
	for _, game := range games {
		label := widget.NewLabel(game.Title)
		label.Wrapping = fyne.TextWrapWord
		label.Resize(fyne.NewSize(100, 20))
		items = append(items, label)
		// var installedIcon fyne.Resource
		// if game.Installed {
		// 	installedIcon = theme.CheckButtonCheckedIcon()
		// }

		// button := widget.NewButtonWithIcon(game.Title, installedIcon, nil)
		// items = append(items, button)
		// currentGame := game // capture
		// items = append(items, widget.NewButton(game.Title, func() {
		// 	scr.Manager.RunGame(&currentGame)
		// }))
	}
	container := widget.NewVBox(items...)
	// container := fyne.NewContainerWithLayout(
	// 	layout.NewFixedGridLayout(fyne.NewSize(150, 200)),
	// )
	// var items []fyne.CanvasObject = nil
	// for _, game := range games {
	// 	container.Append(widget.NewButtonWithIcon(game.Title, theme.CheckButtonCheckedIcon(), nil))
	// 	// container.AddObject(gameItem(&game, scr))
	// 	// items = append(items, gameItem(game.Title, MainIcon))
	// }
	scroll := widget.NewVScrollContainer(
		container,
	)
	// scroll.Resize(fyne.NewSize(1, 400))

	desc := widget.NewLabel("Лифтёр 3 описание большое очень")
	// desc.Resize(fyne.NewSize(100, 100))
	desc.Wrapping = fyne.TextWrapWord
	info := widget.NewVBox(
		widget.NewLabelWithStyle("Лифтёр 3", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		desc,
	)

	// buttons := fyne.NewContainerWithLayout(
	// 	layout.NewGridLayoutWithColumns(2),
	// 	widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), showSettings),
	// 	widget.NewButtonWithIcon("About", theme.InfoIcon(), showAbout),
	// )
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.InfoIcon(), showAbout),
		widget.NewToolbarAction(theme.SettingsIcon(), showSettings),
	)
	infoContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		info,
		toolbar,
	)

	// buttons := fyne.NewContainerWithLayout(
	// 	layout.NewGridLayoutWithColumns(2),
	// 	// widget.NewCheck("Installed", nil),
	// 	widget.NewButtonWithIcon("Install", theme.ContentAddIcon(), nil),
	// 	widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), showSettings),
	// 	// widget.NewButtonWithIcon("About", theme.InfoIcon(), showAbout),
	// )

	installButton := widget.NewButtonWithIcon("Install", theme.ContentAddIcon(), nil)
	mainContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(search, installButton, nil, nil),
		search,
		scroll,
		installButton,
	)

	scr.Screen = widget.NewHSplitContainer(
		mainContainer,
		infoContainer,
	)
	// fyne.NewContainerWithLayout(
	// 	layout.NewBorderLayout(nil, nil, nil, infoContainer),
	// 	// toolbar,
	// 	mainContainer,
	// 	infoContainer,
	// )

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
