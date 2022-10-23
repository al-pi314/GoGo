package population

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/al-pi314/gogo"
	"github.com/al-pi314/gogo/game"
	"github.com/al-pi314/gogo/nn"
	"github.com/al-pi314/gogo/player"
	"github.com/pkg/errors"
)

type Population struct {
	GameDymension   int
	Enteties        []*Entety
	Age             int
	Size            int
	OutputDirectory string

	outputFileName string   `json:"-"`
	file           *os.File `json:"-"`
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
		GameDymension:   config.Dymension,
		Enteties:        []*Entety{},
		OutputDirectory: strings.TrimSuffix(config.OutputDirectory, "/"),
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

	p.outputFileName = fmt.Sprintf("%s/population.json", p.OutputDirectory)
	if err := os.Mkdir(fmt.Sprintf("%s/games", p.OutputDirectory), os.ModePerm); err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(errors.Wrap(err, "failed to create games directory inside of population directory"))
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

func (p *Population) OpenOutputFile(filePath string) {
	var err error
	p.file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		log.Print(errors.Wrap(err, "failed to open output file!"))
	} else {
		p.outputFileName = filePath
	}
}

func (p *Population) Save() {
	if p.file == nil {
		p.OpenOutputFile(p.outputFileName)
	}
	defer func() {
		p.file.Close()
		p.file = nil
	}()

	now := time.Now()
	bytes, err := json.Marshal(TrainingSave{
		Time:       &now,
		Population: p,
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not marshal population"))
	}

	n, err := p.file.Write(bytes)
	if err != nil || n != len(bytes) {
		log.Fatal(errors.Wrap(err, fmt.Sprintf("writting error or the write was incomplete (attempted to write %d bytes, written %d bytes)", len(bytes), n)))
	}
	fmt.Println("saved population!")
}

type TrainingSettings struct {
	Rounds    int
	Groups    int
	KeepBestN int

	SaveInterval     int
	SaveGameInterval int
}

func (p *Population) Train(settings TrainingSettings) {
	for i := 0; i < settings.Rounds; i++ {
		s := time.Now().UnixMilli()
		fmt.Println("-------------------------------------")
		fmt.Printf("starting round (population age %d) %d\n", p.Age, i)

		// divide population into groups
		groups := p.CreateGroups(settings.Groups)

		// play games among agents inside groups and select top N
		groupsBest := [][]*Entety{}
		saveBestGames := (i+1)%settings.SaveGameInterval == 0
		for i, group := range groups {
			groupsBest = append(groupsBest, p.playMatches(i, group, settings.KeepBestN, saveBestGames))
			fmt.Println("-------------")
		}

		// crossover and mutate group winners to create new population
		p.newPopulation(groupsBest, settings.KeepBestN)

		d := time.Now().UnixMilli() - s
		fmt.Printf("finished round %d (miliseconds spent %d)\n", i, d)

		// check save intervals
		if (i+1)%settings.SaveInterval == 0 {
			p.Save()
		}
	}
}

func (p *Population) CreateGroups(groups int) [][]*Entety {
	groupSize := p.Size / groups
	result := [][]*Entety{}
	for i := 0; i < groups; i++ {
		result = append(result, p.Enteties[i*groupSize:(i+1)*groupSize])
	}
	return result
}

func (p *Population) playMatches(groupID int, enteties []*Entety, toKeep int, saveBest bool) []*Entety {
	s := time.Now().UnixMilli()
	fmt.Printf("[group %d] starting group matches\n", groupID)
	var best *float64
	var bestGame *game.Game
	gameName := ""
	for idOne, entetyOne := range enteties {
		for idTwo, entetyTwo := range enteties {
			if idOne == idTwo {
				continue
			}
			score, game := entetyOne.match(entetyTwo, p.GameDymension)
			if best == nil || *best > math.Abs(score) {
				best = &score
				bestGame = game
				gameName = fmt.Sprintf("group_%d_%d_%d_%d.json", p.Age, groupID, idOne, idTwo)
			}
		}
	}
	if saveBest {
		bestGame.SaveFileName = fmt.Sprintf("%s/games/%s", p.OutputDirectory, gameName)
		bestGame.Save()
	}

	sort.Slice(enteties, func(i, j int) bool {
		return enteties[i].Score >= enteties[j].Score
	})
	fmt.Printf("[group %d] best entety score %f\n", groupID, enteties[0].Score)
	fmt.Printf("[group %d] finished group matches (miliseconds spent %d)\n", groupID, time.Now().UnixMilli()-s)
	return enteties[:toKeep]
}

func (p *Population) newPopulation(groups [][]*Entety, groupsSize int) {
	prevSize := p.Size
	p.Enteties = []*Entety{}
	p.Size = 0

	groupIdx := 0
	inGroupIdx := 0
	for p.Size < prevSize {
		// select parents
		entetyOne := groups[groupIdx][inGroupIdx]
		entetyTwo := groups[rand.Intn(len(groups))][rand.Intn(groupsSize)]

		// crossover & mutate for new entety
		newEntety := Entety{
			Agent: entetyOne.Agent.Crossover(entetyTwo.Agent),
		}
		p.AddEntety(&newEntety)

		// select next entety
		groupIdx++
		if groupIdx >= len(groups) {
			groupIdx = 0
			inGroupIdx = (inGroupIdx + 1) % groupsSize
		}
	}

	p.Age++
}

func (p *Population) BestNPlayer(n int) *player.Agent {
	if n >= p.Size {
		return nil
	}

	return p.Enteties[n].Agent
}

func (e *Entety) match(o *Entety, gameDymension int) (float64, *game.Game) {
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
	return score, g
}
