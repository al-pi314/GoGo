package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"

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

func kthBiggest(scores []float64, k int) float64 {
	sort.Slice(scores, func(i, j int) bool {
		return scores[i] <= scores[j]
	})
	return scores[k]
}

func main() {
	loadConfig()
	hidden_layer := viper.GetIntSlice("HIDDEN_LAYER")
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
					InputNeurons:         3 * dymension * dymension,
					HiddenNeuronsByLayer: hidden_layer,
					OutputNeurons:        3,
				},
				ActivationFuncName: activation,
			}),
		})
		agents[uuid.NewString()] = &agent
	}

	var bestScore float64
	for i := 0; i <= matches; i++ {
		fmt.Printf("Starting match %d\n", i+1)
		agents, bestScore = match(agents, dymension)
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

	file, err := os.Open(agentsFile)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open agents file"))
	}
	defer file.Close()

	file.Write(data)
}

func match(agents map[string]*player.Agent, dymension int) (map[string]*player.Agent, float64) {
	// play games among players
	gamesScore := map[string]float64{}
	for blackPlayerID, blackPlayer := range agents {
		for whitePlayerID, whitePlayer := range agents {
			if blackPlayerID == whitePlayerID {
				continue
			}
			// create new game
			g := game.NewGame(game.Game{
				Dymension:   dymension,
				WhitePlayer: whitePlayer,
				BlackPlayer: blackPlayer,
			})
			// play game
			for g.Update() == nil {
			}
			// add score to the winner
			winner := blackPlayerID
			score := g.Score() + 0.5
			if score >= 0 {
				winner = whitePlayerID
			}
			gamesScore[winner] += math.Abs(score)
		}
	}
	// find best player
	scores := make([]float64, 0, len(gamesScore))
	for _, v := range gamesScore {
		scores = append(scores, v)
	}
	limit := kthBiggest(scores, len(scores)/2)

	newAgents := map[string]*player.Agent{}
	bestScore := 0.0
	for k, v := range gamesScore {
		if v >= limit {
			newAgents[uuid.NewString()] = agents[k].Offsprint()
			newAgents[uuid.NewString()] = agents[k].Offsprint()
		}
		if v >= bestScore {
			bestScore = v
		}
	}
	return newAgents, bestScore
}
