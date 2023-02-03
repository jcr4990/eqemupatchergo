package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2/app"
	"github.com/xackery/eqgzi-manager/client"
)

var (
	Version string
)

func main() {
	if Version == "" {
		Version = string(client.VersionText.Content())
	}
	log.Println("initializing", Version)

	a := app.New()

	serverNameRaw := client.NameText.Content()
	lines := strings.Split(string(serverNameRaw), "\n")
	serverName := "eqemupatcher"
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		serverName = line
		break
	}

	w := a.NewWindow(fmt.Sprintf("%s v%s", serverName, Version))
	c, err := client.New(w)
	if err != nil {
		fmt.Println("client new:", err)
		os.Exit(1)
	}

	w.SetContent(c.GetContent())
	w.CenterOnScreen()
	w.ShowAndRun()
}
