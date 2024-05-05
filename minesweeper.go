package main

import (
	"math/rand"
)

// size of minefield, width and height
const (
	defaultWidth  = 10
	defaultHeight = 10
	defaultMines  = 10
)

// Grid Square
type Square struct {
	hasMine       bool
	adjacentMines int
	isRevealed    bool
	isFlagged     bool
	color         string
}

type Guess struct {
	action string
	x, y   int
}

// Minefield
type Minefield struct {
	grid       [][]Square
	height     int
	width      int
	minesLeft  int
	history    []Guess
	isGameover bool
	isWin      bool
}

// Initialize minefield based on width, height, and number of mines
func (m *Minefield) init(width, height, numMines int) {
	m.grid = make([][]Square, height)
	m.width = width
	m.height = height
	m.minesLeft = numMines
	m.isGameover = false
	m.isWin = false

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
	m.history = append([]Guess{guess}, m.history...)
}

// Reveal square at x, y
func (m *Minefield) reveal(x, y int) {
	if x < 0 || x >= len(m.grid[0]) || y < 0 || y >= len(m.grid) {
		return
	}

	if m.grid[y][x].isRevealed {
		return
	}

	m.grid[y][x].isRevealed = true
	if m.grid[y][x].hasMine {
		m.isGameover = true
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

// All non-mines revealed
func (m *Minefield) allNonMinesRevealed() bool {
	for y := range m.grid {
		for x := range m.grid[y] {
			if !m.grid[y][x].hasMine && !m.grid[y][x].isRevealed {
				return false
			}
		}
	}
	m.isGameover = true
	m.isWin = true
	return true
}
