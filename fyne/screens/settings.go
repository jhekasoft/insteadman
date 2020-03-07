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

func makeMainTab(config *configurator.InsteadmanConfig, configurator *configurator.Configurator, manager *manager.Manager) fyne.CanvasObject {
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

func makeAboutTab(mainIcon fyne.Resource, version string) fyne.CanvasObject {
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

// SettingsScreen is screen with settings
func SettingsScreen(
	config *configurator.InsteadmanConfig,
	configurator *configurator.Configurator,
	manager *manager.Manager,
	mainIcon fyne.Resource,
	version string) fyne.CanvasObject {
	return fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		widget.NewTabContainer(
			widget.NewTabItem("Main", makeMainTab(config, configurator, manager)),
			widget.NewTabItem("About", makeAboutTab(mainIcon, version)),
		),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), nil),
	)
}
