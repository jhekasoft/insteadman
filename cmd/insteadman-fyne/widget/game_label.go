package widget

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/core/manager"
)

type GameLabel struct {
	widget.Label
	Game     *manager.Game
	OnTapped func() `json:"-"`
}

func NewGameLabel(game *manager.Game, tapped func()) *GameLabel {
	item := &GameLabel{}
	item.ExtendBaseWidget(item)
	item.SetText(game.Title)
	item.TextStyle = fyne.TextStyle{Bold: game.Installed}
	// item.Alignment = widget.ButtonAlignLeading
	item.Wrapping = fyne.TextWrapWord
	item.Game = game
	item.OnTapped = tapped

	return item
}

func (l *GameLabel) Tapped(_ *fyne.PointEvent) {
	l.OnTapped()
}

func (l *GameLabel) TappedSecondary(_ *fyne.PointEvent) {
}
