package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CachedData struct {
	username string
	oauth    string
	channel  string
}

func runIRCGame() {

	// Read cached data
	cachedData := readCachedData()

	app := app.New()

	gameWindow := app.NewWindow("Configure Twitch Sweeps Mines")

	gameWindow.Resize(fyne.NewSize(600, 400))

	configGrid := container.NewGridWithColumns(2)

	usernameLabel := widget.NewLabel("Twitch Username:")
	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(cachedData.username)
	configGrid.Add(usernameLabel)
	configGrid.Add(usernameEntry)

	oauthLabel := widget.NewLabel("Twitch OAuth Token:")
	oauthEntry := widget.NewPasswordEntry()
	oauthEntry.SetText(cachedData.oauth)
	configGrid.Add(oauthLabel)
	configGrid.Add(oauthEntry)

	channelLabel := widget.NewLabel("Twitch Channel:")
	channelEntry := widget.NewEntry()
	channelEntry.SetText(cachedData.channel)
	configGrid.Add(channelLabel)
	configGrid.Add(channelEntry)

	widthLabel := widget.NewLabel("Width:")
	widthEntry := widget.NewEntry()
	widthEntry.SetText(strconv.Itoa(defaultWidth))
	widthEntry.Disable()
	configGrid.Add(widthLabel)
	configGrid.Add(widthEntry)

	heightLabel := widget.NewLabel("Height:")
	heightEntry := widget.NewEntry()
	heightEntry.SetText(strconv.Itoa(defaultHeight))
	heightEntry.Disable()
	configGrid.Add(heightLabel)
	configGrid.Add(heightEntry)

	numMinesLabel := widget.NewLabel("Number of Mines:")
	numMinesEntry := widget.NewEntry()
	numMinesEntry.SetText(strconv.Itoa(defaultMines))
	numMinesEntry.Disable()
	configGrid.Add(numMinesLabel)
	configGrid.Add(numMinesEntry)

	// Button to start the game
	startButton := widget.NewButton("Start", func() {
		username := usernameEntry.Text
		channel := channelEntry.Text
		oauth := oauthEntry.Text

		// username must be lowecase
		username = strings.ToLower(username)

		// Write cached data
		writeCachedData(CachedData{username, oauth, channel})

		width, _ := strconv.Atoi(widthEntry.Text)
		height, _ := strconv.Atoi(heightEntry.Text)
		numMines, _ := strconv.Atoi(numMinesEntry.Text)

		// Connect to the Twitch IRC server
		conn := initTwitchConnection(username, oauth, channel)

		if conn == nil {
			fmt.Println("Failed to connect to Twitch IRC server")
			return
		}

		// Initialize minefield
		minefield := TwitchSweepsMines{}
		minefield.init(width, height, numMines)

		// Start the IRC game
		playIRCGame(conn, minefield, gameWindow)
	})

	exitButton := widget.NewButton("Exit", func() {
		gameWindow.Close()
	})

	configGrid.Add(startButton)
	configGrid.Add(exitButton)

	// Set up the GUI game)
	gameWindow.SetContent(configGrid)

	// Show the window and run the application
	gameWindow.ShowAndRun()
}

func initTwitchConnection(username, oauth, channel string) net.Conn {
	fmt.Println("Initializing Twitch IRC connection: ", channel)

	// Connect to the Twitch IRC server
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		fmt.Println("Failed to connect to Twitch IRC server:", err)
		return nil
	}

	// Send authentication details
	fmt.Fprintf(conn, "PASS oauth:%s\r\n", oauth)
	fmt.Fprintf(conn, "NICK %s\r\n", username)

	// Join the desired Twitch channel
	fmt.Fprintf(conn, "JOIN #%s\r\n", channel)

	fmt.Println("Connected to Twitch IRC server: ", channel)

	return conn
}

func playIRCGame(conn net.Conn, minefield TwitchSweepsMines, window fyne.Window) {

	// Start goroutine for handling Twitch IRC messages
	go func() {
		// Create a reader to read messages from the server
		reader := bufio.NewReader(conn)

		// Start reading messages from the server
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Failed to read message:", err)
				return
			}

			minefield.processMessage(message)

			refreshGameWindow(window, minefield)
		}
	}()

	// Start goroutine for handling game logic
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(minefield.countdown))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("Time is up!")

				// Execute the action with the most votes
				actionExecuted := minefield.executeAction()

				// print the minefield
				minefield.print(false)

				// Refresh the game window
				if actionExecuted {
					refreshGameWindow(window, minefield)

					// Check if the game is over
					if minefield.isGameover {
						if minefield.isWin {
							fmt.Println("Congratulations! You have won the game!")
						} else {
							fmt.Println("Game over! You have hit a mine!")

							// print the minefield
							minefield.print(true)

							// Exit the game
							return
						}
					}
				}
			}
		}
	}()

	// Wait indefinitely or until an exit condition is met
	select {}
}

func readCachedData() CachedData {
	//If directory cache does not exist, create it
	if _, err := os.Stat("cache"); os.IsNotExist(err) {
		os.Mkdir("cache", 0755)
	}

	// if file does not exist, create it
	if _, err := os.Stat("cache/irc"); os.IsNotExist(err) {
		os.Create("cache/irc")
	}

	// read from file
	file, err := os.Open("cache/irc")
	if err != nil {
		return CachedData{}
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	username := scanner.Text()
	scanner.Scan()
	oauth := scanner.Text()
	scanner.Scan()
	channel := scanner.Text()

	return CachedData{username, oauth, channel}
}

// write cached data
func writeCachedData(data CachedData) {
	file, err := os.Create("cache/irc")
	if err != nil {
		return
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(data.username + "\n")
	writer.WriteString(data.oauth + "\n")
	writer.WriteString(data.channel + "\n")
	writer.Flush()
}

func refreshGameWindow(window fyne.Window, minefield TwitchSweepsMines) {
	// Create a new grid layout
	gameGrid := container.NewGridWithColumns(minefield.width + 1)

	mostVotedAction := minefield.getMostVotedAction()

	mostVotedCol := -1
	mostVotedRow := -1
	if len(mostVotedAction) != 0 {
		mostVotedCol = int(mostVotedAction[0] - 'A')
		mostVotedRow, _ = strconv.Atoi(string(mostVotedAction[1]))
	}

	// Add Column Letter
	for x := -1; x < minefield.width; x++ {
		if x >= 0 {
			colLetter := widget.NewLabel(fmt.Sprintf("%c", 'A'+x))
			colLetter.Alignment = fyne.TextAlignCenter
			gameGrid.Add(colLetter)
		} else {
			gameGrid.Add(widget.NewLabel(""))
		}
	}

	// Loop through each cell in the minefield
	for y := 0; y < minefield.height; y++ {

		// Row Numbers
		rowNumber := widget.NewLabel(fmt.Sprint(y))
		gameGrid.Add(rowNumber)

		for x := 0; x < minefield.width; x++ {
			cell := minefield.grid[y][x]

			buttonText := ""
			if cell.isRevealed || minefield.isGameover {

				if minefield.isGameover && cell.isFlagged {
					if cell.hasMine {
						buttonText = "🚩" // Correct flag
					} else {
						buttonText = "❌" // Wrong flag
					}
				} else if cell.hasMine {
					buttonText = "💣"

					// Last action
					lastAction := minefield.history[0]

					if minefield.isGameover && (lastAction.x == x && lastAction.y == y) {
						buttonText = "💥"
					}

				} else if cell.adjacentMines > 0 {
					buttonText = strconv.Itoa(cell.adjacentMines)
				} else {
					buttonText = "."
				}
			} else if cell.isFlagged {
				buttonText = "🚩"
			}

			// Create a new button widget
			button := widget.NewButton(buttonText, func() {

				if minefield.isGameover {
					return
				}

				if cell.isRevealed {
					return
				}

				// handle left click
				if cell.isFlagged {
					return
				}

				minefield.reveal(x, y, true)

				minefield.drawMineField(window)
			})
			if x == mostVotedCol && y == mostVotedRow {
				// TODO highlight the most voted cell
			}
			button.Disable()

			// Add the button to the grid layout
			gameGrid.Add(button)
		}
	}

	//Information Panel
	infoGrid := container.NewGridWithColumns(2)

	// Mines left
	minesLeftLabel := widget.NewLabel("Total Mines:")
	minesLeftValue := widget.NewLabel(fmt.Sprintf("%d", minefield.minesLeft))
	infoGrid.Add(minesLeftLabel)
	infoGrid.Add(minesLeftValue)

	// History
	historyLabel := widget.NewLabel("History:")
	infoGrid.Add(historyLabel)

	historyValues := []string{}
	i := 0
	for i < 5 && i < len(minefield.history) {
		action := minefield.history[i]

		//Convert action.x to column letter
		colLetter := string('A' + action.x)

		historyValues = append(historyValues, fmt.Sprintf("%s%d ", colLetter, action.y))

		i++
	}

	historyValueText := strings.Join(historyValues, "\n")
	historyValue := widget.NewLabel(historyValueText)
	infoGrid.Add(historyValue)

	// infoGrid is right, gameGrid is center
	window.SetContent(container.NewBorder(nil, nil, nil, infoGrid, gameGrid))
}
