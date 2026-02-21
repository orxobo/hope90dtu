package main

import (
	"e90dtu/ui"

	"fyne.io/fyne/v2/app"
)

func main() {

	myApp := app.NewWithID("com.icsoft.e90config")
	mainWindow := ui.NewMainWindow(myApp)

	mainWindow.Show()
	myApp.Run()
}
