package player

type Agent struct{}

func NewAgent(p Agent) Player {
	return &p
}

func (p *Agent) Place(board [][]*bool) (bool, int, int, bool) {
	return false, -1, -1, false
}
