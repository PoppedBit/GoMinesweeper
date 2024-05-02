package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// size of minefield, width and height
const (
	width  = 10
	height = 10
)

func main() {
	println("Minefield")

	// Initialize minefield
	m := Minefield{}
	m.init(width, height, 10)
	println("Revealed View")
	m.print(true)

	// Create a new Fyne application
	app := app.New()

	// Create a widget to display the minefield
	minefieldWidget := widget.NewLabel("Minefield")

	// Create a container to hold the minefield widget
	content := container.NewVBox(minefieldWidget)

	// Set the minefield widget's text to the minefield's contents
	minefieldWidget.SetText(m.toString(false))

	// Create a new window with the minefield container as its content
	window := app.NewWindow("Minefield")
	window.Resize(fyne.NewSize(600, 400))

	window.SetContent(content)

	// Show the window and run the application
	window.ShowAndRun()
}
