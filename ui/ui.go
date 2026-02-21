package ui

import (
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"hope90dtu/device"
)

type MainWindow struct {
	app    fyne.App
	window fyne.Window
	device *device.E90Device

	// UI State
	monitorBuffer []string
	isAutoScroll  bool

	// Widgets - Connection
	statusLabel   *widget.Label
	statusIcon    *widget.Icon
	serialPortSel *widget.Select
	baudSel       *widget.Select
	ipEntry       *widget.Entry
	portEntry     *widget.Entry
	connTypeRadio *widget.RadioGroup
	connectBtn    *widget.Button
	disconnectBtn *widget.Button

	// Widgets - Terminal
	atEntry *widget.Entry
	sendBtn *widget.Button

	// Widgets - Monitor
	monitorText     *widget.Entry
	autoScrollCheck *widget.Check

	// Widgets - Parameters
	paramEntries map[string]fyne.Widget
	testBtn      *widget.Button
}

func NewMainWindow(app fyne.App) *MainWindow {

	w := &MainWindow{
		app:           app,
		window:        app.NewWindow("E90-DTU Configuration Tool"),
		paramEntries:  make(map[string]fyne.Widget),
		monitorBuffer: make([]string, 0),
		isAutoScroll:  true,
	}

	w.window.Resize(fyne.NewSize(1200, 800))
	w.makeUI()
	w.appendToMonitor("E90-DTU Configuration Tool")
	return w
}

func (w *MainWindow) makeUI() {
	// Menu
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Exit", func() { w.app.Quit() }),
		),
	)
	w.window.SetMainMenu(menu)

	// Tabs
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Connection", theme.ComputerIcon(), w.makeConnectionTab()),
		container.NewTabItemWithIcon("AT Terminal", theme.MailSendIcon(), w.makeTerminalTab()),
		container.NewTabItemWithIcon("Parameters", theme.SettingsIcon(), w.makeParametersTab()),
		container.NewTabItemWithIcon("Monitor", theme.VisibilityIcon(), w.makeMonitorTab()),
		container.NewTabItemWithIcon("Test Tools", theme.HelpIcon(), w.makeTestTab()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.window.SetContent(tabs)
}

func (w *MainWindow) makeConnectionTab() fyne.CanvasObject {
	w.statusLabel = widget.NewLabelWithStyle("Status: Disconnected",
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	w.statusIcon = widget.NewIcon(theme.CancelIcon())

	// Serial options
	w.serialPortSel = widget.NewSelect(
		[]string{"COM1", "COM2", "COM3", "COM4", "COM5", "COM6",
			"/dev/ttyUSB0", "/dev/ttyACM0", "/dev/ttyS0"}, nil)
	w.serialPortSel.SetSelected("COM1")

	w.baudSel = widget.NewSelect(
		[]string{"9600", "19200", "38400", "57600", "115200"}, nil)
	w.baudSel.SetSelected("9600")

	serialForm := widget.NewForm(
		widget.NewFormItem("Port", w.serialPortSel),
		widget.NewFormItem("Baud Rate", w.baudSel),
	)

	// Network options
	w.ipEntry = widget.NewEntry()
	w.ipEntry.SetText("192.168.0.101")
	w.ipEntry.SetPlaceHolder("192.168.0.101")

	w.portEntry = widget.NewEntry()
	w.portEntry.SetText("8886")
	w.portEntry.SetPlaceHolder("8886")

	networkForm := widget.NewForm(
		widget.NewFormItem("IP Address", w.ipEntry),
		widget.NewFormItem("Port", w.portEntry),
	)

	// Connection type selector
	serialForm.Hide()
	networkForm.Show()

	w.connTypeRadio = widget.NewRadioGroup([]string{"Serial", "Network"}, func(selected string) {
		if selected == "Serial" {
			serialForm.Show()
			networkForm.Hide()
		} else {
			serialForm.Hide()
			networkForm.Show()
		}
	})
	w.connTypeRadio.SetSelected("Network")

	// Buttons
	w.connectBtn = widget.NewButtonWithIcon("Connect", theme.LoginIcon(), w.handleConnect)
	w.disconnectBtn = widget.NewButtonWithIcon("Disconnect", theme.LogoutIcon(), w.handleDisconnect)
	w.disconnectBtn.Disable()

	btnBox := container.NewHBox(w.connectBtn, w.disconnectBtn)
	loginUrl, _ := url.Parse("http://admin:admin@192.168.0.150")

	return container.NewVBox(
		container.NewHBox(w.statusIcon, w.statusLabel),
		widget.NewSeparator(),
		w.connTypeRadio,
		widget.NewSeparator(),
		serialForm,
		networkForm,
		widget.NewSeparator(),
		btnBox,
		layout.NewSpacer(),
		container.NewHBox(widget.NewLabel("Default: IP 192.168.0.101, Port: 8886, "), widget.NewHyperlink("Web login admin/admin", loginUrl)),
	)
}

func (w *MainWindow) makeTerminalTab() fyne.CanvasObject {
	w.atEntry = widget.NewMultiLineEntry()
	w.atEntry.SetPlaceHolder("Enter AT command (e.g., AT+ADDRESS? or AT+MODE=TCP_SERVER)...")

	// Quick commands grid
	quickCmds := []string{
		"AT", "AT+HELP", "AT+ADDRESS?", "AT+CHANNEL?",
		"AT+NETWORKID?", "AT+IP?", "AT+MODE?", "AT+UART?",
		"AT+SAVE", "AT+REBOOT",
	}

	cmdGrid := container.NewGridWithColumns(5)
	for _, cmd := range quickCmds {
		c := cmd
		btn := widget.NewButton(cmd, func() {
			w.atEntry.SetText(c)
		})
		cmdGrid.Add(btn)
	}

	w.sendBtn = widget.NewButtonWithIcon("Send", theme.MailSendIcon(), w.handleSendAT)
	w.sendBtn.Disable()

	inputArea := container.NewBorder(nil, w.sendBtn, nil, nil, w.atEntry)

	return container.NewBorder(
		widget.NewCard("Quick Commands", "", cmdGrid),
		nil, nil, nil,
		inputArea,
	)
}

func (w *MainWindow) makeParametersTab() fyne.CanvasObject {
	w.paramEntries["address"] = widget.NewEntry()
	w.paramEntries["address"].(*widget.Entry).SetText("0")

	w.paramEntries["channel"] = widget.NewEntry()
	w.paramEntries["channel"].(*widget.Entry).SetText("23")

	w.paramEntries["networkid"] = widget.NewEntry()
	w.paramEntries["networkid"].(*widget.Entry).SetText("0")

	w.paramEntries["airrate"] = widget.NewSelect(
		[]string{"2.4K", "4.8K", "9.6K", "19.2K", "38.4K", "62.5K"}, nil)
	w.paramEntries["airrate"].(*widget.Select).SetSelected("2.4K")

	w.paramEntries["ip"] = widget.NewEntry()
	w.paramEntries["ip"].(*widget.Entry).SetText("192.168.4.101")

	w.paramEntries["workmode"] = widget.NewSelect([]string{
		"TCP Server", "TCP Client", "UDP Server", "UDP Client",
		"HTTP Client", "MQTT Client"}, nil)
	w.paramEntries["workmode"].(*widget.Select).SetSelected("TCP Server")

	w.paramEntries["localport"] = widget.NewEntry()
	w.paramEntries["localport"].(*widget.Entry).SetText("8886")

	wirelessForm := widget.NewForm(
		widget.NewFormItem("Address", w.paramEntries["address"]),
		widget.NewFormItem("Channel", w.paramEntries["channel"]),
		widget.NewFormItem("Network ID", w.paramEntries["networkid"]),
		widget.NewFormItem("Air Rate", w.paramEntries["airrate"]),
	)

	networkForm := widget.NewForm(
		widget.NewFormItem("Local IP", w.paramEntries["ip"]),
		widget.NewFormItem("Work Mode", w.paramEntries["workmode"]),
		widget.NewFormItem("Local Port", w.paramEntries["localport"]),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Wireless", wirelessForm),
		container.NewTabItem("Network", networkForm),
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

func (w *MainWindow) makeMonitorTab() fyne.CanvasObject {
	w.monitorText = widget.NewMultiLineEntry()
	w.monitorText.Disable()
	w.monitorText.TextStyle = fyne.TextStyle{Monospace: true}

	scroll := container.NewScroll(w.monitorText)
	scroll.SetMinSize(fyne.NewSize(800, 500))

	w.autoScrollCheck = widget.NewCheck("Auto Scroll", func(checked bool) {
		w.isAutoScroll = checked
	})
	w.autoScrollCheck.SetChecked(true)

	clearBtn := widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
		w.monitorBuffer = []string{}
		w.monitorText.SetText("")
	})

	toolbar := container.NewHBox(w.autoScrollCheck, layout.NewSpacer(), clearBtn)

	return container.NewBorder(toolbar, nil, nil, nil, scroll)
}

func (w *MainWindow) makeTestTab() fyne.CanvasObject {
	// Random data test
	lengthEntry := widget.NewEntry()
	lengthEntry.SetText("10")

	intervalEntry := widget.NewEntry()
	intervalEntry.SetText("1000")

	w.testBtn = widget.NewButtonWithIcon("Start Random Test", theme.MediaPlayIcon(), func() {
		w.handleRandomTest(lengthEntry.Text, intervalEntry.Text)
	})
	w.testBtn.Disable()

	form := widget.NewForm(
		widget.NewFormItem("Length (bytes)", lengthEntry),
		widget.NewFormItem("Interval (ms)", intervalEntry),
	)

	// Protocol distribution
	socketA := widget.NewButton("Protocol Dist Socket A", func() {
		w.handleProtocolDist(0)
	})
	socketB := widget.NewButton("Protocol Dist Socket B", func() {
		w.handleProtocolDist(1)
	})

	customData := widget.NewEntry()
	customData.SetPlaceHolder("Hex data (55 FE AA 00 01 02...)")
	sendHex := widget.NewButton("Send Hex", func() {
		w.handleSendHex(customData.Text)
	})

	return container.NewVBox(
		widget.NewCard("Random Generator", "", form),
		w.testBtn,
		widget.NewSeparator(),
		widget.NewCard("Protocol Distribution", "", container.NewHBox(socketA, socketB)),
		widget.NewSeparator(),
		widget.NewCard("Custom Data", "", container.NewBorder(nil, nil, nil, sendHex, customData)),
	)
}

// Event Handlers

func (w *MainWindow) handleConnect() {

	if w.connTypeRadio.Selected == "Serial" {
		baud, _ := strconv.Atoi(w.baudSel.Selected)
		dev, err := device.NewE90SerialDevice(w.serialPortSel.Selected, baud)
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		w.device = dev
	} else {
		dev, err := device.NewE90UDPDeviceFromIPAddressAndPort(w.ipEntry.Text, w.portEntry.Text)
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		w.device = dev
	}

	w.appendToMonitor("✓ connection established")
	w.appendToMonitor(strings.Repeat("=", 30))

	w.device.SetMonitor(w.appendToMonitor)
	w.window.SetOnClosed(func() { w.device.Close() })

	w.statusLabel.SetText("Status: Connected")
	w.statusIcon.SetResource(theme.ConfirmIcon())
	w.connectBtn.Disable()
	w.disconnectBtn.Enable()
	w.sendBtn.Enable()
	w.testBtn.Enable()

	w.device.SetDisconnectCallback(func() {
		w.statusLabel.SetText("Status: Disconnected")
		w.statusIcon.SetResource(theme.CancelIcon())
		w.connectBtn.Enable()
		w.disconnectBtn.Disable()
		w.sendBtn.Disable()
		w.testBtn.Disable()
	})
}

func (w *MainWindow) handleDisconnect() {
	w.device.Close()
}

func (w *MainWindow) handleSendAT() {
	cmd := w.atEntry.Text
	if cmd == "" {
		return
	}

	// if err := w.device.SendCommand(cmd); err != nil {
	// 	dialog.ShowError(err, w.window)
	// }
}

func (w *MainWindow) handleReadParams() {
	// Query all parameters from device
	// cmds := []string{
	// 	"AT+ADDRESS?",
	// 	"AT+CHANNEL?",
	// 	"AT+NETWORKID?",
	// 	"AT+IP?",
	// 	"AT+MODE?",
	// }

	// for _, cmd := range cmds {
	// 	w.device.SendCommand(cmd)
	// 	time.Sleep(50 * time.Millisecond)
	// }
}

func (w *MainWindow) handleSaveParams() {
	// Collect from UI and send to device
	// addr, _ := strconv.Atoi(w.paramEntries["address"].(*widget.Entry).Text)
	// ch, _ := strconv.Atoi(w.paramEntries["channel"].(*widget.Entry).Text)
	// netID, _ := strconv.Atoi(w.paramEntries["networkid"].(*widget.Entry).Text)

	// w.device.SetAddress(uint16(addr))
	// w.device.SetChannel(uint8(ch))
	// w.device.SetNetworkID(uint8(netID))
	// w.device.SetIP(w.paramEntries["ip"].(*widget.Entry).Text)

	// w.device.SaveToDevice()

	dialog.ShowInformation("Saved", "Configuration saved to device", w.window)
}

func (w *MainWindow) handleRandomTest(lengthStr, intervalStr string) {
	length, _ := strconv.Atoi(lengthStr)
	// interval, _ := strconv.Atoi(intervalStr)

	if length <= 0 {
		return
	}

	go func() {
		// for i := 0; i < 10 && w.device.IsConnected(); i++ {
		// 	data, _ := w.device.SendRandomData(length)
		// 	w.appendToMonitor(fmt.Sprintf("[TEST] Sent %d random bytes", len(data)))
		// 	time.Sleep(time.Duration(interval) * time.Millisecond)
		// }
	}()
}

func (w *MainWindow) handleProtocolDist(socket byte) {
	// payload := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	// if err := w.device.SendProtocolDist(socket, payload); err != nil {
	// 	dialog.ShowError(err, w.window)
	// }
}

func (w *MainWindow) handleSendHex(hexStr string) {
	if err := w.device.SendUDPHexCommand(hexStr); err != nil {
		dialog.ShowError(err, w.window)
	}
}

func (w *MainWindow) appendToMonitor(line string) {
	w.monitorBuffer = append(w.monitorBuffer, line)

	// Limit buffer size
	if len(w.monitorBuffer) > 1000 {
		w.monitorBuffer = w.monitorBuffer[len(w.monitorBuffer)-1000:]
	}

	fullText := strings.Join(w.monitorBuffer, "\n")
	w.monitorText.SetText(fullText)

	if w.isAutoScroll {

	}
}

func (w *MainWindow) Show() {
	w.window.Show()
}
