package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	// Screen constants
	screenWidth  = 300
	screenHeight = 300
	lineWidth    = 5
	boardRows    = 3
	boardCols    = 3
	squareSize   = screenWidth / boardCols

	// Drawing constants
	circleRadius = squareSize / 3
	circleWidth  = 15
	crossWidth   = 25
)

var (
	// Colors
	white = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	gray  = color.RGBA{R: 180, G: 180, B: 180, A: 255}
	red   = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	green = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	black = color.RGBA{R: 0, G: 0, B: 0, A: 255}
)

type Game struct {
	board    [boardRows][boardCols]int
	player   int
	gameOver bool
}

func NewGame() *Game {
	return &Game{
		player: 1,
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.restartGame()
		}
		return nil
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		clickedRow := y / squareSize
		clickedCol := x / squareSize

		if clickedRow >= 0 && clickedRow < boardRows && clickedCol >= 0 && clickedCol < boardCols {
			if g.isSquareEmpty(clickedRow, clickedCol) {
				g.markSquare(clickedRow, clickedCol, g.player)

				if g.checkWin(g.player) {
					g.gameOver = true
				}
				g.player = 3 - g.player // Switch between 1 and 2

				if !g.gameOver {
					if g.bestMove() {
						if g.checkWin(2) {
							g.gameOver = true
						}
						g.player = 3 - g.player
					}
				}

				if !g.gameOver && g.isBoardFull() {
					g.gameOver = true
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(black)

	// Draw grid lines
	drawColor := white
	if g.gameOver {
		if g.checkWin(1) {
			drawColor = green
		} else if g.checkWin(2) {
			drawColor = red
		} else {
			drawColor = gray
		}
	}

	// Draw vertical lines
	for i := 1; i < boardCols; i++ {
		vector.StrokeLine(screen,
			float32(i*squareSize), 0,
			float32(i*squareSize), float32(screenHeight),
			1, drawColor, false)
	}

	// Draw horizontal lines
	for i := 1; i < boardRows; i++ {
		vector.StrokeLine(screen,
			0, float32(i*squareSize),
			float32(screenWidth), float32(i*squareSize),
			1, drawColor, false)
	}

	// Draw figures (X's and O's)
	for row := 0; row < boardRows; row++ {
		for col := 0; col < boardCols; col++ {
			centerX := float64(col*squareSize + squareSize/2)
			centerY := float64(row*squareSize + squareSize/2)

			switch g.board[row][col] {
			case 1: // Circle
				drawCircle(screen, centerX, centerY, circleRadius, drawColor)
			case 2: // Cross
				offset := float32(squareSize) / 3
				vector.StrokeLine(screen,
					float32(centerX)-offset, float32(centerY)-offset,
					float32(centerX)+offset, float32(centerY)+offset,
					1, drawColor, false)
				vector.StrokeLine(screen,
					float32(centerX)-offset, float32(centerY)+offset,
					float32(centerX)+offset, float32(centerY)-offset,
					1, drawColor, false)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) markSquare(row, col, player int) {
	g.board[row][col] = player
}

func (g *Game) isSquareEmpty(row, col int) bool {
	return g.board[row][col] == 0
}

func (g *Game) isBoardFull() bool {
	for row := 0; row < boardRows; row++ {
		for col := 0; col < boardCols; col++ {
			if g.board[row][col] == 0 {
				return false
			}
		}
	}
	return true
}

func (g *Game) checkWin(player int) bool {
	// Check rows
	for row := 0; row < boardRows; row++ {
		if g.board[row][0] == player && g.board[row][1] == player && g.board[row][2] == player {
			return true
		}
	}

	// Check columns
	for col := 0; col < boardCols; col++ {
		if g.board[0][col] == player && g.board[1][col] == player && g.board[2][col] == player {
			return true
		}
	}

	// Check diagonals
	if g.board[0][0] == player && g.board[1][1] == player && g.board[2][2] == player {
		return true
	}
	if g.board[0][2] == player && g.board[1][1] == player && g.board[2][0] == player {
		return true
	}

	return false
}

func (g *Game) minimax(depth int, isMaximizing bool) float64 {
	if g.checkWin(2) {
		return math.Inf(1)
	}
	if g.checkWin(1) {
		return math.Inf(-1)
	}
	if g.isBoardFull() {
		return 0
	}

	if isMaximizing {
		bestScore := math.Inf(-1)
		for row := 0; row < boardRows; row++ {
			for col := 0; col < boardCols; col++ {
				if g.isSquareEmpty(row, col) {
					g.board[row][col] = 2
					score := g.minimax(depth+1, false)
					g.board[row][col] = 0
					bestScore = math.Max(score, bestScore)
				}
			}
		}
		return bestScore
	}

	bestScore := math.Inf(1)
	for row := 0; row < boardRows; row++ {
		for col := 0; col < boardCols; col++ {
			if g.isSquareEmpty(row, col) {
				g.board[row][col] = 1
				score := g.minimax(depth+1, true)
				g.board[row][col] = 0
				bestScore = math.Min(score, bestScore)
			}
		}
	}
	return bestScore
}

func (g *Game) bestMove() bool {
	bestScore := math.Inf(-1)
	bestMove := struct{ row, col int }{-1, -1}

	for row := 0; row < boardRows; row++ {
		for col := 0; col < boardCols; col++ {
			if g.isSquareEmpty(row, col) {
				g.board[row][col] = 2
				score := g.minimax(0, false)
				g.board[row][col] = 0
				if score > bestScore {
					bestScore = score
					bestMove.row = row
					bestMove.col = col
				}
			}
		}
	}

	if bestMove.row != -1 && bestMove.col != -1 {
		g.markSquare(bestMove.row, bestMove.col, 2)
		return true
	}
	return false
}

func (g *Game) restartGame() {
	g.board = [boardRows][boardCols]int{}
	g.gameOver = false
	g.player = 1
}

// Helper function to draw a circle
func drawCircle(screen *ebiten.Image, x, y, radius float64, clr color.Color) {
	steps := 32
	for i := 0; i < steps; i++ {
		angle1 := float64(i) * 2 * math.Pi / float64(steps)
		angle2 := float64(i+1) * 2 * math.Pi / float64(steps)

		x1 := x + radius*math.Cos(angle1)
		y1 := y + radius*math.Sin(angle1)
		x2 := x + radius*math.Cos(angle2)
		y2 := y + radius*math.Sin(angle2)

		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 1, clr, false)
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tic Tac Toe with MiniMax")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
