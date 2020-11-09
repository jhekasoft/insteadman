package screen

import (
	"fmt"
	"net/url"
	"path"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
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
	language := widget.NewSelect([]string{"system", "en", "ru", "uk"}, nil)
	if scr.Manager.Config.Lang != "" {
		language.SetSelected(scr.Manager.Config.Lang)
	}

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("INSTEAD path")
	pathEntry.SetText(scr.Manager.Config.InterpreterCommand)
	pathBrowseButton := widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
		// TODO: Move to function
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, scr.Window)
				return
			}
			if reader == nil {
				return
			}

			pathEntry.SetText(reader.URI().String())

			err = reader.Close()
			if err != nil {
				fyne.LogError("Failed to close stream", err)
			}
		}, scr.Window)
		fileURL := storage.NewFileURI(path.Dir(scr.Manager.Config.InterpreterCommand))
		dir, err := storage.ListerForURI(fileURL)
		if err == nil {
			fd.SetLocation(dir)
		} else {
			fyne.LogError("File dialog error", err)
		}
		fd.Show()
	})
	pathContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, pathBrowseButton),
		pathEntry,
		pathBrowseButton,
	)

	pathInfo := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	pathInfo.Hide()

	// browseButton := widget.NewButton("Browse...", nil)
	// browseButton.Disable()
	pathButtons := fyne.NewContainerWithLayout(
		layout.NewAdaptiveGridLayout(3),
		// browseButton,
		widget.NewButton("Use built-in", nil),
		widget.NewButtonWithIcon("Detect", theme.SearchIcon(), func() {
			pathInfo.SetText("Detecting...")
			pathInfo.Show()
			command := scr.Manager.InterpreterFinder.Find()
			if command != nil {
				pathEntry.SetText(*command)
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

	path := widget.NewVBox(
		pathContainer,
		pathButtons,
		pathInfo,
	)

	gamesPathEntry := widget.NewEntry()
	gamesPathEntry.SetPlaceHolder(scr.Manager.Config.CalculatedGamesPath)
	gamesPathEntry.SetText(scr.Manager.Config.GamesPath)
	gamesPathBrowseButton := widget.NewButtonWithIcon("", theme.FolderIcon(), nil)
	gamesPath := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, gamesPathBrowseButton),
		gamesPathEntry,
		gamesPathBrowseButton,
	)

	clearCache := widget.NewButtonWithIcon("Clear", theme.DeleteIcon(), nil)

	configPathEntry := widget.NewEntry()
	configPathEntry.SetText(scr.Configurator.FilePath)
	configPathEntry.Disable()

	form := &widget.Form{}
	form.Append("Language", language)
	form.Append("INSTEAD path", path)
	form.Append("Games path", gamesPath)
	form.Append("Config path", configPathEntry)
	form.Append("Cache", clearCache)

	return form
}

func (scr *SettingsScreen) makeRepositoriesTab() fyne.CanvasObject {
	repositories := scr.Manager.GetRepositories()

	listHeader := widget.NewVBox(
		fyne.NewContainerWithLayout(
			layout.NewGridLayoutWithColumns(2),
			widget.NewLabelWithStyle("Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("URL", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		),
		widget.NewSeparator(),
	)

	list := widget.NewList(
		func() int {
			return len(repositories)
		},
		func() fyne.CanvasObject {
			return fyne.NewContainerWithLayout(
				layout.NewGridLayoutWithColumns(2),
				widget.NewLabel(""),
				widget.NewLabel(""),
			)
		},
		func(index int, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(repositories[index].Name)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(repositories[index].Url)
		},
	)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), nil),
		widget.NewToolbarAction(theme.ContentAddIcon(), nil),
		widget.NewToolbarAction(theme.DeleteIcon(), nil),
		// widget.NewToolbarAction(theme.MoveUpIcon(), nil),
		// widget.NewToolbarAction(theme.MoveDownIcon(), nil),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ContentUndoIcon(), nil),
	)

	return fyne.NewContainerWithLayout(
		layout.NewBorderLayout(listHeader, toolbar, nil, nil),
		listHeader,
		list,
		toolbar,
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
