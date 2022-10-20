package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Human struct {
	XSnap int
	YSnap int
}

func NewHuman(p Human) *Human {
	return &p
}

func (p *Human) IsHuman() bool {
	return true
}

// Place implements player logic for placing their piece. Returns wether to place the piece or not, piece position and weather to skip move.
func (p *Human) Place(state *GameState) (bool, *int, *int) {
	if ebiten.IsFocused() {
		switch {
		case inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft):
			x, y := ebiten.CursorPosition()
			x /= p.XSnap
			y /= p.YSnap
			return false, &x, &y
		case inpututil.IsKeyJustPressed(ebiten.KeySpace):
			return true, nil, nil
		}
	}
	return false, nil, nil
}
