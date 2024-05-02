package main

import (
	"math/rand"
)

// Grid Square
type Square struct {
	hasMine         bool
	adjacentMines   int
	contentsVisible bool
	isFlagged       bool
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

// Print Minefield to CLI
func (m *Minefield) print(viewAll bool) {
	for y := range m.grid {
		for x := range m.grid[y] {
			if m.grid[y][x].contentsVisible || viewAll {
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

// Convert Minefield to string
func (m *Minefield) toString(viewAll bool) string {
	var result string
	for y := range m.grid {
		for x := range m.grid[y] {
			if m.grid[y][x].contentsVisible || viewAll {
				if m.grid[y][x].hasMine {
					result += "X"
				} else {
					result += string(m.grid[y][x].adjacentMines + '0')
				}
			} else if m.grid[y][x].isFlagged {
				result += "F"
			} else {
				result += "."
			}
		}
		result += "\n"
	}
	return result
}
