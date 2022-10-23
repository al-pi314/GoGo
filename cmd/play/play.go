package main

import (
	"flag"
	"log"
	"math/rand"

	"github.com/al-pi314/gogo"
	"github.com/al-pi314/gogo/game"
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

func isArgSet(arg *string) bool {
	return arg != nil && *arg != ""
}

func main() {
	config := loadConfig()

	white := flag.String("white", "human", "set to 'human' or to 'agent'")
	black := flag.String("black", "human", "set to 'human' or to 'agent'")
	populationFile := flag.String("population", "", "population file to use for agent players")
	moveDelay := flag.Int("delay", 0, "miliseconds to wait after each AI move")
	flag.Parse()

	var whitePlayer player.Player
	var blackPlayer player.Player
	whitePlayer = player.NewHuman(player.Human{
		XSnap: config.SquareSize + config.BorderSize,
		YSnap: config.SquareSize + config.BorderSize,
	})
	blackPlayer = player.NewHuman(player.Human{
		XSnap: config.SquareSize + config.BorderSize,
		YSnap: config.SquareSize + config.BorderSize,
	})

	population := population.NewPopulation(config)
	if isArgSet(populationFile) {
		population.LoadFromFile(populationFile)

		n := 0
		if isArgSet(white) && *white == "agent" {
			if agent := population.BestNPlayer(n); agent != nil {
				whitePlayer = agent
			}
			n++
		}

		if isArgSet(black) && *black == "agent" {
			if agent := population.BestNPlayer(n); agent != nil {
				blackPlayer = agent
			}
			n++
		}
	}

	game := game.NewGame(game.Game{
		Dymension:      config.Dymension,
		SquareSize:     config.SquareSize,
		BorderSize:     config.BorderSize,
		WhitePlayer:    whitePlayer,
		BlackPlayer:    blackPlayer,
		AgentMoveDelay: moveDelay,
	})

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(game.Size())
	ebiten.SetWindowTitle("GoGo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
