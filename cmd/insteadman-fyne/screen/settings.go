package screen

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

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
	tabs *widget.TabContainer
}

// NewSettingsScreen is constructor for Settings screen
func NewSettingsScreen(
	win fyne.Window,
	m *manager.Manager,
	c *configurator.Configurator) *SettingsScreen {
	scr := SettingsScreen{win: win, m: m, c: c}

	scr.tabs = widget.NewTabContainer(
		widget.NewTabItem("Main", scr.makeMainTab()),
		widget.NewTabItem("Repositories", scr.makeRepositoriesTab()),
		widget.NewTabItem("About", scr.makeAboutTab()),
	)

	okButton := widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
		// Don't use Close() because it will crash app
		scr.win.Hide()
		scr.win = nil
	})
	okButton.Style = widget.PrimaryButton

	scr.Screen = fyne.NewContainerWithLayout(
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
