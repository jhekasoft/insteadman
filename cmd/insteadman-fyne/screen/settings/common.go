package settings

import (
	"fmt"
	"path"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

func NewCommonScreen(win fyne.Window, m *manager.Manager, c *configurator.Configurator) fyne.CanvasObject {

	language := widget.NewSelect([]string{"system", "en", "ru", "uk"}, nil)
	if m.Config.Lang != "" {
		language.SetSelected(m.Config.Lang)
	}

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("INSTEAD path")
	pathEntry.SetText(m.Config.InterpreterCommand)
	pathEntryContainer := widget.NewHScrollContainer(pathEntry)
	pathBrowseButton := widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
		// TODO: Move to function
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, win)
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
		}, win)
		fileURL := storage.NewFileURI(path.Dir(m.Config.InterpreterCommand))
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
		pathEntryContainer,
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
			command := m.InterpreterFinder.Find()
			if command != nil {
				pathEntry.SetText(*command)
				pathInfo.SetText("INSTEAD has detected!")
			} else {
				pathInfo.SetText("INSTEAD hasn't detected!")
			}
		}),
		widget.NewButtonWithIcon("Check", theme.ConfirmIcon(), func() {
			version, checkErr := m.InterpreterFinder.Check(m.InterpreterCommand())
			var txt string
			if checkErr != nil {
				if m.IsBuiltinInterpreterCommand() {
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

	insteadPath := widget.NewVBox(
		pathContainer,
		pathButtons,
		pathInfo,
	)

	gamesPathEntry := widget.NewEntry()
	gamesPathEntry.SetPlaceHolder(m.Config.CalculatedGamesPath)
	gamesPathEntry.SetText(m.Config.GamesPath)
	gamesPathEntryContainer := widget.NewHScrollContainer(gamesPathEntry)
	gamesPathBrowseButton := widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
		// TODO: Move to function
		fd := dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				return
			}

			gamesPathEntry.SetText(list.String())
		}, win)
		fileURL := storage.NewFileURI(m.Config.CalculatedGamesPath)
		dir, err := storage.ListerForURI(fileURL)
		if err == nil {
			fd.SetLocation(dir)
		} else {
			fyne.LogError("File dialog error", err)
		}
		fd.Show()
	})
	gamesPath := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, gamesPathBrowseButton),
		gamesPathEntryContainer,
		gamesPathBrowseButton,
	)

	clearCache := widget.NewButtonWithIcon("Clear", theme.DeleteIcon(), nil)

	configPathEntry := widget.NewEntry()
	configPathEntry.SetText(c.FilePath)
	configPathEntry.Disable()

	form := &widget.Form{}
	form.Append("Language", language)
	form.Append("INSTEAD path", insteadPath)
	form.Append("Games path", gamesPath)
	form.Append("Config path", configPathEntry)
	form.Append("Cache", clearCache)

	return form
}
