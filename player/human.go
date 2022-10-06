package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Human struct {
	XSnap int
	YSnap int
}

func NewHuman(p Human) Player {
	return &p
}

func (p *Human) Place(board [][]*bool) (bool, int, int, bool) {
	if ebiten.IsFocused() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		return true, x / p.XSnap, y / p.YSnap, false
	}
	return false, -1, -1, false
}
