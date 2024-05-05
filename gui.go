package main

import (
	"fmt"
	"strconv"

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
	grid := container.NewGridWithColumns(m.width + 1)

	// Add Column Letter
	for x := -1; x < m.width; x++ {
		if x >= 0 {
			colLetter := widget.NewLabel(fmt.Sprintf("%c", 'A'+x))
			grid.Add(colLetter)
		} else {
			grid.Add(widget.NewLabel(""))
		}
	}

	// Loop through each cell in the minefield
	for y := 0; y < m.height; y++ {

		// Row Numbers
		rowNumber := widget.NewLabel(fmt.Sprint(y))
		grid.Add(rowNumber)

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
					// TODO fix this
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

				m.reveal(x, y)

				m.drawMineField(window)
			})

			// Add the button to the grid layout
			grid.Add(button)
		}
	}

	window.SetContent(grid)
}
