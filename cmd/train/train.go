package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func loadConfig() {
	viper.SetEnvPrefix("X")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	rand.Seed(viper.GetInt64("RANDOM_SEED"))
}

func main() {
	loadConfig()
	hidden_layer := viper.GetIntSlice("HIDDEN_LAYERS")
	activation := viper.GetString("ACTIVATION")
	populationSize := viper.GetInt("POPULATION_SIZE")
	mutationRate := viper.GetFloat64("MUTATION_RATE")
	stabilizationRate := viper.GetFloat64("STABLIZATION_RATE")
	dymension := viper.GetInt("DYMENSION")
	matches := viper.GetInt("MATCHES")
	agentsFile := viper.GetString("AGENTS_FILE")

	// fill population
	agents := map[string]*player.Agent{}
	for len(agents) < populationSize {
		agent := player.NewAgent(player.Agent{
			StabilizationRate: stabilizationRate,
			MutationRate:      mutationRate,
			Logic: nn.NewNeuralNetwork(nn.NeuralNetwork{
				Structure: nn.Structure{
					InputNeurons:         3*dymension*dymension + gogo.GameStateSize(),
					HiddenNeuronsByLayer: hidden_layer,
					OutputNeurons:        dymension*dymension + 1,
				},
				ActivationFuncName: activation,
			}),
		})
		agents[uuid.NewString()] = &agent
	}

	var bestScore float64
	for i := 0; i <= matches; i++ {
		fmt.Printf("Starting match %d\n", i+1)
		agents, bestScore = playTurnament(agents, dymension)
		fmt.Printf("Match %d finished by %d players. Best Score %.2f\n", i+1, len(agents), bestScore)
	}

	agentsList := make([]*player.Agent, 0, len(agents))
	for _, v := range agents {
		agentsList = append(agentsList, v)
	}

	data, err := json.Marshal(agentsList)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to masthal agents"))
	}

	file, err := os.OpenFile(agentsFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open agents file"))
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
			if score >= 0 {
				gameScores[whitePlayerID] += score
			} else {
				gameScores[blackPlayerID] += (-score)
			}

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
	return g.Score()
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
	for playerID, score := range gameScores {
		if !bestPlayerSet {
			bestPlayer = PlayerPreformanceLinked{
				Element: PlayerPreformance{
					Agent: agents[playerID],
					Score: score,
				},
			}
			continue
		}
		bestPlayer.Add(PlayerPreformance{
			Agent: agents[playerID],
			Score: score,
		})
	}

	newAgents := map[string]*player.Agent{}
	bestScore := bestPlayer.Element.Score
	for i := 0; i < len(agents)/2; i++ {
		newAgents[uuid.NewString()] = bestPlayer.Element.Agent.Offsprint()
		newAgents[uuid.NewString()] = bestPlayer.Element.Agent.Offsprint()
		if bestPlayer.Next != nil {
			bestPlayer = *bestPlayer.Next
		}
	}
	return newAgents, bestScore
}
