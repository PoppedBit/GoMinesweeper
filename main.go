package main

import (
	"fmt"
	"os"
	"strconv"
)

// mode can be passed in a a command line argument
func main() {
	var mode int

	// Check if a command line argument for mode is provided
	if len(os.Args) > 1 {
		arg, err := strconv.Atoi(os.Args[1])
		if err == nil {
			mode = arg
		}
	}

	for mode != 1 && mode != 2 {
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
