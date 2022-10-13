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

func interperate(output *mat.Dense, dymension int) (bool, *int, *int) {
	if output == nil {
		return false, nil, nil
	}

	x := int(float64(dymension) * output.At(0, 1))
	y := int(float64(dymension) * output.At(0, 1))
	return output.At(0, 0) >= 0.5, &x, &y
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

	result := p.Logic.Predict(encodeState(state))
	return interperate(result, len(state.Board))
}
