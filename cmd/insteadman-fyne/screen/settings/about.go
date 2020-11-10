package settings

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/cmd/insteadman-fyne/data"
	"github.com/jhekasoft/insteadman3/core/manager"
)

func NewAboutScreen() fyne.CanvasObject {
	mainIcon := data.InsteadManLogo

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
