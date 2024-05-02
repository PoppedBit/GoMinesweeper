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
}

// Minefield
type Minefield struct {
	grid      [][]Square
	minesLeft int
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
	for (m.minesLeft > 0 && !m.allNonMinesRevealed()) && !m.mineRevealed() {
		println("Mines left:", m.minesLeft)
		m.print(false)

		var col, row int
		var action string

		println("Enter action (r/f x y):")
		_, _ = fmt.Scanf("%s %c %d", &action, &col, &row)

		// if col is lowercase, convert to uppercase
		if col >= 'a' && col <= 'z' {
			col -= 'a' - 'A'
		}

		// Convert x to 0-indexed integer
		col -= 'A'

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
}

// Print Minefield to CLI, with column letters and row numbers
func (m *Minefield) print(viewAll bool) {

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
				print(".")
			}
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
