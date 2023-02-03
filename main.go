package main

import (
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/xackery/eqemupatchergo/client"
)

func main() {
	version := string(client.VersionText.Content())
	log.Println("initializing", version)

	a := app.New()

	serverName := client.Parse(client.NameText.Content())
	url := client.Parse(client.UrlText.Content())
	w := a.NewWindow(fmt.Sprintf("%s v%s", serverName, version))
	c, err := client.New(w, url)
	if err != nil {
		fmt.Println("client new:", err)
		os.Exit(1)
	}

	w.SetContent(c.GetContent())
	w.Resize(fyne.NewSize(305, 371))
	w.CenterOnScreen()
	w.ShowAndRun()
}
