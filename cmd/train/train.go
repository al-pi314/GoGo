package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/al-pi314/gogo"
	"github.com/al-pi314/gogo/game"
	"github.com/al-pi314/gogo/nn"
	"github.com/al-pi314/gogo/player"
	"github.com/google/uuid"
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

func main() {
	config := loadConfig()

	// fill population
	agents := map[string]*player.Agent{}
	for len(agents) < config.PopulationSize {
		agent := player.NewAgent(player.Agent{
			StabilizationRate: config.StabilizationRate,
			MutationRate:      config.MutationRate,
			Logic: nn.NewNeuralNetwork(nn.NeuralNetwork{
				Structure: nn.Structure{
					InputNeurons:         3*config.Dymension*config.Dymension + gogo.GameStateSize(),
					HiddenNeuronsByLayer: config.HiddenLayers,
					OutputNeurons:        config.Dymension*config.Dymension + 1,
				},
				ActivationFuncName: config.Activation,
			}),
		})
		agents[uuid.NewString()] = &agent
	}

	var bestScore float64
	for i := 0; i <= config.Matches; i++ {
		fmt.Printf("Starting match %d\n", i+1)
		agents, bestScore = playTurnament(agents, config.Dymension)
		fmt.Printf("Match %d finished by %d players. Best Score %.2f\n", i+1, len(agents), bestScore)
	}

	agentsList := make([]*player.Agent, 0, len(agents))
	for _, v := range agents {
		agentsList = append(agentsList, v)
	}

	data, err := json.Marshal(agentsList)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to marshal agents"))
	}

	file, err := os.OpenFile(config.AgentsFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		log.Print(errors.Wrap(err, "failed to open env agents file"))
		file, err = os.OpenFile("./agents.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to open or create ./agents file"))
		}
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to write to file"))
	}
	fmt.Println(file.Fd())
}

func playTurnament(agents map[string]*player.Agent, dymension int) (map[string]*player.Agent, float64) {
	gameScores := map[string]float64{}
	for blackPlayerID, blackPlayer := range agents {
		for whitePlayerID, whitePlayer := range agents {
			if blackPlayerID == whitePlayerID {
				continue
			}
			score := playGame(game.Game{
				Dymension:   dymension,
				WhitePlayer: whitePlayer,
				BlackPlayer: blackPlayer,
			})

			gameScores[whitePlayerID] += score
			gameScores[blackPlayerID] += (-score)
		}
	}
	return findBest(agents, gameScores)
}

func playGame(newGame game.Game) float64 {
	// create new game
	g := game.NewGame(newGame)
	// play game
	for g.Update() == nil {
	}
	// return adjusted score
	return g.Score() / math.Max(1, float64(g.Moves()))
}

type PlayerPreformanceLinked = gogo.LinkedList[PlayerPreformance]

type PlayerPreformance struct {
	Agent *player.Agent
	Score float64
}

func (pp PlayerPreformance) Less(other interface{}) bool {
	switch otherType := other.(type) {
	case PlayerPreformance:
		return pp.Score < otherType.Score
	}
	return false

}

func findBest(agents map[string]*player.Agent, gameScores map[string]float64) (map[string]*player.Agent, float64) {
	var bestPlayer PlayerPreformanceLinked
	var bestPlayerSet bool
	avgScore := 0.0
	bestScore := 0.0
	for playerID, score := range gameScores {
		if !bestPlayerSet {
			bestPlayer = PlayerPreformanceLinked{
				Element: PlayerPreformance{
					Agent: agents[playerID],
					Score: score,
				},
			}
			bestPlayerSet = true
			continue
		}
		bestPlayer = bestPlayer.Add(PlayerPreformance{
			Agent: agents[playerID],
			Score: score,
		})
		avgScore += score
		if score > bestScore {
			bestScore = score
		}
		fmt.Printf("current score %f; best score %f; avf score %f\n", score, bestScore, avgScore)
	}
	avgScore /= float64(len(gameScores))

	newAgents := map[string]*player.Agent{}
	for i := 0; i < len(agents)/2; i++ {
		fmt.Printf("choosing agent with score %f\n", bestPlayer.Element.Score)
		newAgents[uuid.NewString()] = bestPlayer.Element.Agent.Offsprint()
		newAgents[uuid.NewString()] = bestPlayer.Element.Agent.Offsprint()
		if bestPlayer.Next != nil {
			bestPlayer = *bestPlayer.Next
		}
	}
	return newAgents, avgScore
}
