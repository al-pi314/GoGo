package main

import (
	"flag"
	"fmt"
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
	moveDelay := flag.Int("delay", 0, "miliseconds to wait after each move not made by human")
	replay := flag.String("replay", "", "game save file to replay")
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

		if isArgSet(white) && *white == "agent" {
			whitePlayer = population.FirstNthAgent(0)
		}

		if isArgSet(black) && *black == "agent" {
			blackPlayer = population.FirstNthAgent(1)
		}
	}

	game := game.NewGame(game.Game{
		Dymension:   config.Dymension,
		SquareSize:  config.SquareSize,
		BorderSize:  config.BorderSize,
		WhitePlayer: whitePlayer,
		BlackPlayer: blackPlayer,
		MoveDelay:   moveDelay,
	})

	if isArgSet(replay) {
		game.ReplayFromFile(*replay)
		if moveDelay == nil {
			fmt.Println("WARRNING: move delay not set, use -delay flag to set replay pace")
		}
	}

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(game.Size())
	ebiten.SetWindowTitle("GoGo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
