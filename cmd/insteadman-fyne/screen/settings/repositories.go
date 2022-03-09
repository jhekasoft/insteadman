package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

func NewRepositoriesScreen(m *manager.Manager, c *configurator.Configurator) fyne.CanvasObject {
	repositories := m.GetRepositories()

	listHeader := container.NewVBox(
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
