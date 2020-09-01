package screen

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

type GameInfoScreen struct {
	Manager       *manager.Manager
	Configurator  *configurator.Configurator
	MainIcon      fyne.Resource
	Title         *widget.Label
	Desc          *widget.Label
	Screen        fyne.CanvasObject
	Image         *widget.Icon
	Hyperlink     *widget.Hyperlink
	InstallButton *widget.Button
	RunButton     *widget.Button
	DeleteButton  *widget.Button
	Game          *manager.Game
}

func (scr *GameInfoScreen) UpdateInfo(g *manager.Game) {
	scr.Game = g

	scr.Title.SetText(g.Title)
	scr.Desc.SetText(g.Description)

	if g.Descurl != "" {
		scr.Hyperlink.SetURLFromString(g.Descurl)
		scr.Hyperlink.Show()
	}

	if g.Installed {
		scr.InstallButton.Hide()
		scr.RunButton.Show()
		scr.DeleteButton.Show()
	} else {
		scr.InstallButton.Show()
		scr.RunButton.Hide()
		scr.DeleteButton.Hide()
	}

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
	scr.Title = widget.NewLabelWithStyle("InsteadMan", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	scr.Desc = widget.NewLabel("Выберите игру слева в списке")
	scr.Desc.Wrapping = fyne.TextWrapWord

	descScroll := widget.NewVScrollContainer(
		scr.Desc,
	)
	// descScroll.SetMinSize(fyne.NewSize(0, 100))

	scr.Hyperlink = widget.NewHyperlink("Website", nil)
	scr.Hyperlink.Hide()
	scr.InstallButton = widget.NewButtonWithIcon("Install", theme.ContentAddIcon(), func() {
		scr.Manager.InstallGame(scr.Game, nil)
	})
	scr.InstallButton.Style = widget.PrimaryButton
	scr.InstallButton.Hide()
	scr.RunButton = widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), func() {
		scr.Manager.RunGame(scr.Game)
	})
	scr.RunButton.Style = widget.PrimaryButton
	scr.RunButton.Hide()
	scr.DeleteButton = widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		scr.Manager.RemoveGame(scr.Game)
	})
	scr.DeleteButton.Hide()
	buttonsContainer := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		scr.InstallButton,
		scr.RunButton,
		scr.DeleteButton,
		scr.Hyperlink,
	)

	contentContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(scr.Title, buttonsContainer, nil, nil),
		descScroll,
		scr.Title,
		buttonsContainer,
	)

	allContainer := widget.NewVSplitContainer(scr.Image, contentContainer)

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil),
		allContainer,
	)

	return &scr
}
