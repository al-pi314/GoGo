package player

import (
	"reflect"

	"github.com/al-pi314/gogo/nn"
	"gonum.org/v1/gonum/mat"
)

type Agent struct {
	StabilizationRate float64
	MutationRate      float64
	Logic             *nn.NeuralNetwork
	SuggestedOnMove   int
	SuggestedMoves    *MoveSuggestion
}

func NewAgent(p Agent) Agent {
	return p
}

func (p *Agent) IsHuman() bool {
	return false
}

func encode(v interface{}) float64 {
	switch t := v.(type) {
	case bool:
		val := 0.0
		if t {
			val = 1.0
		}
		return val
	case int:
		return float64(t)
	case float64:
		return t
	default:
		return 0.0
	}
}

func encodeState(state *GameState) *mat.Dense {
	raw := []float64{}
	addEqualityResult := func(a, b *bool) {
		if (a == nil) != (b == nil) {
			raw = append(raw, 0)
			return
		}

		if a == b || *a == *b {
			raw = append(raw, 1)
			return
		}

		raw = append(raw, 0)
	}
	boolFalse := false
	boolTrue := true

	// encode game state
	for i := range state.Board {
		for j := range state.Board[i] {
			addEqualityResult(state.Board[i][j], nil)
			addEqualityResult(state.Board[i][j], &boolFalse)
			addEqualityResult(state.Board[i][j], &boolTrue)
		}
	}

	// encode game state
	t := reflect.TypeOf(*state)
	v := reflect.ValueOf(*state)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("encode")
		if tag == "true" && v.Field(i).CanInterface() {
			raw = append(raw, encode(v.Field(i).Interface()))
		}
	}

	return mat.NewDense(1, len(raw), raw)
}

type MoveSuggestion struct {
	X              int
	Y              int
	Effectivness   float64
	NextSuggestion *MoveSuggestion
}

func (ms *MoveSuggestion) add(suggestion MoveSuggestion) MoveSuggestion {
	// fix head
	if ms.Effectivness < suggestion.Effectivness {
		c := *ms
		suggestion.NextSuggestion = &c
		return suggestion
	}

	// fix current
	if ms.NextSuggestion == nil || ms.NextSuggestion.Effectivness < suggestion.Effectivness {
		suggestion.NextSuggestion = ms.NextSuggestion
		ms.NextSuggestion = &suggestion
		return *ms
	}

	// fix recursively
	ms.NextSuggestion.add(suggestion)
	return *ms
}

func interperate(output *mat.Dense, dymension int) (bool, *MoveSuggestion) {
	if output == nil {
		return false, nil
	}

	suggestions := MoveSuggestion{
		X:            0,
		Y:            0,
		Effectivness: output.At(0, 0),
	}
	for y := 0; y < dymension; y++ {
		for x := 1; x < dymension; x++ {
			suggestions = suggestions.add(MoveSuggestion{
				X:            x,
				Y:            y,
				Effectivness: output.At(0, y*dymension+x),
			})
		}
	}

	return output.At(0, dymension*dymension) >= 0.5, &suggestions
}

func (p *Agent) Offsprint() *Agent {
	return &Agent{
		StabilizationRate: p.StabilizationRate,
		MutationRate:      (1 - p.StabilizationRate) * p.MutationRate,
		Logic:             p.Logic.Mutate(p.MutationRate),
	}
}

func (p *Agent) Place(state *GameState) (bool, *int, *int) {
	if state == nil {
		return false, nil, nil
	}

	if p.SuggestedOnMove != state.Moves {
		// refresh cached moves suggestions
		var skip bool
		result := p.Logic.Predict(encodeState(state))
		skip, p.SuggestedMoves = interperate(result, len(state.Board))
		if skip || p.SuggestedMoves == nil {
			return true, nil, nil
		}
		p.SuggestedOnMove = state.Moves
	}

	// no more suggeste moves means no possible moves
	if p.SuggestedMoves == nil {
		return true, nil, nil
	}

	// pick best move from cached suggestions
	bestMove := *p.SuggestedMoves
	p.SuggestedMoves = bestMove.NextSuggestion

	return false, &bestMove.X, &bestMove.Y
}
