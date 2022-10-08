package player

import (
	"math/rand"

	"github.com/al-pi314/gogo/nn"
)

type Agent struct {
	Logic nn.NeuralNetwork
}

func NewAgent(p Agent) Player {
	return &p
}

func (p *Agent) Place(board [][]*bool) (bool, *int, *int) {
	y := rand.Intn(len(board))
	x := rand.Intn(len(board[y]))
	return false, &y, &x
}
