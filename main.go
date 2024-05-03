package main

import "fmt"

func main() {

	var mode int
	for mode != 1 {
		println("Welcome to Minesweeper")
		println("[1] CLI ")
		println("[2] GUI ")
		print("Choose game mode: ")
		_, _ = fmt.Scanf("%d", &mode)

		if mode != 1 && mode != 2 {
			println("Invalid action, try again")
		}
	}

	switch mode {
	case 1:
		runCLIGame()
	case 2:
		runGUIGame()
	}
}
