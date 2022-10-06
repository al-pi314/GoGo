package main

import (
	"log"

	"github.com/al-pi314/GoGo/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := game.NewGame(game.Game{
		Rows:       19,
		Columns:    19,
		SquareSize: 30,
		BorderSize: 3,
	})

	game.PlacePiece(0, 0, true)
	game.PlacePiece(8, 7, true)
	game.PlacePiece(7, 8, false)
	game.PlacePiece(18, 18, false)

	ebiten.SetWindowSize(game.Size())
	ebiten.SetWindowTitle("GoGo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
