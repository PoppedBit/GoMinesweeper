package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultCountdown = 15
)

type TwitchSweepsMines struct {
	Minefield
	users       []string
	actionVotes map[string]int
	paused      bool
	countdown   int
}

func (m *TwitchSweepsMines) init(width, height, numMines int) {
	m.Minefield.init(width, height, numMines)

	m.users = []string{}
	m.actionVotes = make(map[string]int)
	m.paused = false
	m.countdown = defaultCountdown
}

func (minefield *TwitchSweepsMines) processMessage(message string) {

	username, content := parseTwitchMessage(message)

	// Check if the user is already in the game
	for _, u := range minefield.users {
		if u == username {
			return
		}
	}

	// Check message is a valid command
	// command must be in format XY where X is a letter and Y is a number
	if len(content) != 2 {
		return
	}

	// Convert message to coordinates
	message = strings.ToUpper(content)

	col := int(message[0] - 'A')
	row, err := strconv.Atoi(string(message[1]))
	if err != nil {
		return
	}

	// Check if the coordinates are within the minefield
	if col < 0 || col >= minefield.width || row < 0 || row >= minefield.height {
		return
	}

	// Check if the square has already been revealed
	if minefield.grid[row][col].isRevealed {
		return
	}

	// Add the user to the list of users
	minefield.users = append(minefield.users, username)

	// Add the user's action to the actionVotes map
	minefield.actionVotes[message]++

	fmt.Printf("Action recorded - %s: %s\n", username, message)
}

func (minefield *TwitchSweepsMines) printActionsQueue() {
	fmt.Println("Actions queue:")
	for a, v := range minefield.actionVotes {
		fmt.Printf("%s: %d\n", a, v)
	}
}

func (minefield *TwitchSweepsMines) executeAction() {

	minefield.printActionsQueue()

	// Find the most voted action
	maxVotes := 0
	var action string
	for a, v := range minefield.actionVotes {
		if v > maxVotes {
			maxVotes = v
			action = a
		}
	}

	if maxVotes == 0 {
		fmt.Println("No action to execute")
		return
	}

	fmt.Printf("Executing action: %s\n", action)

	// Convert action to coordinates
	col := int(action[0] - 'A')
	row, _ := strconv.Atoi(string(action[1]))

	fmt.Printf("Coordinates: %d, %d\n", col, row)

	// Reveal the square
	minefield.reveal(col, row, true)

	// Reset the action queue
	minefield.resetActionsQueue()
}

func (minefield *TwitchSweepsMines) resetActionsQueue() {
	minefield.users = []string{}
	minefield.actionVotes = make(map[string]int)
	println("Actions queue reset")
}

func parseTwitchMessage(message string) (string, string) {

	// Remove the trailing newline character
	message = strings.TrimSuffix(message, "\r\n")

	// Parse the message to extract the username and the content
	parts := strings.Split(message, " ")
	if len(parts) < 4 {
		return "", ""
	}

	username := strings.Split(parts[0], "!")[0][1:]
	content := strings.Join(parts[3:], " ")[1:]

	return username, content
}
