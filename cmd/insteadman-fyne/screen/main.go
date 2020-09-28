package screen

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

type MainScreen struct {
	Manager        *manager.Manager
	Configurator   *configurator.Configurator
	MainIcon       fyne.Resource
	MainContainer  *fyne.Container
	GamesContainer *widget.List
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

	scr.GamesContainer = widget.NewList(
		func() int {
			return len(games)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(
				layout.NewHBoxLayout(),
				widget.NewLabel("Game"),
				layout.NewSpacer(),
				widget.NewIcon(nil),
			)
		},
		func(index int, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(games[index].Title)

			iconRes := theme.ConfirmIcon()
			if !games[index].Installed {
				iconRes = nil
			}
			item.(*fyne.Container).Objects[2].(*widget.Icon).SetResource(iconRes)
		},
	)

	scr.GamesContainer.OnItemSelected = func(index int) {
		scr.GameInfo.UpdateInfo(&games[index])
	}

	if scr.MainContainer != nil {
		scr.MainContainer.Refresh()
	}
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
		GamesContainer: nil,
		GameInfo:       NewGameInfoScreen(m, c, mainIcon, window),
		Window:         window,
	}

	scr.GameInfo.UpdateF = scr.RefreshList

	search := widget.NewEntry()
	search.SetPlaceHolder("Search")

	// TODO: move to the goroutine
	scr.Manager.UpdateRepositories()
	scr.RefreshList()

	scr.MainContainer = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(search, nil, nil, nil),
		search,
		scr.GamesContainer,
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
		scr.MainContainer,
		scr.GameInfo.Screen,
	)

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar,
		contentContainer,
	)

	return &scr
}
