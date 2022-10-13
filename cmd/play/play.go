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

func loadConfig() *gogo.Config {
	viper.SetEnvPrefix("X")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	config := gogo.Config{}
	viper.Unmarshal(&config)

	rand.Seed(config.RandomSeed)
	return &config
}

func main() {
	config := loadConfig()

	whitePlayer := player.NewHuman(player.Human{
		XSnap: config.SquareSize + config.BorderSize,
		YSnap: config.SquareSize + config.BorderSize,
	})
	blackPlayer := player.NewAgent(player.Agent{
		Logic: nn.NewNeuralNetwork(nn.NeuralNetwork{
			Structure: nn.Structure{
				InputNeurons:         3*config.Dymension*config.Dymension + gogo.GameStateSize(),
				HiddenNeuronsByLayer: config.HiddenLayers,
				OutputNeurons:        config.Dymension*config.Dymension + 1,
			},
			ActivationFuncName: config.Activation,
		},
		),
	})

	game := game.NewGame(game.Game{
		Dymension:   config.Dymension,
		SquareSize:  config.SquareSize,
		BorderSize:  config.BorderSize,
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
