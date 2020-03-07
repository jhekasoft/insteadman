package screens

import (
	"fyne.io/fyne"
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

	scr.Screen = widget.NewVBox(
		widget.NewEntry(),
		widget.NewLabel(manager.Config.CalculatedGamesPath),
		widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), showSettings),
		widget.NewButtonWithIcon("About", theme.InfoIcon(), showAbout),
	)

	return &scr
}
