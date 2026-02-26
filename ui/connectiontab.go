// UI for the Connection tab
package ui

import (
	"context"
	"fmt"
	"hope90dtu/device"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (w *MainWindow) makeConnectionTab() fyne.CanvasObject {
	w.statusLabel = widget.NewLabelWithStyle("Status: Disconnected", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	w.statusIcon = widget.NewIcon(theme.CancelIcon())

	// --- Serial Form ---
	w.serialPortSel = widget.NewSelect(device.SerialPorts(), nil)
	w.serialPortSel.SetSelectedIndex(0)

	w.baudSel = widget.NewSelect(device.BaudRates(), nil)
	w.baudSel.SetSelectedIndex(0)

	serialForm := widget.NewForm(
		widget.NewFormItem("Port", w.serialPortSel),
		widget.NewFormItem("Baud Rate", w.baudSel),
	)

	// --- Network Form ---
	w.ipEntry = widget.NewEntry()
	w.ipEntry.SetPlaceHolder("Factory Default: 192.168.4.101")

	w.portEntry = widget.NewEntry()
	w.portEntry.SetPlaceHolder("Factory Default: 8886")

	currSub := widget.NewEntry()
	currSub.SetText(getLocalSubnet().String())
	currSub.OnChanged = func(s string) {
		currSub.Undo()
		dialog.ShowInformation("Current Subnet", "This is the IP Address of your current subnet.\nClick 'Search Subnet' to try to find your E90-DTU.", w.window)
	}

	w.actWaiting = widget.NewActivity()
	w.actWaiting.Hide()

	subSearch := widget.NewButtonWithIcon("Search Subnet", theme.SearchIcon(), func() {
		subNet := net.ParseIP(currSub.Text)
		w.InitiateSearch(subNet)
	})
	subSearchContainer := container.NewBorder(nil, nil, nil, container.NewHBox(w.actWaiting, subSearch), currSub)

	networkForm := widget.NewForm(
		widget.NewFormItem("IP Address", w.ipEntry),
		widget.NewFormItem("UDP Port", w.portEntry),
		widget.NewFormItem("Current Subnet", subSearchContainer),
	)

	// --- Tab UI ---
	serialForm.Hide()
	networkForm.Show()

	w.connTypeRadio = widget.NewRadioGroup([]string{"Serial", "Network"}, func(selected string) {
		// TODO: if connected
		if selected == "Serial" {
			serialForm.Show()
			networkForm.Hide()
		} else {
			serialForm.Hide()
			networkForm.Show()
		}
	})
	w.connTypeRadio.SetSelected("Network")

	w.connectBtn = widget.NewButtonWithIcon("Connect", theme.LoginIcon(), w.handleConnect)
	w.disconnectBtn = widget.NewButtonWithIcon("Disconnect", theme.LogoutIcon(), w.device.Close)
	w.disconnectBtn.Disable()

	w.connectionMonitor = widget.NewMultiLineEntry()
	w.connectionMonitor.TextStyle = fyne.TextStyle{Monospace: true}
	w.connectionMonitor.OnChanged = func(s string) {
		//w.connectionMonitor.Undo()
	}

	loginUrl, _ := url.Parse("http://admin:admin@192.168.0.59")

	return container.NewBorder(
		container.NewVBox(
			container.NewHBox(w.statusIcon, w.statusLabel),
			widget.NewSeparator(),
			w.connTypeRadio,
			widget.NewSeparator(),
			widget.NewLabelWithStyle("Device Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			serialForm,
			networkForm,
			widget.NewSeparator(),
			container.NewHBox(w.connectBtn, w.disconnectBtn),
		),
		nil, nil, nil,
		container.NewBorder(
			nil,
			container.NewHBox(
				widget.NewLabel(fmt.Sprintf("IP %s, UDP Port: 8886", loginUrl.Hostname())),
				widget.NewHyperlink("Web login admin/admin", loginUrl),
			),
			nil, nil,
			w.connectionMonitor,
		),
	)
}

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

	w.appendToConnectMonitor("✓ Connection established")
	w.appendToConnectMonitor(strings.Repeat("=", 30))

	w.device.SetMonitor(w.appendToMonitor)
	w.window.SetOnClosed(w.device.Close)

	w.statusLabel.SetText("Status: Connected")
	w.statusIcon.SetResource(theme.ConfirmIcon())
	w.connectBtn.Disable()
	w.disconnectBtn.Enable()
	w.initATClientBtn.Enable()
	w.testBtn.Enable()
	w.listenerBtn.Enable()

	w.device.SetDisconnectCallback(func() {
		w.statusLabel.SetText("Status: Disconnected")
		w.statusIcon.SetResource(theme.CancelIcon())
		w.connectBtn.Enable()
		w.disconnectBtn.Disable()
		w.initATClientBtn.Disable()
		w.testBtn.Disable()
		w.listenerBtn.Disable()
	})
}

func (w *MainWindow) appendToConnectMonitorf(message string, a ...any) {
	w.appendToConnectMonitor(fmt.Sprintf(message, a...))
}

func (w *MainWindow) appendToConnectMonitor(message string) {
	w.connectionMonitor.Append(fmt.Sprintf("[%s] %s\n", time.Now().Format("15:04:05.000"), message))
	w.appendToMonitor(message)
}

func (w *MainWindow) InitiateSearch(subNet net.IP) {
	w.appendToConnectMonitorf("Search subnet %v started", subNet)
	w.actWaiting.Show()
	w.actWaiting.Start()

	go func() {

		foundIP, err := searchSubnet(subNet)
		if err != nil {
			dialog.ShowError(err, w.window)
			fyne.Do(func() { w.appendToConnectMonitor("E90-DTU was not found.") })
			return
		}
		fyne.Do(func() {
			w.ipEntry.SetText(foundIP.String())
			if w.portEntry.Text == "" {
				w.portEntry.SetText("8886")
			}
			w.appendToConnectMonitorf("E90-DTU found at %v", foundIP)
			w.actWaiting.Stop()
			w.actWaiting.Hide()
		})
	}()
}

// Get preferred outbound subnet of this machine
func getLocalSubnet() net.IP {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	lip := localAddr.IP
	return net.IPv4(lip[0], lip[1], lip[2], 0)
}

// searchSubnet searches all 254 computers on the local subnet
// for a port 80 that returns the fingerprint of the E90-DTU
func searchSubnet(subNet net.IP) (net.IP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	resultChan := make(chan net.IP, 254)
	errorChan := make(chan error, 254)
	maxConcurrent := 20
	semaphore := make(chan struct{}, maxConcurrent)

	ip := subNet.To4()
	if ip == nil {
		return nil, fmt.Errorf("invalid IPv4 address")
	}

	ip[3] = 0
	var wg sync.WaitGroup

	// check each IP
	for i := 1; i <= 254; i++ {
		wg.Add(1)
		go func(octet int) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				return
			}

			testIP := make(net.IP, len(ip))
			copy(testIP, ip)
			testIP[3] = byte(octet)

			found, err := checkIP(ctx, testIP)
			if err == nil && found != nil {
				select {
				case resultChan <- found:
				case <-ctx.Done():
				}
			} else if err != nil {
				select {
				case errorChan <- err:
				case <-ctx.Done():
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	select {
	case foundIP := <-resultChan:
		return foundIP, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// checkIP checks a single IP address for the E90-DTU fingerprint
func checkIP(ctx context.Context, ip net.IP) (net.IP, error) {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	url := fmt.Sprintf("http://%s", ip.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pageBody, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}

	if strings.Contains(string(pageBody), "ebyte_token") {
		return ip, nil
	}

	return nil, fmt.Errorf("E90-DTU not found at %s", ip.String())
}
