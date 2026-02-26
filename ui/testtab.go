package ui

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (w *MainWindow) makeTestTab() fyne.CanvasObject {
	// Random data test
	lengthEntry := widget.NewEntry()
	lengthEntry.SetText("10")

	repeatsEntry := widget.NewEntry()
	repeatsEntry.SetText("10")

	intervalEntry := widget.NewEntry()
	intervalEntry.SetText("1000")

	w.testBtn = widget.NewButtonWithIcon("Start Random Test", theme.MediaPlayIcon(), func() {
		w.handleRandomTest(lengthEntry.Text, intervalEntry.Text, repeatsEntry.Text)
	})
	w.testBtn.Disable()

	randomForm := widget.NewForm(
		widget.NewFormItem("Length (bytes)",
			container.NewHBox(
				lengthEntry,
				layout.NewSpacer(),
				widget.NewLabelWithStyle("# of Repeats", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				repeatsEntry,
			),
		),
		widget.NewFormItem("Interval (ms)", container.NewBorder(nil, nil, nil, w.testBtn, intervalEntry)),
	)

	// RSSI
	rssiBackEntry := widget.NewEntry()
	rssiLastEntry := widget.NewEntry()

	rssiForm := widget.NewForm(
		widget.NewFormItem("RSSI Background Noise",
			container.NewBorder(
				nil, nil, nil,
				widget.NewButton("Get RSSI Background Noise", func() {
					rssiBackEntry.SetText(w.device.GetBackgroundNoise())
				}),
				rssiBackEntry,
			),
		),
		widget.NewFormItem("RSSI Last Responce Noise",
			container.NewBorder(
				nil, nil, nil,
				widget.NewButton("Get RSSI Last Response", func() {
					rssiLastEntry.SetText(w.device.GetLastResponseNoise())
				}),
				rssiLastEntry,
			),
		),
	)

	// Custom HEX
	hexData := widget.NewEntry()
	hexData.SetPlaceHolder("55 AA 00 02 DE AD BE EF")

	stringData := widget.NewEntry()
	stringData.SetPlaceHolder("TEST STRING")

	hexForm := widget.NewForm(
		widget.NewFormItem("Hex Data",
			container.NewBorder(
				nil, nil, nil,
				widget.NewButton("Send Hex", func() {
					w.handleSendHex(hexData.Text)
				}),
				hexData,
			),
		),
		widget.NewFormItem("String Data",
			container.NewBorder(
				nil, nil, nil,
				widget.NewButton("Send String", func() {
					w.handleSendString(stringData.Text)
				}),
				stringData,
			),
		),
	)

	return container.NewVBox(
		widget.NewCard("Random Generator", "Send random bytes to test response.", randomForm),
		widget.NewSeparator(),
		widget.NewCard("RSSI", "Get signal noise in dBm on current channel and the noise from the last response.", rssiForm),
		widget.NewSeparator(),
		widget.NewCard("Custom Data", "Send specified data.", hexForm),
	)
}

func (w *MainWindow) handleRandomTest(lengthStr, intervalStr, repeatsStr string) {
	length, _ := strconv.Atoi(lengthStr)
	interval, _ := strconv.Atoi(intervalStr)
	repeats, _ := strconv.Atoi(repeatsStr)

	if length <= 0 || interval <= 0 || repeats <= 0 {
		return
	}

	go func() {
		for i := 0; i < repeats && w.device != nil; i++ {
			w.appendToMonitorf("Random Data test, Sent %d random bytes", length)
			len, err := w.device.SendRandomData(length)
			if err != nil {
				w.appendToMonitorf("Random data test failure: %v", err)
			}
			w.appendToMonitorf("Random Data test, Received %d random bytes", len)
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()
}

func (w *MainWindow) handleSendHex(hexStr string) {
	err := w.device.SendUDPHexCommand(hexStr)
	if err != nil {
		dialog.ShowError(err, w.window)
	}
}

func (w *MainWindow) handleSendString(data string) {
	err := w.device.SendUDPASCIICommand(data)
	if err != nil {
		dialog.ShowError(err, w.window)
	}
}
