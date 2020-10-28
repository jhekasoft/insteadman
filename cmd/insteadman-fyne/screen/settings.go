package screen

import (
	"fmt"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

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

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		scr.tabs,
		layout.NewSpacer(),
		widget.NewButtonWithIcon("OK", theme.ConfirmIcon(), func() {
			scr.Window.Close()
		}),
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
	path := widget.NewEntry()
	path.SetPlaceHolder("INSTEAD path")
	path.SetText(scr.Manager.Config.InterpreterCommand)

	pathInfo := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
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
			command := scr.Manager.InterpreterFinder.Find()
			if command != nil {
				path.SetText(*command)
				pathInfo.SetText("INSTEAD has detected!")
			} else {
				pathInfo.SetText("INSTEAD hasn't detected!")
			}
		}),
		widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), func() {
			version, checkErr := scr.Manager.InterpreterFinder.Check(scr.Manager.InterpreterCommand())
			var txt string
			if checkErr != nil {
				if scr.Manager.IsBuiltinInterpreterCommand() {
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
	if scr.Manager.Config.Lang != "" {
		language.SetSelected(scr.Manager.Config.Lang)
	}

	cleanCache := widget.NewButtonWithIcon("Clean", theme.DeleteIcon(), nil)

	configPathEntry := widget.NewEntry()
	configPathEntry.SetText(scr.Configurator.FilePath)
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

func (scr *SettingsScreen) makeRepositoriesTab() fyne.CanvasObject {
	repositories := scr.Manager.GetRepositories()

	return widget.NewList(
		func() int {
			return len(repositories)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(
				layout.NewGridLayoutWithColumns(2),
				widget.NewEntry(),
				widget.NewEntry(),
			)
		},
		func(index int, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Entry).SetText(repositories[index].Name)
			item.(*fyne.Container).Objects[1].(*widget.Entry).SetText(repositories[index].Url)
		},
	)
}

func (scr *SettingsScreen) makeAboutTab() fyne.CanvasObject {
	mainIcon := scr.MainIcon

	siteURL := "https://jhekasoft.github.io/insteadman/"
	link, err := url.Parse(siteURL)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return fyne.NewContainerWithLayout(
		layout.NewCenterLayout(),
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(160, 160)), canvas.NewImageFromResource(mainIcon)),
			widget.NewVBox(
				widget.NewLabelWithStyle("InsteadMan", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("Version: "+manager.Version),
				widget.NewHyperlink(siteURL, link),
				widget.NewLabel("License: MIT"),
				widget.NewLabel("© 2015-2020 InsteadMan"),
			),
		),
	)
}
