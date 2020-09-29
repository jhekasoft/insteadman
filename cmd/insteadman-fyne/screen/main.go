package screen

import (
	"fmt"

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
	SearchEntry    *widget.Entry
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

	var keyword *string
	if scr.SearchEntry.Text != "" {
		keyword = &scr.SearchEntry.Text
	}
	games = manager.FilterGames(games, keyword, nil, nil, false)
	fmt.Println(games)

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
	scr.GamesContainer.Show()
	scr.GamesContainer.Refresh()

	if scr.MainContainer != nil {
		scr.MainContainer.Refresh()
	}

	// scr.Window.Canvas().Refresh(scr.Window.Content())
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
		GameInfo:     NewGameInfoScreen(m, c, mainIcon, window),
		Window:       window,
	}

	scr.GameInfo.UpdateF = scr.RefreshList

	scr.SearchEntry = widget.NewEntry()
	scr.SearchEntry.SetPlaceHolder("Search")
	scr.SearchEntry.OnChanged = func(s string) {
		scr.RefreshList()
	}

	// TODO: move to the goroutine
	scr.Manager.UpdateRepositories()
	scr.RefreshList()

	scr.MainContainer = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(scr.SearchEntry, nil, nil, nil),
		scr.SearchEntry,
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
