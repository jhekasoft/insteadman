package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/jhekasoft/insteadman3/cmd/insteadman-fyne/screen/settings"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

// SettingsScreen is structure for Settings screen
type SettingsScreen struct {
	win    fyne.Window
	m      *manager.Manager
	c      *configurator.Configurator
	Screen fyne.CanvasObject

	// Widgets
	tabs *container.AppTabs
}

// NewSettingsScreen is constructor for Settings screen
func NewSettingsScreen(
	win fyne.Window,
	m *manager.Manager,
	c *configurator.Configurator) *SettingsScreen {
	scr := SettingsScreen{win: win, m: m, c: c}

	scr.tabs = container.NewAppTabs(
		container.NewTabItem("Main", scr.makeMainTab()),
		container.NewTabItem("Repositories", scr.makeRepositoriesTab()),
		container.NewTabItem("About", scr.makeAboutTab()),
	)

	okButton := widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
		// Don't use Close() because it will crash app
		scr.win.Hide()
		scr.win = nil
	})
	okButton.Importance = widget.HighImportance

	scr.Screen = container.New(
		layout.NewBorderLayout(nil, okButton, nil, nil),
		scr.tabs,
		okButton,
	)

	return &scr
}

func (scr *SettingsScreen) SetMainTab() {
	scr.tabs.SelectTabIndex(0)
}

func (scr *SettingsScreen) SetRepositoriesTab() {
	scr.tabs.SelectTabIndex(1)
}

func (scr *SettingsScreen) SetAboutTab() {
	scr.tabs.SelectTabIndex(2)
}

func (scr *SettingsScreen) makeMainTab() fyne.CanvasObject {
	return settings.NewCommonScreen(scr.win, scr.m, scr.c)
}

func (scr *SettingsScreen) makeRepositoriesTab() fyne.CanvasObject {
	return settings.NewRepositoriesScreen(scr.m, scr.c)
}

func (scr *SettingsScreen) makeAboutTab() fyne.CanvasObject {
	return settings.NewAboutScreen()
}
