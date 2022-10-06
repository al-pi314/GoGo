package main

import (
	"log"

	"github.com/al-pi314/gogo/game"
	"github.com/al-pi314/gogo/player"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	whitePlayer := player.NewHuman(player.Human{
		XSnap: 33,
		YSnap: 33,
	})
	blackPlayer := player.NewHuman(player.Human{
		XSnap: 33,
		YSnap: 33,
	})

	game := game.NewGame(game.Game{
		Rows:        19,
		Columns:     19,
		SquareSize:  30,
		BorderSize:  3,
		WhitePlayer: whitePlayer,
		BlackPlayer: blackPlayer,
	})

	ebiten.SetWindowSize(game.Size())
	ebiten.SetWindowTitle("GoGo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
