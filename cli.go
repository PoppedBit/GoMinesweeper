package main

import (
	"fmt"
)

// size of minefield, width and height
const (
	defaultWidth  = 10
	defaultHeight = 10
)

// run CLI game
func runCLIGame() {

	// Initialize minefield
	m := Minefield{}
	m.init(defaultWidth, defaultHeight, 10)

	// Clear screen
	print("\033[H\033[2J")

	// If mines left or all non-mines revealed and no mine revealed, continue game
	for (!m.allNonMinesRevealed()) && !m.mineRevealed() {
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