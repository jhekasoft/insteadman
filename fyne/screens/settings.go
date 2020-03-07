package screens

import (
	"fmt"
	"net/url"

	"github.com/jhekasoft/insteadman3/core/manager"

	"fyne.io/fyne/theme"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/jhekasoft/insteadman3/core/configurator"
)

// SettingsScreen is structure for Settings screen
type SettingsScreen struct {
	Manager      *manager.Manager
	Configurator *configurator.Configurator
	MainIcon     fyne.Resource
	Version      string
	Window       fyne.Window
	Screen       fyne.CanvasObject
}

// NewSettingsScreen is constructor for Settings screen
func NewSettingsScreen(
	manager *manager.Manager,
	configurator *configurator.Configurator,
	mainIcon fyne.Resource,
	version string,
	window fyne.Window) *SettingsScreen {
	scr := SettingsScreen{
		manager,
		configurator,
		mainIcon,
		version,
		window,
		nil,
	}

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		widget.NewTabContainer(
			widget.NewTabItem("Main", scr.makeMainTab()),
			widget.NewTabItem("About", scr.makeAboutTab()),
		),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
			scr.Window.Close()
		}),
	)

	return &scr
}

func (win *SettingsScreen) makeMainTab() fyne.CanvasObject {
	manager := win.Manager
	config := win.Manager.Config
	configurator := win.Configurator

	path := widget.NewEntry()
	path.SetPlaceHolder("INSTEAD path")
	path.SetText(config.InterpreterCommand)

	pathInfo := widget.NewLabel("")
	pathInfo.Hide()

	browseButton := widget.NewButton("Browse...", nil)
	browseButton.Disable()
	pathButtons := fyne.NewContainerWithLayout(
		layout.NewAdaptiveGridLayout(4),
		browseButton,
		widget.NewButton("Use built-in", nil),
		widget.NewButtonWithIcon("Detect", theme.SearchIcon(), func() {
			pathInfo.SetText("Detecting...")
			pathInfo.Show()
			command := manager.InterpreterFinder.Find()
			if command != nil {
				path.SetText(*command)
				pathInfo.SetText("INSTEAD has detected!")
			} else {
				pathInfo.SetText("INSTEAD hasn't detected!")
			}
		}),
		widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), func() {
			version, checkErr := manager.InterpreterFinder.Check(manager.InterpreterCommand())
			var txt string
			if checkErr != nil {
				if manager.IsBuiltinInterpreterCommand() {
					txt = "INSTEAD built-in check failed!"
				} else {
					txt = "INSTEAD check failed!"
				}
			} else {
				txt = fmt.Sprintf("INSTEAD %s has found!", version)
			}
			pathInfo.SetText(txt)
			pathInfo.Show()
		}),
	)

	language := widget.NewSelect([]string{"system", "en", "ru", "uk"}, nil)

	// Language
	if config.Lang != "" {
		language.SetSelected(config.Lang)
	}

	cleanCache := widget.NewButtonWithIcon("Clean", theme.DeleteIcon(), nil)

	configPathEntry := widget.NewEntry()
	configPathEntry.SetText(configurator.FilePath)
	configPathEntry.Disable()

	form := &widget.Form{}
	form.Append("INSTEAD path", widget.NewVBox(
		path,
		pathButtons,
		pathInfo,
	))
	form.Append("Language", language)
	form.Append("Cache", cleanCache)
	form.Append("Config path", configPathEntry)

	return form
}

func (win *SettingsScreen) makeAboutTab() fyne.CanvasObject {
	mainIcon := win.MainIcon
	version := win.Version

	siteURL := "https://jhekasoft.github.io/insteadman/"
	link, err := url.Parse(siteURL)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return fyne.NewContainerWithLayout(
		layout.NewCenterLayout(),
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(120, 120)), canvas.NewImageFromResource(mainIcon)),
			widget.NewVBox(
				widget.NewLabel("InsteadMan"),
				widget.NewLabel("Version: "+version),
				widget.NewHyperlinkWithStyle(siteURL, link, fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabel("© 2015-2020 InsteadMan"),
			),
		),
	)
}
