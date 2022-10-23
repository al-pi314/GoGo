package gogo

import "reflect"

type GameState struct {
	Board               [][]*bool
	Moves               [][2]*int
	MovesCount          int
	OpponentSkipped     bool `encode:"true"`
	BlackStones         int  `encode:"true"`
	WhiteStones         int  `encode:"true"`
	BlackStonesCaptured int  `encode:"true"`
	WhiteStonesCaptured int  `encode:"true"`
}

func GameStateSize() int {
	t := reflect.TypeOf(GameState{})
	v := reflect.ValueOf(GameState{})
	encodable := 0
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("encode")
		if tag == "true" && v.Field(i).CanInterface() {
			encodable++
		}
	}
	return encodable
}
