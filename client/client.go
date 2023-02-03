package client

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/xackery/eqgzi-manager/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Client wraps the entire UI
type Client struct {
	mu             sync.RWMutex
	currentPath    string
	progress       float64
	canvas         fyne.CanvasObject
	mainCanvas     fyne.CanvasObject
	downloadButton *widget.Button
	statusLabel    *widget.Label
	window         fyne.Window
	cfg            *config.Config
}

// New creates a new client
func New(window fyne.Window) (*Client, error) {
	var err error
	c := &Client{
		window: window,
	}

	c.cfg, err = config.New(context.Background())
	if err != nil {
		return nil, fmt.Errorf("config.new: %w", err)
	}

	c.currentPath, err = os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("wd invalid: %w", err)
	}

	//c.currentPath = `C:\src\eqp\client\zones`

	c.downloadButton = widget.NewButton("Download", c.onDownloadButton)

	c.statusLabel = widget.NewLabel("Status")

	c.mainCanvas = container.NewVBox(
		c.downloadButton,
		c.statusLabel,
	)

	c.canvas = c.mainCanvas

	return c, nil
}

// GetContent returns the current canvas, and is used by SetContent
func (c *Client) GetContent() fyne.CanvasObject {
	return c.canvas
}

func (c *Client) logf(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	fmt.Println(text)
	c.statusLabel.SetText(text)
}

func (c *Client) addProgress(amount float64) float64 {
	c.progress += amount

	if c.progress > 1 {
		fmt.Printf("progress > 1: %0.2f\n", c.progress)
		c.progress = 1
	}
	return c.progress
}

func (c *Client) onDownloadButton() {

}
