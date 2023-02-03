package client

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"strings"
	"sync"

	"github.com/xackery/eqemupatchergo/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Client wraps the entire UI
type Client struct {
	mu          sync.RWMutex
	url         string
	currentPath string
	progressBar *widget.ProgressBar
	canvas      fyne.CanvasObject
	mainCanvas  fyne.CanvasObject
	patchButton *widget.Button
	playButton  *widget.Button
	splashImage *canvas.Image
	statusLabel *widget.Label
	window      fyne.Window
	cfg         *config.Config
}

// New creates a new client
func New(window fyne.Window, url string) (*Client, error) {
	var err error
	c := &Client{
		window: window,
		url:    url,
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cfg, err = config.New(context.Background())
	if err != nil {
		return nil, fmt.Errorf("config.new: %w", err)
	}

	c.currentPath, err = os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("wd invalid: %w", err)
	}
	c.progressBar = widget.NewProgressBar()

	//c.currentPath = `C:\src\eqp\client\zones`

	c.patchButton = widget.NewButton("Patch", c.onPatchButton)
	c.playButton = widget.NewButton("Play", c.onPlayButton)

	c.statusLabel = widget.NewLabel("Status")
	c.statusLabel.Wrapping = fyne.TextWrapBreak
	c.statusLabel.Alignment = fyne.TextAlignCenter

	img, err := png.Decode(bytes.NewReader(RoFImage.Content()))
	if err != nil {
		return nil, fmt.Errorf("png decode: %w", err)
	}
	c.splashImage = canvas.NewImageFromImage(img)
	c.splashImage.FillMode = canvas.ImageFillOriginal

	c.mainCanvas = container.NewVBox(
		c.splashImage,
		container.NewBorder(
			nil,
			nil,
			c.patchButton,
			c.playButton,
		),
		c.progressBar,
		c.statusLabel,
	)

	c.canvas = c.mainCanvas

	c.logf("Press Patch to download.")
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

func (c *Client) addProgress(amount float64) {
	c.progressBar.Value += amount
	if c.progressBar.Value > 1 {
		fmt.Printf("progress > 1: %0.2f\n", c.progressBar.Value)
		c.progressBar.Value = 1
	}
	c.progressBar.SetValue(c.progressBar.Value)
}

func (c *Client) onPatchButton() {
	c.addProgress(1)
}

func (c *Client) onPlayButton() {
	c.addProgress(1)
}

// Parse will parse a []byte and turn it into the first element
func Parse(in []byte) string {
	lines := strings.Split(string(in), "\n")
	out := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		out = line
		break
	}
	return out
}
