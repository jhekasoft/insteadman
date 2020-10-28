package screen

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
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
	Version       *widget.Label
	Lang          *widget.Label
	Repository    *widget.Label
	Screen        fyne.CanvasObject
	Container     *widget.SplitContainer
	Image         *widget.Icon
	Hyperlink     *widget.Hyperlink
	InstallButton *widget.Button
	RunButton     *widget.Button
	DeleteButton  *widget.Button
	Game          *manager.Game
	UpdateF       func()
}

func (scr *GameInfoScreen) UpdateInfo(g *manager.Game) {
	scr.Game = g

	scr.Title.SetText(g.Title)
	scr.Desc.SetText(g.Description)

	// Labels
	scr.Version.SetText(g.Version)
	scr.Lang.SetText(strings.Join(g.Languages, ", "))
	scr.Repository.SetText(g.RepositoryName)

	// URL
	if g.Descurl != "" {
		scr.Hyperlink.SetURLFromString(g.Descurl)
		scr.Hyperlink.Show()
	}

	// Buttons
	// TODO: add Update button
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

	// if scr.Container != nil {
	// 	scr.Container.Refresh()
	// }
}

func NewGameInfoScreen(
	m *manager.Manager,
	c *configurator.Configurator,
	mainIcon fyne.Resource,
	window fyne.Window) *GameInfoScreen {
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
		progDialog := dialog.NewProgress(scr.Game.Title, "Installing...", window)
		progDialog.Show()
		err := scr.Manager.InstallGame(scr.Game, func(size uint64) {
			percents := float64(size) / float64(scr.Game.Size)
			progDialog.SetValue(percents)
			if float64(size) >= float64(scr.Game.Size) {
				progDialog.SetValue(1)
				progDialog.Hide()
			}
		})

		if err != nil {
			progDialog.Hide()
			dialog.ShowError(err, window)
			return
		}

		scr.Game.Installed = true
		scr.UpdateInfo(scr.Game)

		if scr.UpdateF != nil {
			scr.UpdateF()
		}
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

		// TODO: Check error
		scr.Game.Installed = false
		scr.UpdateInfo(scr.Game)
		if scr.UpdateF != nil {
			scr.UpdateF()
		}
	})
	scr.DeleteButton.Hide()
	scr.Version = widget.NewLabel("")
	scr.Lang = widget.NewLabel("")
	scr.Repository = widget.NewLabel("")
	scr.Repository.Wrapping = fyne.TextWrapWord
	buttonsContainer := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		scr.InstallButton,
		scr.RunButton,
		scr.DeleteButton,
		scr.Hyperlink,
		scr.Version,
		scr.Lang,
		scr.Repository,
	)

	contentContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(scr.Title, buttonsContainer, nil, nil),
		descScroll,
		scr.Title,
		buttonsContainer,
	)

	scr.Container = widget.NewVSplitContainer(scr.Image, contentContainer)

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil),
		scr.Container,
	)

	return &scr
}
