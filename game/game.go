package game

import (
	"image/color"

	"github.com/al-pi314/gogo/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Cordinate struct {
	X int
	Y int
}

type Game struct {
	Rows        int
	Columns     int
	SquareSize  int
	BorderSize  int
	WhitePlayer player.Player
	BlackPlayer player.Player

	board       [][]*bool
	delay_lock  bool
	locked      Cordinate
	whiteToMove bool
}

func NewGame(g Game) *Game {
	g.board = make([][]*bool, g.Rows)
	for y := range g.board {
		g.board[y] = make([]*bool, g.Columns)
	}
	g.locked = Cordinate{-1, -1}

	return &g
}

// ------------------------------------ Helper Functions ------------------------------------ \\
func (g *Game) Size() (int, int) {
	return g.Columns*(g.SquareSize+g.BorderSize) + g.BorderSize, g.Rows*(g.SquareSize+g.BorderSize) + g.BorderSize
}

// drawSquare draws a square on the image.
func (g *Game) drawSquare(screen *ebiten.Image, x1, y1, x2, y2 int, clr color.Color) {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			screen.Set(x, y, clr)
		}
	}
}

func (g *Game) pieceAt(x, y int) *bool {
	if y >= len(g.board) || y < 0 || x >= len(g.board[y]) || x < 0 {
		return nil
	}
	return g.board[y][x]
}

func (g *Game) hasRoom(x, y int, white bool, checked map[Cordinate]bool, group *[]Cordinate) bool {
	c := Cordinate{x, y}
	if checked, ok := checked[c]; ok && checked {
		return false
	}
	checked[c] = true

	if y >= len(g.board) || y < 0 || x >= len(g.board[y]) || x < 0 {
		return false
	}

	if g.board[y][x] == nil {
		return true
	}

	if *g.board[y][x] != white {
		return false
	}

	// record group
	if group != nil {
		*group = append(*group, c)
	}

	// execute all checks
	results := []bool{g.hasRoom(x-1, y, white, checked, group), g.hasRoom(x+1, y, white, checked, group), g.hasRoom(x, y-1, white, checked, group), g.hasRoom(x, y+1, white, checked, group)}
	for _, r := range results {
		if r {
			return true
		}
	}
	return false
}

// ------------------------------------ ----------------- ------------------------------------ \\

// -------------------------------------- Game Functions ------------------------------------- \\
func (g *Game) placePiece(x, y int, white bool) bool {
	// out of bounds or already occupied spaces are invalid
	if y >= len(g.board) || y < 0 || x >= len(g.board[y]) || x < 0 || g.board[y][x] != nil {
		return false
	}
	// place piece
	g.board[y][x] = &white

	// check for opponent group eliminations
	if g.caputreOpponent(x, y, white) {
		return true
	}

	// would be eliminated when placed
	hasRoom := g.hasRoom(x, y, white, map[Cordinate]bool{}, nil)
	if !hasRoom {
		g.board[y][x] = nil
		return false
	}

	return true
}

func (g *Game) caputreOpponent(x, y int, white bool) bool {
	toRemove := []Cordinate{}
	if g := g.findGroup(x-1, y, !white); g != nil {
		toRemove = append(toRemove, g...)
	}
	if g := g.findGroup(x+1, y, !white); g != nil {
		toRemove = append(toRemove, g...)
	}
	if g := g.findGroup(x, y-1, !white); g != nil {
		toRemove = append(toRemove, g...)
	}
	if g := g.findGroup(x, y+1, !white); g != nil {
		toRemove = append(toRemove, g...)
	}

	// ko rule
	if len(toRemove) == 1 {
		if toRemove[0].X == g.locked.X && toRemove[0].Y == g.locked.Y {
			return false
		}
		g.delay_lock = true
		g.locked.X = x
		g.locked.Y = y
	}

	for _, c := range toRemove {
		g.board[c.Y][c.X] = nil
	}
	return len(toRemove) != 0
}

func (g *Game) findGroup(x, y int, white bool) []Cordinate {
	pieceColor := g.pieceAt(x, y)
	if pieceColor == nil || *pieceColor != white {
		return nil
	}

	group := []Cordinate{}
	if g.hasRoom(x, y, white, map[Cordinate]bool{}, &group) {
		return nil
	}

	return group
}

// -------------------------------------- -------------- ------------------------------------- \\

// --------------------------- Functions required by ebiten engine --------------------------- \\
func (g *Game) Update() error {
	player := g.WhitePlayer
	if !g.whiteToMove {
		player = g.BlackPlayer
	}

	place, x, y, skip := player.Place(g.board)
	if place && (skip || g.placePiece(x, y, g.whiteToMove)) {
		// lock unlocks after next successful move
		if !g.delay_lock {
			g.locked.X = -1
			g.locked.Y = -1
		}
		g.delay_lock = false
		// change player to move
		g.whiteToMove = !g.whiteToMove
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw board - squares with left and top borders
	for x := 0; x <= g.Columns+1; x++ {
		for y := 0; y <= g.Rows+1; y++ {
			x1 := x * (g.BorderSize + g.SquareSize)
			y1 := y * (g.BorderSize + g.SquareSize)
			// draw left border
			if y <= g.Rows {
				x2 := x1 + g.BorderSize
				y2 := y1 + g.SquareSize + g.BorderSize
				g.drawSquare(screen, x1, y1, x2, y2, color.RGBA{160, 175, 190, 1})
			}
			// draw top border
			if x <= g.Columns {
				x2 := x1 + g.SquareSize + g.BorderSize
				y2 := y1 + g.BorderSize
				g.drawSquare(screen, x1, y1, x2, y2, color.RGBA{160, 175, 190, 1})
			}

			// draw empty square when inside the board
			if x <= g.Columns && y <= g.Rows {
				x1 += g.BorderSize
				y1 += g.BorderSize
				x2 := x1 + g.SquareSize
				y2 := y1 + g.SquareSize
				g.drawSquare(screen, x1, y1, x2, y2, color.RGBA{180, 90, 30, 1})
			}
		}
	}

	// draw pieces
	for piece_y, row := range g.board {
		for piece_x, piece := range row {
			if piece == nil {
				continue
			}
			x := (piece_x+1)*(g.SquareSize+g.BorderSize) - g.SquareSize/2
			y := (piece_y+1)*(g.SquareSize+g.BorderSize) - g.SquareSize/2
			clr := color.White
			if !*piece {
				clr = color.Black
			}
			ebitenutil.DrawCircle(screen, float64(x), float64(y), float64(g.SquareSize/2)*0.8, clr)
		}

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Rows*(g.BorderSize+g.SquareSize) + g.BorderSize, g.Columns*(g.BorderSize+g.SquareSize) + g.BorderSize
}

// --------------------------- ------------------------------------ --------------------------- \\
