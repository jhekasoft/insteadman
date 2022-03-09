package settings

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
		container.NewHBox(
			container.New(layout.NewGridWrapLayout(fyne.NewSize(160, 160)), canvas.NewImageFromResource(mainIcon)),
			container.NewVBox(
				widget.NewLabelWithStyle("InsteadMan", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("Version: "+manager.Version),
				widget.NewHyperlink(siteURL, link),
				widget.NewLabel("License: MIT"),
				widget.NewLabel("© 2015-2020 InsteadMan"),
			),
		),
	)
}
