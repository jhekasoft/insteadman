package screen

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

type GameInfoScreen struct {
	Manager      *manager.Manager
	Configurator *configurator.Configurator
	MainIcon     fyne.Resource
	Title        *widget.Label
	Desc         *widget.Label
	Screen       fyne.CanvasObject
	Image        *widget.Icon
}

func (scr *GameInfoScreen) UpdateInfo(g *manager.Game) {
	scr.Title.SetText(g.Title)
	scr.Desc.SetText(g.Description)

	var icon fyne.Resource = nil
	var b []byte = nil

	fileName, e := scr.Manager.GetGameImage(g)
	if e == nil {
		iconFile, e := os.Open(scr.Configurator.DataResourcePath(fileName))
		if e == nil {
			r := bufio.NewReader(iconFile)

			b, e = ioutil.ReadAll(r)
		}

		if e != nil {
			// dialog.ShowError(e, scr.Window)
			fmt.Printf("Error: %v\n", e)
			icon = scr.MainIcon
		} else {
			icon = fyne.NewStaticResource("game_"+g.Name, b)
		}
	}

	scr.Image.SetResource(icon)
}

func NewGameInfoScreen(
	m *manager.Manager,
	c *configurator.Configurator,
	mainIcon fyne.Resource) *GameInfoScreen {
	scr := GameInfoScreen{
		Manager:      m,
		Configurator: c,
		MainIcon:     mainIcon,
	}

	scr.Image = widget.NewIcon(mainIcon)
	imageContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil),
		scr.Image,
	)
	scr.Title = widget.NewLabelWithStyle("InsteadMan", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	scr.Desc = widget.NewLabel("Выберите игру слева в списке")
	scr.Desc.Wrapping = fyne.TextWrapWord

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(scr.Title, scr.Desc, nil, nil),
		scr.Title,
		imageContainer,
		scr.Desc,
	)

	return &scr
}
