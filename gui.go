package main

import (
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

	// Field For User to Enter Width of minefield
	widthEntry := widget.NewEntry()
	widthEntry.SetPlaceHolder("Width")
	widthEntry.SetText(strconv.Itoa(defaultWidth))

	// Field For User to Enter Height of minefield
	heightEntry := widget.NewEntry()
	heightEntry.SetPlaceHolder("Height")
	heightEntry.SetText(strconv.Itoa(defaultHeight))

	// Field For User to Enter Number of Mines
	numMinesEntry := widget.NewEntry()
	numMinesEntry.SetPlaceHolder("Number of Mines")
	numMinesEntry.SetText(strconv.Itoa(defaultMines))

	// Button to start the game
	startButton := widget.NewButton("Start", func() {
		width, _ := strconv.Atoi(widthEntry.Text)
		height, _ := strconv.Atoi(heightEntry.Text)
		numMines, _ := strconv.Atoi(numMinesEntry.Text)

		// Initialize minefield
		m := Minefield{}
		m.init(width, height, numMines)

		m.drawMineField(window)
	})

	// Set up the GUI game
	content := container.NewVBox(widthEntry, heightEntry, numMinesEntry, startButton)
	window.SetContent(content)

	// // Show the window and run the application
	window.ShowAndRun()
}

// Render minefield on GUI
func (m *Minefield) drawMineField(window fyne.Window) {
	// Create a new grid layout
	grid := container.NewGridWithColumns(m.width)

	// Loop through each cell in the minefield
	for i := 0; i < m.height; i++ {
		for j := 0; j < m.width; j++ {
			cell := m.grid[j][i]

			buttonText := ""
			if cell.isRevealed {
				if cell.hasMine {
					buttonText = "X" //"💣"
				} else if cell.adjacentMines > 0 {
					buttonText = strconv.Itoa(cell.adjacentMines)
				} else {
					buttonText = "."
				}
			} else if cell.isFlagged {
				buttonText = "F" // "🚩"
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

				m.reveal(i, j)

				if cell.hasMine {
					// Game over
					m.revealAll()
					m.drawMineField(window)
					return
				}

				if m.allNonMinesRevealed() {
					// Game won
					m.revealAll()
				}

				m.drawMineField(window)
			})

			// Add the button to the grid layout
			grid.Add(button)
		}
	}

	window.SetContent(grid)
}
