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

// Place implements player logic for placing their piece. Returns wether to place the piece or not, piece position and weather to skip move.
func (p *Human) Place(board [][]*bool) (bool, int, int, bool) {
	if ebiten.IsFocused() {
		switch {
		case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
			x, y := ebiten.CursorPosition()
			return true, x / p.XSnap, y / p.YSnap, false
		case inpututil.IsKeyJustPressed(ebiten.KeySpace):
			return false, -1, -1, true
		}
	}
	return false, -1, -1, false
}
