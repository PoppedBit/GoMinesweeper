package main

import (
	"fmt"
	"math/rand"
)

// Grid Square
type Square struct {
	hasMine       bool
	adjacentMines int
	revealed      bool
	isFlagged     bool
	color         string
}

type Guess struct {
	action string
	x, y   int
}

// Minefield
type Minefield struct {
	grid      [][]Square
	minesLeft int
	history   []Guess
}

// size of minefield, width and height
const (
	width  = 10
	height = 10
)

// run CLI game
func runCLIGame() {

	// Initialize minefield
	m := Minefield{}
	m.init(width, height, 10)

	// Clear screen
	print("\033[H\033[2J")

	// If mines left or all non-mines revealed and no mine revealed, continue game
	for (!m.allNonMinesRevealed()) && !m.mineRevealed() {
		m.print(false)

		var action string
		for action != "r" && action != "f" {
			print("Enter action [r/f]: ")
			_, _ = fmt.Scanf("%s %c ", &action)

			if action != "r" && action != "f" {
				println("Invalid action, try again")
			}
		}

		col := -1
		row := -1
		for col < 0 || col >= width || row < 0 || row >= height {
			print("Enter coordinates [x y]: ")
			_, _ = fmt.Scanf("%c %d", &col, &row)

			print("col: ", col, " row: ", row, "\n")

			if col >= 'a' && col <= 'z' {
				col -= 'a' - 'A'
			}

			col -= 'A'

			if col < 0 || col >= width || row < 0 || row >= height {
				println("Invalid coordinates, try again")
			}
		}

		switch action {
		case "r":
			m.reveal(col, row)
		case "f":
			m.flag(col, row)
		}

		// Clear screen
		print("\033[H\033[2J")
	}

	println("Mines left:", m.minesLeft)
	m.print(true)

	if m.mineRevealed() {
		println("You lose!")
	} else {
		println("You win!")
	}
}

// Initialize minefield based on width, height, and number of mines
func (m *Minefield) init(width, height, numMines int) {
	m.grid = make([][]Square, height)
	m.minesLeft = numMines

	for i := range m.grid {
		m.grid[i] = make([]Square, width)
	}

	// Randomly distribute mines
	for numMines > 0 {
		x := rand.Intn(width)
		y := rand.Intn(height)
		if !m.grid[y][x].hasMine {
			m.grid[y][x].hasMine = true
			numMines--

			// Increment adjacent mines for surrounding squares
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if x+i >= 0 && x+i < width && y+j >= 0 && y+j < height {
						m.grid[y+j][x+i].adjacentMines++
					}
				}
			}

		}
	}

	// Checkerboard pattern
	for i := range m.grid {
		for j := range m.grid[i] {
			if (i+j)%2 == 0 {
				m.grid[i][j].color = "black"
			} else {
				m.grid[i][j].color = "white"
			}
		}
	}
}

// Print Minefield to CLI, with column letters and row numbers
func (m *Minefield) print(viewAll bool) {

	// white := color.New(color.BgWhite).SprintFunc() // Fg* was the text color
	// black := color.New(color.BgHiMagenta).SprintFunc()

	// Print column letters
	print(" ")
	for i := 0; i < len(m.grid[0]); i++ {
		print(" ")
		print(string('A' + i))
	}
	println()

	for y := range m.grid {

		// Print row numbers
		print(y)

		for x := range m.grid[y] {
			print(" ")

			if m.grid[y][x].revealed || viewAll {
				if m.grid[y][x].hasMine {
					print("X")
				} else {
					print(m.grid[y][x].adjacentMines)
				}
			} else if m.grid[y][x].isFlagged {
				print("F")
			} else {
				if m.grid[y][x].color == "white" {
					// print(white(" "))
					print(".")
				} else {
					// print(black(" "))
					print(".")
				}
			}
		}

		print(" ")
		if y == 0 {
			print("Mines left: ", m.minesLeft)
		}

		println()
	}
}

// Flag square at x, y
func (m *Minefield) flag(x, y int) {
	if x < 0 || x >= len(m.grid[0]) || y < 0 || y >= len(m.grid) {
		return
	}

	m.grid[y][x].isFlagged = !m.grid[y][x].isFlagged

	if m.grid[y][x].isFlagged {
		m.minesLeft--
	} else {
		m.minesLeft++
	}

	guess := Guess{"f", x, y}
	m.history = append(m.history, guess)
}

// Reveal square at x, y
func (m *Minefield) reveal(x, y int) {
	if x < 0 || x >= len(m.grid[0]) || y < 0 || y >= len(m.grid) {
		return
	}

	if m.grid[y][x].revealed {
		return
	}

	m.grid[y][x].revealed = true
	if m.grid[y][x].hasMine {
		m.minesLeft--
	}

	if m.grid[y][x].adjacentMines == 0 {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				m.reveal(x+i, y+j)
			}
		}
	}

	guess := Guess{"r", x, y}
	m.history = append(m.history, guess)
}

// Check if mine has been revealed
func (m *Minefield) mineRevealed() bool {
	for y := range m.grid {
		for x := range m.grid[y] {
			if m.grid[y][x].hasMine && m.grid[y][x].revealed {
				return true
			}
		}
	}
	return false
}

// All non-mines revealed
func (m *Minefield) allNonMinesRevealed() bool {
	for y := range m.grid {
		for x := range m.grid[y] {
			if !m.grid[y][x].hasMine && !m.grid[y][x].revealed {
				return false
			}
		}
	}
	return true
}
