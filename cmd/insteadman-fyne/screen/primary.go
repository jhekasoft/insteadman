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
	Manager        *manager.Manager
	Configurator   *configurator.Configurator
	MainIcon       fyne.Resource
	MainContainer  *fyne.Container
	SearchEntry    *widget.Entry
	GamesContainer *widget.List
	GameInfo       *primary.GameInfoScreen
	Window         fyne.Window
	Screen         fyne.CanvasObject
	Games          []manager.Game
}

func (scr *MainScreen) RefreshList() {
	games, e := scr.Manager.GetSortedGamesByDateDesc()
	if e != nil {
		dialog.ShowError(e, scr.Window)
		return
	}

	var keyword *string
	if scr.SearchEntry.Text != "" {
		keyword = &scr.SearchEntry.Text
	}
	games = manager.FilterGames(games, keyword, nil, nil, false)
	scr.Games = games

	if scr.GamesContainer != nil {
		scr.GamesContainer.Refresh()
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
		Manager:      m,
		Configurator: c,
		MainIcon:     mainIcon,
		GameInfo:     primary.NewGameInfoScreen(m, c, mainIcon, window),
		Window:       window,
	}

	scr.GameInfo.UpdateF = scr.RefreshList

	scr.SearchEntry = widget.NewEntry()
	scr.SearchEntry.SetPlaceHolder("Search")
	scr.SearchEntry.OnChanged = func(s string) {
		scr.RefreshList()
	}

	scr.GamesContainer = widget.NewList(
		func() int {
			return len(scr.Games)
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
			// Title
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(scr.Games[index].Title)

			// Icon
			icon := item.(*fyne.Container).Objects[2].(*widget.Icon)
			icon.Hide()
			if scr.Games[index].Installed {
				icon.SetResource(theme.ConfirmIcon())
				icon.Show()
			}
		},
	)

	scr.GamesContainer.OnSelected = func(index int) {
		scr.GameInfo.UpdateInfo(&scr.Games[index])
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
			loadingDialog := dialog.NewProgressInfinite("Refreshing", "Refreshing games...", window)
			loadingDialog.Show()

			scr.Manager.UpdateRepositories()
			scr.RefreshList()

			loadingDialog.Hide()
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
