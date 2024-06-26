package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Run GUI game
func runGUIGame() {
	// Create a new Fyne application
	app := app.New()

	// // Create a new window
	window := app.NewWindow("Minesweeper GUI")
	window.Resize(fyne.NewSize(600, 400))

	widthLabel := widget.NewLabel("Width:")
	widthEntry := widget.NewEntry()
	widthEntry.SetPlaceHolder("Width")
	widthEntry.SetText(strconv.Itoa(defaultWidth))
	widthField := container.NewHBox(widthLabel, widthEntry)

	// Field For User to Enter Height of minefield
	heightLabel := widget.NewLabel("Height:")
	heightEntry := widget.NewEntry()
	heightEntry.SetPlaceHolder("Height")
	heightEntry.SetText(strconv.Itoa(defaultHeight))
	heightField := container.NewHBox(heightLabel, heightEntry)

	// Field For User to Enter Number of Mines
	numMinesLabel := widget.NewLabel("Number of Mines:")
	numMinesEntry := widget.NewEntry()
	numMinesEntry.SetPlaceHolder("Number of Mines")
	numMinesEntry.SetText(strconv.Itoa(defaultMines))
	numMinesField := container.NewHBox(numMinesLabel, numMinesEntry)

	// Button to start the game
	startButton := widget.NewButton("Start", func() {
		width, _ := strconv.Atoi(widthEntry.Text)
		height, _ := strconv.Atoi(heightEntry.Text)
		numMines, _ := strconv.Atoi(numMinesEntry.Text)

		// Initialize minefield
		minefield := Minefield{}
		minefield.init(width, height, numMines)

		minefield.drawMineField(window)
	})

	// Set up the GUI game
	content := container.NewVBox(widthField, heightField, numMinesField, startButton)
	window.SetContent(content)

	// // Show the window and run the application
	window.ShowAndRun()
}

// Render minefield on GUI
func (m *Minefield) drawMineField(window fyne.Window) {
	// Create a new grid layout
	gameGrid := container.NewGridWithColumns(m.width + 1)

	// Add Column Letter
	for x := -1; x < m.width; x++ {
		if x >= 0 {
			colLetter := widget.NewLabel(fmt.Sprintf("%c", 'A'+x))
			colLetter.Alignment = fyne.TextAlignCenter
			gameGrid.Add(colLetter)
		} else {
			gameGrid.Add(widget.NewLabel(""))
		}
	}

	// Loop through each cell in the minefield
	for y := 0; y < m.height; y++ {

		// Row Numbers
		rowNumber := widget.NewLabel(fmt.Sprint(y))
		gameGrid.Add(rowNumber)

		for x := 0; x < m.width; x++ {
			cell := m.grid[y][x]

			buttonText := ""
			if cell.isRevealed || m.isGameover {

				if m.isGameover && cell.isFlagged {
					if cell.hasMine {
						buttonText = "🚩" // Correct flag
					} else {
						buttonText = "❌" // Wrong flag
					}
				} else if cell.hasMine {
					buttonText = "💣"

					// Last action
					lastAction := m.history[0]

					if m.isGameover && (lastAction.x == x && lastAction.y == y) {
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

				if m.isGameover {
					return
				}

				if cell.isRevealed {
					return
				}

				// handle left click
				if cell.isFlagged {
					return
				}

				m.reveal(x, y, true)

				m.drawMineField(window)
			})

			// Add the button to the grid layout
			gameGrid.Add(button)
		}
	}

	//Information Panel
	infoGrid := container.NewGridWithColumns(2)

	// Mines left
	minesLeftLabel := widget.NewLabel("Mines Left:")
	minesLeftValue := widget.NewLabel(fmt.Sprintf("%d", m.minesLeft))
	infoGrid.Add(minesLeftLabel)
	infoGrid.Add(minesLeftValue)

	// History
	historyLabel := widget.NewLabel("History:")
	infoGrid.Add(historyLabel)

	historyValues := []string{}
	i := 0
	for i < 5 && i < len(m.history) {
		action := m.history[i]

		//Convert action.x to column letter
		colLetter := string('A' + action.x)

		historyValues = append(historyValues, fmt.Sprintf("%s%d ", colLetter, action.y))

		i++
	}

	historyValueText := strings.Join(historyValues, "\n")
	historyValue := widget.NewLabel(historyValueText)
	infoGrid.Add(historyValue)

	buttonGrid := container.NewGridWithColumns(2)

	restartButton := widget.NewButton("Play Again?", func() {
		m.init(m.width, m.height, m.mines)
		m.drawMineField(window)
	})
	if !m.isGameover {
		restartButton.Disable()
	}
	buttonGrid.Add(restartButton)

	exitButton := widget.NewButton("Exit", func() {
		window.Close()
	})
	buttonGrid.Add(exitButton)

	// Add grids to the window
	// infoGrid is right, gameGrid is center
	window.SetContent(container.NewBorder(nil, buttonGrid, nil, infoGrid, gameGrid))
}
