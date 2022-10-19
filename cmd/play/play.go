package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"

	"github.com/al-pi314/gogo"
	"github.com/al-pi314/gogo/game"
	"github.com/al-pi314/gogo/nn"
	"github.com/al-pi314/gogo/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
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

func loadAgents(a *string) []player.Agent {
	if a == nil {
		return nil
	}

	data, err := os.ReadFile(*a)
	if err != nil {
		log.Print(errors.Wrap(err, "invalid agents file provided"))
		return nil
	}
	agents := []player.Agent{}
	err = json.Unmarshal(data, &agents)
	if err != nil {
		log.Print(errors.Wrap(err, "invalid file structure"))
		return nil
	}
	return agents
}

func createPlayer(config *gogo.Config, agents []player.Agent, t *string) player.Player {
	if *t == "agent" {
		if len(agents) == 0 {
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
		p := player.NewAgent(agents[rand.Intn(len(agents))])
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

	agents := loadAgents(a)
	whitePlayer := createPlayer(config, agents, w)
	blackPlayer := createPlayer(config, agents, b)

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
