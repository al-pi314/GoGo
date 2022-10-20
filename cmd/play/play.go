package main

import (
	"flag"
	"log"
	"math/rand"

	"github.com/al-pi314/gogo"
	"github.com/al-pi314/gogo/game"
	"github.com/al-pi314/gogo/nn"
	"github.com/al-pi314/gogo/player"
	"github.com/al-pi314/gogo/population"
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

func createPlayer(config *gogo.Config, population *population.Population, t *string) player.Player {
	if *t == "agent" {
		if len(population.Enteties) == 0 {
			p := player.NewAgent(player.Agent{
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
			return &p
		}
		p := player.NewAgent(*population.Enteties[0].Agent)
		return &p
	}
	p := player.NewHuman(player.Human{
		XSnap: config.SquareSize + config.BorderSize,
		YSnap: config.SquareSize + config.BorderSize,
	})
	return &p
}

func main() {
	config := loadConfig()

	w := flag.String("white", "human", "set to 'human' or to 'agent'")
	b := flag.String("black", "human", "set to 'human' or to 'agent'")
	a := flag.String("agents", config.OutputFile, "agents file to use for agent players")
	flag.Parse()

	population := population.NewPopulation(config)
	if a != nil && *a != "" {
		population.LoadFromFile(a)
	}
	whitePlayer := createPlayer(config, population, w)
	blackPlayer := createPlayer(config, population, b)

	game := game.NewGame(game.Game{
		Dymension:   config.Dymension,
		SquareSize:  config.SquareSize,
		BorderSize:  config.BorderSize,
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
