package player

import (
	"github.com/al-pi314/gogo/nn"
	"gonum.org/v1/gonum/mat"
)

type Agent struct {
	Logic *nn.NeuralNetwork
}

func NewAgent(p Agent) Player {
	return &p
}

func encodeBoard(board [][]*bool) *mat.Dense {
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

	for i := range board {
		for j := range board[i] {
			addEqualityResult(board[i][j], nil)
			addEqualityResult(board[i][j], &boolFalse)
			addEqualityResult(board[i][j], &boolTrue)
		}
	}

	return mat.NewDense(1, len(raw), raw)
}

func interperate(output *mat.Dense, columns, rows int) (bool, *int, *int) {
	if output == nil {
		return false, nil, nil
	}

	x := int(float64(columns) * output.At(0, 1))
	y := int(float64(rows) * output.At(0, 1))
	return output.At(0, 0) >= 0.5, &x, &y
}

func (p *Agent) Place(board [][]*bool) (bool, *int, *int) {
	if len(board) == 0 || board == nil {
		return false, nil, nil
	}

	result := p.Logic.Predict(encodeBoard(board))
	return interperate(result, len(board), len(board[0]))
}
