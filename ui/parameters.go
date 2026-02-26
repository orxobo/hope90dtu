// Tab to update device parameters
package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (w *MainWindow) makeParametersTab() fyne.CanvasObject {

	// Lora Parameters
	w.wirelessForm = widget.NewForm(
		widget.NewFormItem("Module Address", widget.NewEntry()),
		widget.NewFormItem("Network ID", widget.NewEntry()),
		widget.NewFormItem("Air Baudrate", widget.NewEntry()),
		widget.NewFormItem("Packet Length", widget.NewEntry()),
		widget.NewFormItem("Channel RSSI", widget.NewCheck("Enable Channel RSSI", nil)),
		widget.NewFormItem("Transmission Power", widget.NewEntry()),
		widget.NewFormItem("Channel", widget.NewEntry()),
		widget.NewFormItem("Packet RSSI", widget.NewCheck("Enable Packet RSSI", nil)),
		widget.NewFormItem("Transmission Mode", widget.NewEntry()),
		widget.NewFormItem("Relay Enable", widget.NewCheck("Enable Relay", nil)),
		widget.NewFormItem("LBT Enable", widget.NewCheck("Enable LBT", nil)),
		widget.NewFormItem("WOR Role", widget.NewEntry()),
		widget.NewFormItem("WOR Cycle", widget.NewEntry()),
		widget.NewFormItem("Key", widget.NewEntry()),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Wireless", container.NewPadded(w.wirelessForm)),
		container.NewTabItem("Advanced", widget.NewLabel("Advanced settings...")),
	)

	readBtn := widget.NewButtonWithIcon("Read from Device", theme.ViewRefreshIcon(), w.handleReadParams)
	saveBtn := widget.NewButtonWithIcon("Save to Device", theme.DocumentSaveIcon(), w.handleSaveParams)

	return container.NewBorder(
		container.NewHBox(readBtn, saveBtn),
		nil, nil, nil,
		tabs,
	)
}

// handleReadParams uses the AT+LORA command to retrieve the details from the E90
func (w *MainWindow) handleReadParams() {
	if w.atClient == nil {
		dialog.ShowInformation("Get LoRa parameters", "AT Client must be initialised before using AT commands", w.window)
		return
	}

	loraparams, err := w.atClient.GetLora()
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}

	w.wirelessForm.Items[0].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.ModuleAddress))
	w.wirelessForm.Items[1].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.NetId))
	w.wirelessForm.Items[2].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.AirBaud, "bps"))
	w.wirelessForm.Items[3].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.PacketLength, "bytes"))
	w.wirelessForm.Items[4].Widget.(*widget.Check).SetChecked(loraparams.ChannelRssi)
	w.wirelessForm.Items[5].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.TransmissionPower))
	w.wirelessForm.Items[6].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.Channel))
	w.wirelessForm.Items[7].Widget.(*widget.Check).SetChecked(loraparams.PacketRssi)
	w.wirelessForm.Items[8].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.TransmissionMode))
	w.wirelessForm.Items[9].Widget.(*widget.Check).SetChecked(loraparams.RelayEnable)
	w.wirelessForm.Items[10].Widget.(*widget.Check).SetChecked(loraparams.LBTEnable)
	w.wirelessForm.Items[11].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.WorRole))
	w.wirelessForm.Items[12].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.WorCycle))
	w.wirelessForm.Items[13].Widget.(*widget.Entry).SetText(fmt.Sprint(loraparams.Key))
}

// TODO: create dropdowns for UI for this.
func (w *MainWindow) handleSaveParams() {
	dialog.ShowInformation("Saved", "Configuration saved to device", w.window)
}
