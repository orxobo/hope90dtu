package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (w *MainWindow) makeMonitorTab() fyne.CanvasObject {
	w.monitorText = widget.NewMultiLineEntry()
	//w.monitorText.TextStyle = fyne.TextStyle{Monospace: true}
	w.monitorText.OnChanged = func(s string) {
		//w.monitorText.Undo()
	}

	var ctxCancel context.CancelFunc
	numEntry := widget.NewEntry()
	numEntry.SetText("120")
	countDown := widget.NewLabel("Seconds")
	cancelBtn := widget.NewButton("Cancel", func() {
		if ctxCancel != nil {
			ctxCancel()
		}
	})
	w.listenerBtn = widget.NewButton("Start Listener", func() {
		timeout, _ := strconv.Atoi(numEntry.Text)
		if timeout == 0 {
			return
		}
		w.appendToMonitorf("Starting UDP Listener for %d seconds", timeout)
		var ctx context.Context
		ctx, ctxCancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		go w.device.UDPListener(ctx)
		go w.startCountDown(ctx, countDown)
	})
	w.listenerBtn.Disable()

	checkPointBtn := widget.NewButton("Add Checkpoint", func() {
		lineBreak := strings.Repeat("=", 20)
		w.appendToMonitorf("%s Check Point %s", lineBreak, lineBreak)
	})

	clearBtn := widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
		w.monitorBuffer = []string{}
		w.monitorText.SetText("")
	})

	toolbar := container.NewHBox(
		w.listenerBtn,
		numEntry,
		countDown,
		cancelBtn,
		layout.NewSpacer(),
		checkPointBtn,
		clearBtn,
	)

	return container.NewBorder(toolbar, nil, nil, nil, w.monitorText)
}

func (w *MainWindow) startCountDown(ctx context.Context, countDown *widget.Label) {
	deadline, _ := ctx.Deadline()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			countDown.SetText("Seconds")
			w.appendToMonitor("UDP Listener finished")
			return
		case <-ticker.C:
			remaining := time.Until(deadline)
			countDown.SetText(fmt.Sprintf("Remaining: %d seconds", int(remaining.Seconds())))
		}
	}
}

func (w *MainWindow) appendToMonitorf(line string, a ...any) {
	w.appendToMonitor(fmt.Sprintf(line, a...))
}

func (w *MainWindow) appendToMonitor(line string) {
	timestamp := time.Now().Format("15:04:05.000")

	w.monitorBuffer = append(w.monitorBuffer, fmt.Sprintf("[%s] %s", timestamp, line))

	// Limit buffer size
	if len(w.monitorBuffer) > 1000 {
		w.monitorBuffer = w.monitorBuffer[len(w.monitorBuffer)-1000:]
	}

	fullText := strings.Join(w.monitorBuffer, "\n")
	w.monitorText.SetText(fullText)
}
