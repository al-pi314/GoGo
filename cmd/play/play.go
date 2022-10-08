package main

import (
	"log"
	"math/rand"

	"github.com/al-pi314/gogo/game"
	"github.com/al-pi314/gogo/nn"
	"github.com/al-pi314/gogo/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/viper"
)

func loadConfig() {
	viper.SetEnvPrefix("X")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

func main() {
	loadConfig()
	squareSize := viper.GetInt("SQUARE_SIZE")
	borderSize := viper.GetInt("BORDER_SIZE")
	rows := viper.GetInt("ROWS")
	columns := viper.GetInt("COLUMNS")

	rand.Seed(viper.GetInt64("RANDOM_SEED"))

	whitePlayer := player.NewHuman(player.Human{
		XSnap: squareSize + borderSize,
		YSnap: squareSize + borderSize,
	})
	blackPlayer := player.NewAgent(player.Agent{
		Logic: nn.NewNeuralNetwork(nn.Structure{
			InputNeurons:         3 * rows * columns,
			HiddenNeuronsByLayer: []int{50, 100, 50},
			OutputNeurons:        3,
		}),
	})

	game := game.NewGame(game.Game{
		Rows:        rows,
		Columns:     columns,
		SquareSize:  squareSize,
		BorderSize:  borderSize,
		WhitePlayer: whitePlayer,
		BlackPlayer: blackPlayer,
	})

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(game.Size())
	ebiten.SetWindowTitle("GoGo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
