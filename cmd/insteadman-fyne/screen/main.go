package screen

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	customWidget "github.com/jhekasoft/insteadman3/cmd/insteadman-fyne/widget"
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

	info := NewGameInfoScreen(m, c, mainIcon)

	search := widget.NewEntry()
	search.SetPlaceHolder("Search")

	games, e := scr.Manager.GetSortedGames()
	if e != nil {
		dialog.ShowError(e, scr.Window)
	}
	games = manager.FilterGames(games, nil, nil, nil, false)

	var items []fyne.CanvasObject
	for _, game := range games {
		currentGame := game // capture
		label := customWidget.NewGameLabel(&currentGame, func() {
			// scr.Manager.RunGame(&currentGame)
			info.UpdateInfo(&currentGame)
		})
		label.Wrapping = fyne.TextWrapWord
		label.Resize(fyne.NewSize(100, 20))
		items = append(items, label)
	}
	container := widget.NewVBox(items...)
	scroll := widget.NewVScrollContainer(
		container,
	)
	// scroll.Resize(fyne.NewSize(1, 400))

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.InfoIcon(), showAbout),
		widget.NewToolbarAction(theme.SettingsIcon(), showSettings),
	)

	infoContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		info.Screen,
		toolbar,
	)

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

// func gameItem(g *manager.Game, scr MainScreen) fyne.CanvasObject {
// 	var icon fyne.Resource = nil
// 	var b []byte = nil

// 	fileName, e := scr.Manager.GetGameImage(g)
// 	if e == nil {
// 		iconFile, e := os.Open(scr.Configurator.DataResourcePath(fileName))
// 		if e == nil {
// 			r := bufio.NewReader(iconFile)

// 			b, e = ioutil.ReadAll(r)
// 		}

// 		if e != nil {
// 			dialog.ShowError(e, scr.Window)
// 			icon = scr.MainIcon
// 		} else {
// 			icon = fyne.NewStaticResource("game_"+g.Name, b)
// 		}
// 	}

// 	return widget.NewVBox(
// 		fyne.NewContainerWithLayout(
// 			layout.NewFixedGridLayout(fyne.NewSize(140, 140)),
// 			canvas.NewImageFromResource(icon),
// 		),
// 		// widget.NewButton(title, nil),
// 		widget.NewLabel(g.Title),
// 		// widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{}),
// 	)
// }
