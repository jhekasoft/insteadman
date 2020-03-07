package screens

import (
	"net/url"

	"fyne.io/fyne/theme"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/jhekasoft/insteadman3/core/configurator"
)

func makeMainTab(config *configurator.InsteadmanConfig, configurator *configurator.Configurator) fyne.Widget {
	instead := widget.NewEntry()
	instead.SetPlaceHolder("INSTEAD path")
	instead.SetText(config.InterpreterCommand)

	insteadButtons := fyne.NewContainerWithLayout(
		layout.NewAdaptiveGridLayout(4),
		widget.NewButton("Browse...", nil),
		widget.NewButton("Use built-in", nil),
		widget.NewButtonWithIcon("Detect", theme.SearchIcon(), nil),
		widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), nil),
	)

	language := widget.NewSelect([]string{"system", "English", "Russian", "Ukrainian"}, nil)

	cleanCache := widget.NewButtonWithIcon("Clean", theme.DeleteIcon(), nil)

	form := &widget.Form{}
	form.Append("INSTEAD path", instead)
	form.Append("", insteadButtons)
	form.Append("Language", language)
	form.Append("Cache", cleanCache)
	form.Append("Config path", widget.NewLabel(configurator.FilePath))

	return form
}

func makeAboutTab(mainIcon fyne.Resource) fyne.Widget {
	siteURL := "https://jhekasoft.github.io/insteadman/"
	link, err := url.Parse(siteURL)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return widget.NewHBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(120, 120)), canvas.NewImageFromResource(mainIcon)),
		widget.NewVBox(
			widget.NewLabel("InsteadMan"),
			widget.NewHyperlinkWithStyle(siteURL, link, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabel("© 2015-2020 InsteadMan"),
		),
	)
}

// SettingsScreen is screen with settings
func SettingsScreen(config *configurator.InsteadmanConfig, configurator *configurator.Configurator, mainIcon fyne.Resource) fyne.CanvasObject {
	return fyne.NewContainerWithLayout(
		layout.NewMaxLayout(),
		widget.NewTabContainer(
			widget.NewTabItem("Main", makeMainTab(config, configurator)),
			widget.NewTabItem("About", makeAboutTab(mainIcon)),
		),
	)
}
