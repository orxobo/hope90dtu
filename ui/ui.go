package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"hope90dtu/device"
	"hope90dtu/device/atcommands"
)

type MainWindow struct {
	app    fyne.App
	window fyne.Window
	device *device.E90Device

	// Connection tab
	statusLabel       *widget.Label
	statusIcon        *widget.Icon
	serialPortSel     *widget.Select
	baudSel           *widget.Select
	ipEntry           *widget.Entry
	portEntry         *widget.Entry
	connTypeRadio     *widget.RadioGroup
	connectBtn        *widget.Button
	disconnectBtn     *widget.Button
	connectionMonitor *widget.Entry
	actWaiting        *widget.Activity

	// AT Terminal tab
	atResponse      *widget.Entry
	initATClientBtn *widget.Button
	atSelect        *widget.Select
	atClient        *atcommands.ATClient

	// Monitor tab
	monitorText   *widget.Entry
	monitorBuffer []string
	listenerBtn   *widget.Button

	// Parameters tab
	wirelessForm *widget.Form

	// Test tab
	testBtn *widget.Button
}

func NewMainWindow(app fyne.App) *MainWindow {

	w := &MainWindow{
		app:           app,
		window:        app.NewWindow("HOPE90-DTU Configuration Tool"),
		monitorBuffer: make([]string, 0),
	}

	w.makeUI()
	w.appendToMonitor("E90-DTU Configuration Tool started")
	w.window.Resize(fyne.NewSize(800, 600))
	w.window.CenterOnScreen()
	return w
}

func (w *MainWindow) makeUI() {
	// Menu
	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Save Connection Settings", w.saveConnectionSettings),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Exit", w.app.Quit),
		),
		fyne.NewMenu("Help",
			fyne.NewMenuItem("About", w.showAbout),
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

	w.window.SetContent(tabs)
}

func (w *MainWindow) Show() {
	w.window.Show()
}

type connectionSettings struct {
	ConnectionType string `yaml:"connectiontype"` // Serial or Network
	IPAddress      string `yaml:"ipaddress"`      // IPv4 address
	UDPPort        int    `yaml:"udpport"`        // 1-65535
	SerialPort     string `yaml:"serialport"`     // COM1 /dev/ttyUSB0 etc
	BaudRate       int    `yaml:"baudrate"`       // UINT16 9600
	ATPrefix       string `yaml:"atprefix"`       // AT
}

func (w *MainWindow) saveConnectionSettings() {

}

func (w *MainWindow) loadConnectionSettings() {

}
