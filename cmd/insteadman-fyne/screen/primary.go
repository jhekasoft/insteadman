package screen

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/cmd/insteadman-fyne/screen/primary"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

type MainScreen struct {
	win    fyne.Window
	m      *manager.Manager
	c      *configurator.Configurator
	Screen fyne.CanvasObject
	games  []manager.Game

	// Widgets
	mainContainer  *fyne.Container
	searchEntry    *widget.Entry
	gamesContainer *widget.List
	gameInfo       *primary.GameInfoScreen
}

func (scr *MainScreen) RefreshList() {
	games, e := scr.m.GetSortedGamesByDateDesc()
	if e != nil {
		dialog.ShowError(e, scr.win)
		return
	}

	var keyword *string
	if scr.searchEntry.Text != "" {
		keyword = &scr.searchEntry.Text
	}
	games = manager.FilterGames(games, keyword, nil, nil, false)
	scr.games = games

	if scr.gamesContainer != nil {
		scr.gamesContainer.Refresh()
	}
}

// NewMainScreen is constructor for main screen
func NewMainScreen(
	win fyne.Window,
	m *manager.Manager,
	c *configurator.Configurator,
	onShowSettings func(),
	onShowAbout func()) *MainScreen {
	scr := MainScreen{
		m:   m,
		c:   c,
		win: win,
	}

	scr.gameInfo = primary.NewGameInfoScreen(win, m, c, scr.RefreshList)

	scr.searchEntry = widget.NewEntry()
	scr.searchEntry.SetPlaceHolder("Search")
	scr.searchEntry.OnChanged = func(s string) {
		scr.RefreshList()
	}
	searchEntryContainer := widget.NewHScrollContainer(scr.searchEntry)

	scr.gamesContainer = widget.NewList(
		func() int {
			return len(scr.games)
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
			// Check does game exist
			if scr.games == nil || len(scr.games) < index+1 {
				return
			}

			// Title
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(scr.games[index].Title)

			// Icon
			icon := item.(*fyne.Container).Objects[2].(*widget.Icon)
			icon.Hide()
			if scr.games[index].Installed {
				icon.SetResource(theme.ConfirmIcon())
				icon.Show()
			}
		},
	)

	scr.gamesContainer.OnSelected = func(index int) {
		scr.gameInfo.UpdateInfo(&scr.games[index])
	}

	// TODO: move to the goroutine
	scr.m.UpdateRepositories()
	scr.RefreshList()

	scr.mainContainer = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(searchEntryContainer, nil, nil, nil),
		searchEntryContainer,
		scr.gamesContainer,
	)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			loadingDialog := dialog.NewProgressInfinite("Refreshing", "Refreshing games...", win)
			loadingDialog.Show()

			scr.m.UpdateRepositories()
			scr.RefreshList()

			loadingDialog.Hide()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.InfoIcon(), onShowAbout),
		widget.NewToolbarAction(theme.SettingsIcon(), onShowSettings),
	)

	contentContainer := widget.NewHSplitContainer(
		scr.mainContainer,
		scr.gameInfo.Screen,
	)

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar,
		contentContainer,
	)

	return &scr
}
