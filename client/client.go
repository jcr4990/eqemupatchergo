package client

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/xackery/eqemupatchergo/config"
	"gopkg.in/yaml.v3"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Client wraps the entire UI
type Client struct {
	mu              sync.RWMutex
	cancel          chan bool
	url             string
	currentPath     string
	clientVersion   string
	cacheFileList   *FileList
	logScroll       *container.Scroll
	logLabel        *widget.Label
	copyLogButton   *widget.Button
	progressBar     *widget.ProgressBar
	canvas          fyne.CanvasObject
	mainCanvas      fyne.CanvasObject
	patchButton     *widget.Button
	playButton      *widget.Button
	splashImage     *canvas.Image
	statusLabel     *widget.Label
	window          fyne.Window
	cfg             *config.Config
	isAutoPatchPlay bool
}

// New creates a new client
func New(window fyne.Window, url string) (*Client, error) {
	var err error
	url = strings.TrimSuffix(url, "/")
	c := &Client{
		window:        window,
		url:           url,
		clientVersion: "rof",
		cancel:        make(chan bool, 3),
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

	c.copyLogButton = widget.NewButtonWithIcon("", theme.ContentCopyIcon(), c.onCopyLog)

	c.patchButton = widget.NewButton("Patch", c.onPatchButton)
	c.playButton = widget.NewButton("Play", c.onPlayButton)

	c.logLabel = widget.NewLabel("")
	//c.logLabel.Wrapping = fyne.TextTruncate

	c.logScroll = container.NewScroll(
		c.logLabel,
	)
	c.logScroll.Hide()

	c.statusLabel = widget.NewLabel("")
	c.statusLabel.Wrapping = fyne.TextTruncate
	c.statusLabel.Alignment = fyne.TextAlignCenter
	c.progressBar.Hide()
	c.statusLabel.Hide()

	img, err := png.Decode(bytes.NewReader(RoFImage.Content()))
	if err != nil {
		return nil, fmt.Errorf("png decode: %w", err)
	}
	c.splashImage = canvas.NewImageFromImage(img)
	c.splashImage.FillMode = canvas.ImageFillOriginal

	c.mainCanvas = container.NewBorder(
		nil,
		//bottom
		container.NewVBox(
			container.NewHBox(
				c.patchButton,
				layout.NewSpacer(),
				c.playButton,
			),
			/*container.NewBorder(
				nil,
				nil,
				c.patchButton,
				c.playButton,
			),*/
			c.progressBar,
			c.statusLabel,
		),
		//left
		nil,
		//right
		nil,
		//remaining
		container.NewCenter(
			c.splashImage,
		),
		c.logScroll,
	)
	go c.asyncVersionCheck()

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
	/*c.logHistory = append(c.logHistory, text)
	for len(c.logHistory) > 25 {
		c.logHistory = c.logHistory[1:]
	}

	log := ""
	for _, line := range c.logHistory {
		if log == "" {
			log = line
			continue
		}
		log += "\n" + line
	}

	c.logLabel.SetText(log)*/
	if len(c.logLabel.Text) == 0 {
		c.logLabel.SetText(text)
	} else {
		c.logLabel.SetText(c.logLabel.Text + "\n" + text)
	}

	c.logScroll.ScrollToBottom()
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

func (c *Client) onCopyLog() {

}

func (c *Client) asyncVersionCheck() {
	err := c.refreshFileList()
	if err != nil {
		fmt.Println("ignoring failure to patch:", err)
		return
	}

	c.mu.RLock()
	version := c.cacheFileList.Version
	myVersion := c.cfg.LastPatchedVersion
	autoPlay := c.cfg.AutoPlay
	autoPatch := c.cfg.AutoPatch
	c.mu.RUnlock()
	isAutoPlay := strings.ToLower(autoPlay) == "true"
	isAutoPatch := strings.ToLower(autoPatch) == "true"
	if myVersion != version {
		c.patchButton.Importance = widget.HighImportance
		if isAutoPatch {
			if isAutoPlay {
				c.mu.Lock()
				c.isAutoPatchPlay = true
				c.mu.Unlock()
			}
			c.onPatchButton()
		}
		return
	}

	if isAutoPlay {
		c.onPlayButton()
	}

}

func (c *Client) refreshFileList() error {
	client := http.DefaultClient
	url := fmt.Sprintf("%s/%s/filelist_%s.yml", c.url, c.clientVersion, c.clientVersion)
	fmt.Println("Downloading", url)
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()
	fileList := &FileList{}

	err = yaml.NewDecoder(resp.Body).Decode(fileList)
	if err != nil {
		return fmt.Errorf("decode filelist: %w", err)
	}
	fmt.Println("patch version is", fileList.Version)
	c.mu.Lock()
	c.cacheFileList = fileList
	c.mu.Unlock()

	return nil
}
