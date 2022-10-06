package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Cordinate struct {
	X int
	Y int
}

type Game struct {
	Rows       int
	Columns    int
	SquareSize int
	BorderSize int
	board      [][]*bool
}

func NewGame(g Game) *Game {
	g.board = make([][]*bool, g.Rows)
	for y := range g.board {
		g.board[y] = make([]*bool, g.Columns)
	}

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

func (g *Game) hasRoom(x, y int, white bool, checked map[Cordinate]bool) bool {
	c := Cordinate{x, y}
	if checked, ok := checked[c]; ok && checked {
		return false
	}
	checked[c] = true

	if y >= len(g.board) || x >= len(g.board[y]) {
		return false
	}

	if g.board[y][x] == nil {
		return true
	}

	if *g.board[y][x] != white {
		return false
	}

	return g.hasRoom(x-1, y, white, checked) || g.hasRoom(x+1, y, white, checked) || g.hasRoom(x, y-1, white, checked) || g.hasRoom(x, y+1, white, checked)
}

// ------------------------------------ ----------------- ------------------------------------ \\

// -------------------------------------- Game Functions ------------------------------------- \\
func (g *Game) PlacePiece(x, y int, white bool) bool {
	// out of bounds or already occupied spaces are invalid
	if y >= len(g.board) || x >= len(g.board[y]) || g.board[y][x] != nil {
		return false
	}

	// would be eliminated when placed
	if !g.hasRoom(x, y, white, map[Cordinate]bool{}) {
		return false
	}

	// place piece
	g.board[y][x] = &white

	// check if it eliminates any oponents pieces

	return true
}

// -------------------------------------- -------------- ------------------------------------- \\

// --------------------------- Functions required by ebiten engine --------------------------- \\
func (g *Game) Update() error {
	// Write your game's logical update.
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
			x := (piece_y+1)*(g.SquareSize+g.BorderSize) - g.SquareSize/2
			y := (piece_x+1)*(g.SquareSize+g.BorderSize) - g.SquareSize/2
			clr := color.White
			if !*piece {
				clr = color.Black
			}
			ebitenutil.DrawCircle(screen, float64(x), float64(y), 12, clr)
		}

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.Rows*(g.BorderSize+g.SquareSize) + g.BorderSize, g.Columns*(g.BorderSize+g.SquareSize) + g.BorderSize
}

// --------------------------- ------------------------------------ --------------------------- \\
