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
	Manager      *manager.Manager
	Configurator *configurator.Configurator
	MainIcon     fyne.Resource
	Window       fyne.Window
	Screen       fyne.CanvasObject
	tabs         *widget.TabContainer
}

// NewSettingsScreen is constructor for Settings screen
func NewSettingsScreen(
	m *manager.Manager,
	c *configurator.Configurator,
	mainIcon fyne.Resource,
	window fyne.Window) *SettingsScreen {
	scr := SettingsScreen{
		Manager:      m,
		Configurator: c,
		MainIcon:     mainIcon,
		Window:       window,
	}

	scr.tabs = widget.NewTabContainer(
		widget.NewTabItem("Main", scr.makeMainTab()),
		widget.NewTabItem("Repositories", scr.makeRepositoriesTab()),
		widget.NewTabItem("About", scr.makeAboutTab()),
	)

	okButton := widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
		// Don't use Close() because it will crash app
		scr.Window.Hide()
		scr.Window = nil
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
	return settings.NewCommonScreen(scr.Window, scr.Manager, scr.Configurator)
}

func (scr *SettingsScreen) makeRepositoriesTab() fyne.CanvasObject {
	return settings.NewRepositoriesScreen(scr.Manager, scr.Configurator)
}

func (scr *SettingsScreen) makeAboutTab() fyne.CanvasObject {
	return settings.NewAboutScreen()
}
