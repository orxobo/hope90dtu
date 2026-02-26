package ui

import (
	"fmt"
	"hope90dtu/device/atcommands"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (w *MainWindow) makeTerminalTab() fyne.CanvasObject {
	w.atResponse = widget.NewMultiLineEntry()
	w.atResponse.SetPlaceHolder("AT response displays here . . .")
	w.atResponse.OnChanged = func(s string) {
		//w.atResponse.Undo()
	}

	// Quick commands grid
	quickCmds := []atcommands.ATCmd{
		atcommands.CmdModel,
		atcommands.CmdVer,
		atcommands.CmdReboot,
		atcommands.CmdSn,
	}

	sendBtn := widget.NewButtonWithIcon("Send", theme.MailSendIcon(), w.handleSendAT)
	sendBtn.Disable()

	cmdGrid := container.NewGridWithColumns(5)
	for _, cmd := range quickCmds {
		atcommand, _ := atcommands.GetCommand(cmd) // cant error
		btn := widget.NewButton(atcommand.Description, func() {
			if sendBtn.Disabled() {
				dialog.ShowInformation("Set Prefix", "AT Client must be initialised before using AT commands", w.window)
				return
			}
			w.SendAT(cmd)
		})
		cmdGrid.Add(btn)
	}

	w.atSelect = widget.NewSelect(atcommands.ATCmds(), nil)
	atcommandContainer := container.NewBorder(nil, nil, nil, sendBtn, w.atSelect)

	atPrefix := widget.NewEntry()
	atPrefix.SetPlaceHolder("Factory default: AT / NETWK")
	w.initATClientBtn = widget.NewButton("Initialise AT Client", func() {
		if atPrefix.Text == "" {
			atPrefix.SetText("AT")
		}
		cli, err := atcommands.NewATClient(w.device, atPrefix.Text)
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		w.atClient = cli
		sendBtn.Enable()
	})
	w.initATClientBtn.Disable()
	prefixContainer := container.NewBorder(nil, nil, nil, w.initATClientBtn, atPrefix)

	atForm := widget.NewForm(
		widget.NewFormItem("AT Prefix", prefixContainer),
		widget.NewFormItem("AT Command", atcommandContainer),
	)

	inputArea := container.NewBorder(atForm, nil, nil, nil, w.atResponse)

	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Quick Commands", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			cmdGrid,
			widget.NewSeparator(),
		),
		nil, nil, nil,
		inputArea,
	)
}

func (w *MainWindow) handleSendAT() {
	cmd := w.atSelect.Selected
	if cmd == "" {
		return
	}
	cmd = strings.TrimSpace(strings.Split(cmd, ":")[0])
	w.SendAT(atcommands.ATCmdFromString(cmd))
}

func (w *MainWindow) SendAT(cmd atcommands.ATCmd) {
	w.appendToATResponsef("Sending: %s", cmd.String())
	resp, err := w.atClient.Run(cmd)
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}
	w.appendToATResponse(fmt.Sprint("Recieved: ", resp))
}

func (w *MainWindow) appendToATResponsef(message string, a ...any) {
	w.appendToATResponse(fmt.Sprintf(message, a...))
}

func (w *MainWindow) appendToATResponse(message string) {
	w.atResponse.Append(fmt.Sprintf("[%s] %s\n", time.Now().Format("15:04:05.000"), message))
	w.appendToMonitor(message)
}
