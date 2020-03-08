package screens

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
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

	scroll := widget.NewScrollContainer(
		fyne.NewContainerWithLayout(
			layout.NewFixedGridLayout(fyne.NewSize(100, 130)),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
			gameItem(mainIcon),
		),
	)
	scroll.Resize(fyne.NewSize(400, 400))

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar,
		scroll,
	)

	return &scr
}

func gameItem(mainIcon fyne.Resource) fyne.CanvasObject {
	return widget.NewVBox(
		fyne.NewContainerWithLayout(
			layout.NewFixedGridLayout(fyne.NewSize(90, 90)),
			canvas.NewImageFromResource(mainIcon),
		),
		// widget.NewButton("Лифтёр 2", nil),
		widget.NewLabelWithStyle("Лифтёр 2", fyne.TextAlignCenter, fyne.TextStyle{}),
	)
}
