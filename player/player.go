package player

type Player interface {
	Place(board [][]*bool) (bool, *int, *int)
	IsHuman() bool
}
