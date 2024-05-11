package main

import (
	"strconv"
	"strings"
)

type TwitchSweepsMines struct {
	Minefield
	users       []string
	actionVotes map[string]int
}

func (m *TwitchSweepsMines) init(width, height, numMines int) {
	m.Minefield.init(width, height, numMines)

	m.users = []string{}
	m.actionVotes = make(map[string]int)
}

func (minefield *TwitchSweepsMines) processMessage(username, message string) {

	// Check if the user is already in the game
	for _, u := range minefield.users {
		if u == username {
			return
		}
	}

	// Check message is a valid command
	// command must be in format XY where X is a letter and Y is a number
	if len(message) != 2 {
		return
	}

	// Convert message to coordinates
	message = strings.ToUpper(message)
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
}

func (minefield *TwitchSweepsMines) executeAction() {

	// Find the most voted action
	maxVotes := 0
	var action string
	for a, v := range minefield.actionVotes {
		if v > maxVotes {
			maxVotes = v
			action = a
		}
	}

	// Convert action to coordinates
	col := int(action[0] - 'A')
	row, _ := strconv.Atoi(string(action[1]))

	// Reveal the square
	minefield.reveal(col, row, true)

	// Reset the action queue
	minefield.resetActionsQueue()
}

func (minefield *TwitchSweepsMines) resetActionsQueue() {
	minefield.users = []string{}
	minefield.actionVotes = make(map[string]int)
}
