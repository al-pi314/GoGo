package population

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/al-pi314/gogo"
	"github.com/al-pi314/gogo/game"
	"github.com/al-pi314/gogo/nn"
	"github.com/al-pi314/gogo/player"
	"github.com/pkg/errors"
)

type Population struct {
	GameDymension int
	Enteties      []*Entety
	Age           int
	Size          int

	OutputFileName string   `json:"-"`
	File           *os.File `json:"-"`
}

type Entety struct {
	Agent *player.Agent
	Score float64
}

type TrainingSave struct {
	Time       *time.Time
	Population *Population
}

func NewPopulation(config *gogo.Config) *Population {
	p := Population{
		GameDymension: config.Dymension,
		Enteties:      []*Entety{},
	}
	for len(p.Enteties) < config.PopulationSize {
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
		p.AddEntety(&Entety{
			Agent: &agent,
		})
	}
	return &p
}

func (p *Population) LoadFromFile(filePath *string) bool {
	if filePath == nil {
		return false
	}

	data, err := os.ReadFile(*filePath)
	if err != nil {
		log.Print(errors.Wrap(err, "invalid agents file provided"))
		return false
	}
	saveData := TrainingSave{}
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		log.Print(errors.Wrap(err, "invalid file structure"))
		return false
	}

	for _, e := range saveData.Population.Enteties {
		agent := player.NewAgent(*e.Agent)
		e.Agent = &agent
	}

	fmt.Printf("...loaded population from file (population saved at %s)\n", saveData.Time.String())
	*p = *saveData.Population
	return true
}

func (p *Population) AddEntety(e *Entety) {
	p.Enteties = append(p.Enteties, e)
	p.Size++
}

func (p *Population) OpenOuptutFile(filePath string) {
	var err error
	p.File, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		log.Print(errors.Wrap(err, "failed to open output file!"))
	} else {
		p.OutputFileName = filePath
	}
}

func (p *Population) Save() {
	if p.File == nil {
		p.OpenOuptutFile(p.OutputFileName)
	}
	defer func() {
		p.File.Close()
		p.File = nil
	}()

	now := time.Now()
	bytes, err := json.Marshal(TrainingSave{
		Time:       &now,
		Population: p,
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not marshal population"))
	}

	n, err := p.File.Write(bytes)
	if err != nil || n != len(bytes) {
		log.Fatal(errors.Wrap(err, fmt.Sprintf("writting error or the write was incomplete (attempted to write %d bytes, written %d bytes)", len(bytes), n)))
	}
}

func (p *Population) Train(rounds int, saveInterval int, file *os.File) {
	for i := 0; i < rounds; i++ {
		s := time.Now().UnixMilli()
		fmt.Printf("starting round (population age %d) %d\n", p.Age, i)
		p.playMatches()
		p.fitnessSelection()
		p.Age++
		d := time.Now().UnixMilli() - s
		fmt.Printf("finished round %d (miliseconds spent %d)\n", i, d)

		if (i+1)%saveInterval == 0 {
			fmt.Println("saving population")
			p.Save()
		}
	}
}

func (p *Population) playMatches() {
	for idOne, entetyOne := range p.Enteties {
		for idTwo, entetyTwo := range p.Enteties {
			if idOne == idTwo {
				continue
			}
			entetyOne.match(entetyTwo, p.GameDymension)
		}
	}
}

func (p *Population) sortEnteites() {
	sort.Slice(p.Enteties, func(i, j int) bool {
		return p.Enteties[i].Score >= p.Enteties[j].Score
	})
}

func (p *Population) fitnessSelection() {
	p.sortEnteites()

	newEnteties := []*Entety{}
	averageScore := 0.0
	for i := 0; i < len(p.Enteties)/2; i++ {
		newEnteties = append(newEnteties, &Entety{
			Agent: p.Enteties[i].Agent.Offsprint(),
		})
		newEnteties = append(newEnteties, &Entety{
			Agent: p.Enteties[i].Agent.Offsprint(),
		})
		averageScore += p.Enteties[i].Score
	}
	averageScore /= float64(len(newEnteties) / 2)

	fmt.Println("fitness information:")
	fmt.Printf("...population size %d, selected best %d enteties\n", p.Size, len(p.Enteties)/2)
	fmt.Printf("...best entety score: %f\n", p.Enteties[0].Score)
	fmt.Printf("...worst entety score: %f\n", p.Enteties[len(p.Enteties)-1].Score)
	fmt.Printf("...worst selected score: %f\n", p.Enteties[len(p.Enteties)/2].Score)
	fmt.Printf("...selected enteties average score: %f\n", averageScore)
	p.Enteties = newEnteties
	p.Size = len(newEnteties)
}

func (p *Population) BestNPlayer(n int) *player.Agent {
	if n >= p.Size {
		return nil
	}

	return p.Enteties[n].Agent
}

func (e *Entety) match(o *Entety, gameDymension int) {
	g := game.NewGame(game.Game{
		Dymension:   gameDymension,
		WhitePlayer: e.Agent,
		BlackPlayer: o.Agent,
	})

	for g.Update() == nil {
		// play game
	}

	score := g.Score() / math.Max(1, float64(g.Moves()))
	e.Score += score
	o.Score += -score
}
