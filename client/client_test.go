package client

import (
	"testing"

	"fyne.io/fyne/v2/test"
)

func TestOnDownload(t *testing.T) {
	app := test.NewApp()
	window := app.NewWindow("Test")
	c, err := New(window)
	if err != nil {
		t.Fatalf("new: %s", err)
	}
	c.onDownloadButton()
}
