package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	ipLabel := widget.NewLabel("Prem el botó per obtenir la IP...")
	ipLabel.Alignment = fyne.TextAlignCenter

	refreshBtn := widget.NewButton("Refrescar IP", func() {
		ipLabel.SetText("Obtenint IP...")
		ip, err := getPublicIP()
		if err != nil {
			ipLabel.SetText(fmt.Sprintf("Error: %v", err))
		} else {
			ipLabel.SetText(fmt.Sprintf("La teva adreça IP pública és: %s", ip))
		}
	})

	myWindow.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("Consultor d'IP Pública", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		ipLabel,
		refreshBtn,
	))

	myWindow.Resize(fyne.NewSize(300, 150))
	myWindow.ShowAndRun()
}
