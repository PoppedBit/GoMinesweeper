package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

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

func runIRCGame() {

	// Read cached data
	cachedData := readCachedData()

	app := app.New()

	configWindow := app.NewWindow("Configure Twitch Sweeps Mines")

	configWindow.Resize(fyne.NewSize(600, 400))

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
			// TODO error state
			return
		}

		// Initialize minefield
		minefield := TwitchSweepsMines{}
		minefield.init(width, height, numMines)

		// Start the IRC game
		playIRCGame(conn, minefield)
	})

	exitButton := widget.NewButton("Exit", func() {
		configWindow.Close()
	})

	configGrid.Add(startButton)
	configGrid.Add(exitButton)

	// Set up the GUI game)
	configWindow.SetContent(configGrid)

	// Show the window and run the application
	configWindow.ShowAndRun()
}

func initTwitchConnection(username, oauth, channel string) net.Conn {
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

	return conn
}

func playIRCGame(conn net.Conn, minefield TwitchSweepsMines) {

	// Create a reader to read messages from the server
	reader := bufio.NewReader(conn)

	// Start reading messages from the server
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read message:", err)
			return
		}

		username, content := parseTwitchMessage(message)

		minefield.processMessage(username, content)
	}
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

func (minefield *TwitchSweepsMines) resetActionsQueue() {
	minefield.users = []string{}
	minefield.actionVotes = make(map[string]int)
}
