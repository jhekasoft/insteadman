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
	Manager        *manager.Manager
	Configurator   *configurator.Configurator
	MainIcon       fyne.Resource
	GamesContainer *widget.Box
	GameInfo       *GameInfoScreen
	Window         fyne.Window
	Screen         fyne.CanvasObject
}

func (scr *MainScreen) RefreshList() {
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
			scr.GameInfo.UpdateInfo(&currentGame)
		})
		label.Resize(fyne.NewSize(100, 20))
		items = append(items, label)
	}
	scr.GamesContainer.Children = items
	scr.GamesContainer.Refresh()
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
		Manager:        m,
		Configurator:   c,
		MainIcon:       mainIcon,
		GamesContainer: widget.NewVBox(),
		GameInfo:       NewGameInfoScreen(m, c, mainIcon, window),
		Window:         window,
	}

	scr.GameInfo.UpdateF = scr.RefreshList

	search := widget.NewEntry()
	search.SetPlaceHolder("Search")

	// TODO: move to the goroutine
	scr.Manager.UpdateRepositories()
	scr.RefreshList()

	scroll := widget.NewVScrollContainer(
		scr.GamesContainer,
	)
	// scroll.Resize(fyne.NewSize(1, 400))

	mainContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(search, nil, nil, nil),
		search,
		scroll,
	)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			scr.Manager.UpdateRepositories()
			scr.RefreshList()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.InfoIcon(), showAbout),
		widget.NewToolbarAction(theme.SettingsIcon(), showSettings),
	)

	contentContainer := widget.NewHSplitContainer(
		mainContainer,
		scr.GameInfo.Screen,
	)

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar,
		contentContainer,
	)

	return &scr
}
