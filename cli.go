package main

import (
	"fmt"

	"github.com/fatih/color"
)

// run CLI game
func runCLIGame() {

	// Clear screen
	print("\033[H\033[2J")

	var width, height, numMines int
	for width < 1 {
		print("Enter width: ")
		_, _ = fmt.Scanf("%d", &width)

		if width < 1 {
			println("Invalid width, try again")
		}
	}
	for height < 1 {
		print("Enter height: ")
		_, _ = fmt.Scanf("%d", &height)

		if height < 1 {
			println("Invalid height, try again")
		}
	}

	numMines = width * height / 10

	// Initialize minefield
	m := Minefield{}
	m.init(width, height, numMines)

	// Clear screen
	print("\033[H\033[2J")

	// If mines left or all non-mines revealed and no mine revealed, continue game
	for !m.isGameover {
		m.print(false)

		var action string
		for action != "r" && action != "f" {
			print("Enter action [r/f]: ")
			_, _ = fmt.Scanf("%s", &action)

			if action != "r" && action != "f" {
				println("Invalid action, try again")
			}
		}

		col := -1
		row := -1
		for col < 0 || col >= m.width || row < 0 || row >= m.height {
			print("Enter coordinates [x y]: ")
			_, _ = fmt.Scanf("%c %d", &col, &row)

			if col >= 'a' && col <= 'z' {
				col -= 'a' - 'A'
			}

			col -= 'A'

			if col < 0 || col >= m.width || row < 0 || row >= m.height {
				println("Invalid coordinates, try again")
			}
		}

		switch action {
		case "r":
			m.reveal(col, row, true)
		case "f":
			m.flag(col, row)
		}

		// Clear screen
		print("\033[H\033[2J")
	}

	println("Mines left:", m.minesLeft)
	m.print(true)

	if !m.isWin {
		println("You lose!")
	} else {
		println("You win!")
	}

	var playAgain string
	for playAgain != "y" && playAgain != "n" {
		print("Play again? [y/n]: ")
		_, _ = fmt.Scanf("%s", &playAgain)

		// Make playAgain lowercase
		if playAgain == "Y" || playAgain == "N" {
			playAgain = string(playAgain[0] + 32)
		}

		if playAgain != "y" && playAgain != "n" {
			println("Invalid input, try again")
		}
	}

	if playAgain == "y" {
		runCLIGame()
	}
}

// Print Minefield to CLI, with column letters and row numbers
func (m *Minefield) print(viewAll bool) {

	green := color.New(color.FgGreen).SprintFunc()

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

			if m.grid[y][x].isRevealed || viewAll {
				if m.grid[y][x].hasMine {
					print("X")
				} else {
					print(m.grid[y][x].adjacentMines)
				}
			} else if m.grid[y][x].isFlagged {
				print(green("F"))
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
