package main

import (
	"hope90dtu/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	e90App := app.New()
	mainWindow := ui.NewMainWindow(e90App)
	mainWindow.Show()
	e90App.Run()
}
