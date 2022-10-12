package main

import (
	"log"
	"math"
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
	rand.Seed(viper.GetInt64("RANDOM_SEED"))
}

func main() {
	loadConfig()
	squareSize := viper.GetInt("SQUARE_SIZE")
	borderSize := viper.GetInt("BORDER_SIZE")
	dymension := viper.GetInt("DYMENSION")

	whitePlayer := player.NewHuman(player.Human{
		XSnap: squareSize + borderSize,
		YSnap: squareSize + borderSize,
	})
	blackPlayer := player.NewAgent(player.Agent{
		Logic: nn.NewNeuralNetwork(nn.NeuralNetwork{
			Structure: nn.Structure{
				InputNeurons:         3 * dymension * dymension,
				HiddenNeuronsByLayer: []int{50, 100, 50},
				OutputNeurons:        3,
			},
			ActivationFunc: func(v float64) float64 { return float64(dymension) / (1.0 + math.Exp(-v)) },
		},
		),
	})

	game := game.NewGame(game.Game{
		Dymension:   dymension,
		SquareSize:  squareSize,
		BorderSize:  borderSize,
		WhitePlayer: whitePlayer,
		BlackPlayer: &blackPlayer,
	})

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(game.Size())
	ebiten.SetWindowTitle("GoGo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
