package main

import (
	"log"
	"math/rand"

	"github.com/al-pi314/gogo"
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
	activation := viper.GetString("ACTIVATION")
	hidden_layer := viper.GetIntSlice("HIDDEN_LAYER")

	whitePlayer := player.NewHuman(player.Human{
		XSnap: squareSize + borderSize,
		YSnap: squareSize + borderSize,
	})
	blackPlayer := player.NewAgent(player.Agent{
		Logic: nn.NewNeuralNetwork(nn.NeuralNetwork{
			Structure: nn.Structure{
				InputNeurons:         3*dymension*dymension + gogo.GameStateSize(),
				HiddenNeuronsByLayer: hidden_layer,
				OutputNeurons:        dymension*dymension + 1,
			},
			ActivationFuncName: activation,
		},
		),
	})

	game := game.NewGame(game.Game{
		Dymension:   dymension,
		SquareSize:  squareSize,
		BorderSize:  borderSize,
		WhitePlayer: &whitePlayer,
		BlackPlayer: &blackPlayer,
	})

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(game.Size())
	ebiten.SetWindowTitle("GoGo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
