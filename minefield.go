package main

import (
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
	m.minesLeft--

	if m.grid[y][x].adjacentMines == 0 {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				m.reveal(x+i, y+j)
			}
		}
	}
}

// Reveal all mines
func (m *Minefield) revealAll() {
	for y := range m.grid {
		for x := range m.grid[y] {
			m.grid[y][x].revealed = true
		}
	}
}

// Check if all non-mine squares have been revealed
func (m *Minefield) allRevealed() bool {
	for y := range m.grid {
		for x := range m.grid[y] {
			if !m.grid[y][x].hasMine && !m.grid[y][x].revealed {
				return false
			}
		}
	}
	return true
}
