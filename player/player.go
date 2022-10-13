package player

import (
	"github.com/al-pi314/gogo"
)

type GameState = gogo.GameState

type Player interface {
	Place(*gogo.GameState) (bool, *int, *int)
	IsHuman() bool
}
