package game

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/al-pi314/gogo/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Cordinate struct {
	X int
	Y int
}

type Game struct {
	Dymension   int
	SquareSize  int
	BorderSize  int
	WhitePlayer player.Player
	BlackPlayer player.Player

	active bool
	board  [][]*bool

	delay_lock bool
	locked     Cordinate

	did_skip    bool
	whiteToMove bool
	failedMoves int

	black_stones int
	black_taken  int
	white_stones int
	white_taken  int
}

func NewGame(g Game) *Game {
	g.board = make([][]*bool, g.Dymension)
	for y := range g.board {
		g.board[y] = make([]*bool, g.Dymension)
	}
	g.locked = Cordinate{-1, -1}
	g.active = true

	return &g
}

// ------------------------------------ Helper Functions ------------------------------------ \\
func (g *Game) Size() (int, int) {
	side := g.Dymension*(g.SquareSize+g.BorderSize) + g.BorderSize
	return side, side
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
	if chk, ok := checked[c]; ok && chk {
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

func updateCtr(whitePtr, blackPtr *int, iswhite bool, cnt int) {
	if iswhite {
		blackPtr = whitePtr
	}
	*blackPtr += cnt
}

func (g *Game) asignTeritory(x, y int, checked map[Cordinate]bool) (bool, *bool, int) {
	c := Cordinate{x, y}
	if chk, ok := checked[c]; ok && chk {
		return true, nil, 0
	}
	checked[c] = true

	if y < 0 || y >= len(g.board) || x < 0 || x >= len(g.board[y]) {
		return true, nil, 0
	}
	if g.board[y][x] != nil {
		return true, g.board[y][x], 0
	}

	prev_uniform, prev_owner, prev_cnt := g.asignTeritory(x, y-1, checked)
	for _, c := range []Cordinate{{x, y + 1}, {x - 1, y}, {x + 1, y}} {
		curr_uniform, curr_owner, curr_cnt := g.asignTeritory(c.X, c.Y, checked)
		// previous teritory is not uniform - does not belong to just one player
		if !prev_uniform {
			return false, nil, 0
		}

		// teritories ownership missmatch
		if prev_owner != nil && curr_owner != nil && *prev_owner != *curr_owner {
			return false, nil, 0
		}

		prev_uniform = curr_uniform
		if curr_owner != nil {
			prev_owner = curr_owner
		}
		prev_cnt += curr_cnt
	}
	return prev_uniform, prev_owner, 1 + prev_cnt
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
	updateCtr(&g.white_stones, &g.black_stones, white, 1)

	// check for opponent group eliminations
	if g.caputreOpponent(x, y, white) {
		return true
	}

	// would be eliminated when placed
	hasRoom := g.hasRoom(x, y, white, map[Cordinate]bool{}, nil)
	if !hasRoom {
		g.board[y][x] = nil
		updateCtr(&g.white_stones, &g.black_stones, white, -1)
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

	updateCtr(&g.white_stones, &g.black_stones, !white, -len(toRemove))
	updateCtr(&g.white_taken, &g.black_taken, !white, len(toRemove))
	return len(toRemove) != 0
}

// Score calculates game score based on the current position.
func (g *Game) Score() float64 {
	checked := map[Cordinate]bool{}
	score := -0.5 + float64(g.white_stones) - float64(g.white_taken) - float64(g.black_stones) + float64(g.black_taken)
	for y := range g.board {
		for x := range g.board[y] {
			if chk, ok := checked[Cordinate{x, y}]; g.board[y][x] != nil || (ok && chk) {
				continue
			}
			uniform, owner, size := g.asignTeritory(x, y, checked)
			if uniform && owner != nil {
				sign := 1
				if !(*owner) {
					sign = -1
				}
				score += float64(sign * size)
			}
		}
	}
	return score
}

// -------------------------------------- -------------- ------------------------------------- \\

// --------------------------- Functions required by ebiten engine --------------------------- \\
func (g *Game) Update() error {
	if !g.active {
		return errors.New("game finished")
	}

	player := g.WhitePlayer
	if !g.whiteToMove {
		player = g.BlackPlayer
	}

	skip, x, y := player.Place(g.board)
	skip = skip || g.failedMoves == 3
	if skip || ((x != nil && y != nil) && g.placePiece(*x, *y, g.whiteToMove)) {
		// consequitive skips end the game
		if skip && g.did_skip {
			g.active = false
		}
		g.did_skip = false
		if skip {
			g.did_skip = true
		}

		// lock unlocks after next successful move
		if !g.delay_lock {
			g.locked.X = -1
			g.locked.Y = -1
		}
		g.delay_lock = false
		// change player to move
		g.failedMoves = 0
		g.whiteToMove = !g.whiteToMove
	} else {
		g.failedMoves++
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// end screen
	if !g.active {
		screen.Fill(color.White)
		score := g.Score()
		winner := "Black"
		if score >= -0.5 {
			winner = "White"
		}
		face := basicfont.Face7x13
		txt := fmt.Sprintf("%s player won! Score: %.2f", winner, score)
		centerX := 0.5*float64(g.Dymension*(g.SquareSize+g.BorderSize)+g.BorderSize) - float64(face.Width*len(txt))/2
		centerY := 0.5 * float64(g.Dymension*(g.SquareSize+g.BorderSize)+g.BorderSize)
		text.Draw(screen, txt, face, int(centerX), int(centerY), color.Black)
		return
	}

	// draw board - squares with left and top borders
	for x := 0; x <= g.Dymension+1; x++ {
		for y := 0; y <= g.Dymension+1; y++ {
			x1 := x * (g.BorderSize + g.SquareSize)
			y1 := y * (g.BorderSize + g.SquareSize)
			// draw left border
			if y <= g.Dymension {
				x2 := x1 + g.BorderSize
				y2 := y1 + g.SquareSize + g.BorderSize
				g.drawSquare(screen, x1, y1, x2, y2, color.RGBA{160, 175, 190, 1})
			}
			// draw top border
			if x <= g.Dymension {
				x2 := x1 + g.SquareSize + g.BorderSize
				y2 := y1 + g.BorderSize
				g.drawSquare(screen, x1, y1, x2, y2, color.RGBA{160, 175, 190, 1})
			}

			// draw empty square when inside the board
			if x <= g.Dymension && y <= g.Dymension {
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
	return g.Size()
}

// --------------------------- ------------------------------------ --------------------------- \\
