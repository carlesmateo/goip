package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GoIP Public IP")

	var lastIP string
	var mu sync.Mutex

	// Use binding to update label safely from any thread
	ipBinding := binding.NewString()
	ipBinding.Set("Prem el botó per obtenir la IP...")

	ipLabel := widget.NewLabelWithData(ipBinding)
	ipLabel.Alignment = fyne.TextAlignCenter

	updateIP := func() {
		ipBinding.Set("Obtenint IP...")
		go func() {
			ip, err := getPublicIP()
			if err != nil {
				ipBinding.Set(fmt.Sprintf("Error: %v", err))
				return
			}

			mu.Lock()
			currentLast := lastIP
			lastIP = ip
			mu.Unlock()

			ipBinding.Set(fmt.Sprintf("La teva adreça IP pública és: %s", ip))
			_ = currentLast
		}()
	}

	refreshBtn := widget.NewButton("Refrescar IP", func() {
		updateIP()
	})

	// Background Ticker
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			ip, err := getPublicIP()
			if err != nil {
				continue
			}

			mu.Lock()
			oldIP := lastIP
			if ip != lastIP {
				lastIP = ip
			}
			mu.Unlock()

			if oldIP != "" && oldIP != ip {
				ipBinding.Set(fmt.Sprintf("La teva adreça IP pública és: %s", ip))
				// Note: Window focus requires UI thread.
				// Since direct access caused a crash and Driver().RunOnUIThread was reported undefined,
				// we are currently only updating the text which is thread-safe via binding.
				// To implement focus, we would need a valid way to dispatch to main thread.
				// myWindow.Show()
				// myWindow.RequestFocus()
			} else if oldIP == "" {
				ipBinding.Set(fmt.Sprintf("La teva adreça IP pública és: %s", ip))
			}
		}
	}()

	myWindow.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("Consultor d'IP Pública", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		ipLabel,
		refreshBtn,
	))

	myWindow.Resize(fyne.NewSize(300, 150))
	myWindow.ShowAndRun()
}
