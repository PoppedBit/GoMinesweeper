package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Run GUI game
func runGUIGame() {
	// Create a new Fyne application
	app := app.New()

	// Create a widget to display the minefield
	minefieldWidget := widget.NewLabel("Minefield")

	// Create a container to hold the minefield widget
	content := container.NewVBox(minefieldWidget)

	// // Set the minefield widget's text to the minefield's contents
	minefieldWidget.SetText("")

	// // Create a new window with the minefield container as its content
	window := app.NewWindow("Minefield")
	window.Resize(fyne.NewSize(600, 400))

	window.SetContent(content)

	// // Show the window and run the application
	window.ShowAndRun()
}
